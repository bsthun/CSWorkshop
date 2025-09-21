package studentEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleStudentExamQuestionDetail(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ExamQuestionDetailRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate body
	if err := gut.Validate(body); err != nil {
		return err
	}

	// * get exam attempt
	attemptDetails, err := r.database.P().ExamAttemptGetById(c.Context(), body.ExamAttemptId)
	if err != nil {
		return gut.Err(false, "exam attempt not found", err)
	}

	// * verify user attempt ownership
	if *attemptDetails.User.Id != *u.UserId {
		return gut.Err(false, "access denied", nil)
	}

	// * get exam question with submissions
	results, err := r.database.P().ExamQuestionDetailWithSubmissions(c.Context(), &psql.ExamQuestionDetailWithSubmissionsParams{
		Id:            body.ExamQuestionId,
		ExamAttemptId: body.ExamAttemptId,
	})
	if err != nil {
		return gut.Err(false, "failed to get question details", err)
	}

	if len(results) == 0 {
		return gut.Err(false, "exam question not found", nil)
	}

	// * verify question belongs to the exam
	if *results[0].ExamQuestion.ExamId != *attemptDetails.Exam.Id {
		return gut.Err(false, "question does not belong to this exam", nil)
	}

	// * map results to response
	examQuestion := &payload.ExamQuestion{
		Id:                 results[0].ExamQuestion.Id,
		ExamId:             results[0].ExamQuestion.ExamId,
		OriginalQuestionId: results[0].ExamQuestion.OriginalQuestionId,
		OrderNum:           results[0].ExamQuestion.OrderNum,
		Title:              results[0].ExamQuestion.Title,
		Description:        results[0].ExamQuestion.Description,
		CheckQuery:         results[0].ExamQuestion.CheckQuery,
		CheckPrompt:        results[0].ExamQuestion.CheckPrompt,
		CreatedAt:          results[0].ExamQuestion.CreatedAt,
		UpdatedAt:          results[0].ExamQuestion.UpdatedAt,
	}

	submissions, _ := gut.Iterate(results, func(result psql.ExamQuestionDetailWithSubmissionsRow) (*payload.ExamSubmissionWithStatus, *gut.ErrorInstance) {
		if result.ExamSubmission.Id == nil {
			return &payload.ExamSubmissionWithStatus{
				ExamSubmission: nil,
				Status:         result.Status,
			}, nil
		}

		return &payload.ExamSubmissionWithStatus{
			ExamSubmission: &payload.ExamSubmission{
				Id:                result.ExamSubmission.Id,
				ExamQuestionId:    result.ExamSubmission.ExamQuestionId,
				ExamAttemptId:     result.ExamSubmission.ExamAttemptId,
				Answer:            result.ExamSubmission.Answer,
				Result:            result.ExamSubmission.Result,
				CheckQueryPassed:  result.ExamSubmission.CheckQueryPassed,
				CheckQueryAt:      result.ExamSubmission.CheckQueryAt,
				CheckPromptPassed: result.ExamSubmission.CheckPromptPassed,
				CheckPromptAt:     result.ExamSubmission.CheckPromptAt,
				CreatedAt:         result.ExamSubmission.CreatedAt,
				UpdatedAt:         result.ExamSubmission.UpdatedAt,
			},
			Status: result.Status,
		}, nil
	})

	// * return success response
	return c.JSON(response.Success(c, &payload.StudentExamQuestionDetailResponse{
		ExamQuestion: examQuestion,
		Submissions:  submissions,
	}))
}
