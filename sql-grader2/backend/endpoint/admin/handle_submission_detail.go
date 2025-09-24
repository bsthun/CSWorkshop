package adminEndpoint

import (
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleSubmissionDetail(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.SubmissionIdRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * get user from database
	user, err := r.database.P().UserGetById(c.Context(), u.UserId)
	if err != nil {
		return gut.Err(false, "failed to get user", err)
	}

	// * check if user is admin
	if user.IsAdmin == nil || !*user.IsAdmin {
		return gut.Err(false, "access denied", nil)
	}

	// * get submission detail
	row, err := r.database.P().SubmissionDetail(c.Context(), body.SubmissionId)
	if err != nil {
		return gut.Err(false, "failed to get submission", err)
	}

	// * prepare response
	responsePayload := &payload.SubmissionDetailResponse{
		Submission: &payload.ExamSubmission{
			Id:                row.ExamSubmission.Id,
			ExamQuestionId:    row.ExamSubmission.ExamQuestionId,
			ExamAttemptId:     row.ExamSubmission.ExamAttemptId,
			Answer:            row.ExamSubmission.Answer,
			CheckQueryPassed:  row.ExamSubmission.CheckQueryPassed,
			CheckQueryAt:      row.ExamSubmission.CheckQueryAt,
			CheckPromptPassed: row.ExamSubmission.CheckPromptPassed,
			CheckPromptAt:     row.ExamSubmission.CheckPromptAt,
			CreatedAt:         row.ExamSubmission.CreatedAt,
			UpdatedAt:         row.ExamSubmission.UpdatedAt,
		},
		Question: &payload.ExamQuestion{
			Id:                 row.ExamQuestion.Id,
			ExamId:             row.ExamQuestion.ExamId,
			OriginalQuestionId: row.ExamQuestion.OriginalQuestionId,
			OrderNum:           row.ExamQuestion.OrderNum,
			Title:              row.ExamQuestion.Title,
			Description:        row.ExamQuestion.Description,
			CheckQuery:         row.ExamQuestion.CheckQuery,
			CheckPrompt:        row.ExamQuestion.CheckPrompt,
			CreatedAt:          row.ExamQuestion.CreatedAt,
			UpdatedAt:          row.ExamQuestion.UpdatedAt,
		},
		Attempt: &payload.ExamAttempt{
			Id:            row.ExamAttempt.Id,
			ExamId:        row.ExamAttempt.ExamId,
			ClassJoineeId: row.ExamAttempt.ClassJoineeId,
			StartedAt:     row.ExamAttempt.StartedAt,
			FinishedAt:    row.ExamAttempt.FinishedAt,
			CreatedAt:     row.ExamAttempt.CreatedAt,
			UpdatedAt:     row.ExamAttempt.UpdatedAt,
		},
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}
