package studentEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"backend/type/tuple"
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

	// * check if exam is opened
	if er := r.submissionProcedure.ServeOpened(c.Context(), u.UserId, body.ExamAttemptId); er != nil {
		return er
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

	// * prepare result object
	result := &tuple.ExamSubmissionResult{
		ExecutionError:    "",
		PromptError:       "",
		PromptDescription: "",
	}

	// * prepare timing and status variables
	now := time.Now()
	checkQueryPassed := false
	checkPromptPassed := false
	var checkPromptAt *time.Time

	// * check query execution and comparison
	if attemptDetails.ExamAttempt.DatabaseName != nil {
		queryPassed, er := r.submissionProcedure.ServeCheckQuery(c.Context(), *body.Answer, *examQuestion.CheckQuery, *attemptDetails.ExamAttempt.DatabaseName)
		if er != nil {
			// * query execution failed
			result.ExecutionError = er.Errors[0].Err.Error()

			// * create submission with execution error
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

		// * query executed successfully
		checkQueryPassed = queryPassed

		// * if query passed, check prompt requirements
		if checkQueryPassed {
			promptPassed, promptDescription, er := r.submissionProcedure.ServeCheckPrompt(c.Context(), *body.Answer, *examQuestion.CheckPrompt)
			if er != nil {
				// * prompt check failed due to error, but query passed
				result.PromptError = er.Error()
				checkPromptPassed = false
			} else {
				// * prompt check completed
				checkPromptPassed = promptPassed
				result.PromptDescription = promptDescription
				checkPromptAt = &now
			}
		}
	}

	// * create submission with final results
	var checkPromptPassedPtr *bool
	if checkPromptAt != nil {
		checkPromptPassedPtr = &checkPromptPassed
	}

	submission, err := r.database.P().ExamSubmissionCreate(c.Context(), &psql.ExamSubmissionCreateParams{
		ExamQuestionId:    body.ExamQuestionId,
		ExamAttemptId:     body.ExamAttemptId,
		Answer:            body.Answer,
		Result:            result,
		CheckQueryPassed:  &checkQueryPassed,
		CheckQueryAt:      &now,
		CheckPromptPassed: checkPromptPassedPtr,
		CheckPromptAt:     checkPromptAt,
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
