package studentEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"backend/type/tuple"
	"backend/util/orm"
	"reflect"
	"time"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleExamSubmit(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ExamQuestionSubmitRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate body
	if err := gut.Validate(body); err != nil {
		return err
	}

	// * get exam attempt with verification
	attemptDetails, err := r.database.P().ExamAttemptGetById(c.Context(), body.ExamAttemptId)
	if err != nil {
		return gut.Err(false, "exam attempt not found", err)
	}

	// * verify user owns this attempt
	if *attemptDetails.User.Id != *u.UserId {
		return gut.Err(false, "access denied", nil)
	}

	// * get exam question
	examQuestion, err := r.database.P().ExamQuestionGetById(c.Context(), body.ExamQuestionId)
	if err != nil {
		return gut.Err(false, "exam question not found", err)
	}

	// * verify question belongs to the exam
	if *examQuestion.ExamId != *attemptDetails.Exam.Id {
		return gut.Err(false, "question does not belong to this exam", nil)
	}

	// * execute student answer and check query
	var studentResult []map[string]any
	var checkResult []map[string]any
	checkQueryPassed := false
	now := time.Now()

	// * prepare result object
	result := &tuple.ExamSubmissionResult{
		ExecutionError:    "",
		PromptError:       "",
		PromptDescription: "",
		Rows:              []any{},
	}

	if attemptDetails.ExamAttempt.DatabaseName != nil {
		// * execute student answer query
		gorm, er := orm.Instance(*r.config.MysqlDsn, *attemptDetails.ExamAttempt.DatabaseName)
		if er != nil {
			return gut.Err(false, "failed to connect to exam database instance", er)
		}
		tx := gorm.Raw(*body.Answer).Scan(&studentResult)
		if tx.Error != nil {
			// * query failed to execute
			checkQueryPassed = false
			result.ExecutionError = tx.Error.Error()

			// * create submission with error
			submission, err := r.database.P().ExamSubmissionCreate(c.Context(), &psql.ExamSubmissionCreateParams{
				ExamQuestionId:    body.ExamQuestionId,
				ExamAttemptId:     body.ExamAttemptId,
				Answer:            body.Answer,
				Result:            result,
				CheckQueryPassed:  &checkQueryPassed,
				CheckQueryAt:      &now,
				CheckPromptPassed: nil,
				CheckPromptAt:     nil,
			})
			if err != nil {
				return gut.Err(false, "failed to create submission", err)
			}

			// * return success response with error result
			return c.JSON(response.Success(c, &payload.ExamQuestionSubmitResponse{
				Submission: &payload.ExamSubmission{
					Id:                submission.Id,
					ExamQuestionId:    submission.ExamQuestionId,
					ExamAttemptId:     submission.ExamAttemptId,
					Answer:            submission.Answer,
					CheckQueryPassed:  submission.CheckQueryPassed,
					CheckQueryAt:      submission.CheckQueryAt,
					CheckPromptPassed: submission.CheckPromptPassed,
					CheckPromptAt:     submission.CheckPromptAt,
					CreatedAt:         submission.CreatedAt,
					UpdatedAt:         submission.UpdatedAt,
				},
			}))
		}

		// * execute check query
		tx = gorm.Raw(*examQuestion.CheckQuery).Scan(&checkResult)
		if tx.Error == nil {
			// * compare results using reflect.DeepEqual
			checkQueryPassed = reflect.DeepEqual(studentResult, checkResult)
			result.Rows = make([]any, len(studentResult))
			for i, row := range studentResult {
				result.Rows[i] = row
			}
		}
	}

	// * create submission
	checkPromptPassed := false
	submission, err := r.database.P().ExamSubmissionCreate(c.Context(), &psql.ExamSubmissionCreateParams{
		ExamQuestionId:    body.ExamQuestionId,
		ExamAttemptId:     body.ExamAttemptId,
		Answer:            body.Answer,
		Result:            result,
		CheckQueryPassed:  &checkQueryPassed,
		CheckQueryAt:      &now,
		CheckPromptPassed: &checkPromptPassed,
		CheckPromptAt:     nil,
	})
	if err != nil {
		return gut.Err(false, "failed to create submission", err)
	}

	// * return success response
	return c.JSON(response.Success(c, &payload.ExamQuestionSubmitResponse{
		Submission: &payload.ExamSubmission{
			Id:                submission.Id,
			ExamQuestionId:    submission.ExamQuestionId,
			ExamAttemptId:     submission.ExamAttemptId,
			Answer:            submission.Answer,
			CheckQueryPassed:  submission.CheckQueryPassed,
			CheckQueryAt:      submission.CheckQueryAt,
			CheckPromptPassed: submission.CheckPromptPassed,
			CheckPromptAt:     submission.CheckPromptAt,
			CreatedAt:         submission.CreatedAt,
			UpdatedAt:         submission.UpdatedAt,
		},
	}))
}
