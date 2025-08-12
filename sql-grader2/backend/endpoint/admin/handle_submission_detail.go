package adminEndpoint

import (
	"backend/generate/psql"
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
	body := new(payload.SubmissionDetailRequest)
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
	row, err := r.database.P().SubmissionDetail(c.Context(), body.Id)
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
			OpenedAt:      row.ExamAttempt.OpenedAt,
			StartedAt:     row.ExamAttempt.StartedAt,
			FinishedAt:    row.ExamAttempt.FinishedAt,
			CreatedAt:     row.ExamAttempt.CreatedAt,
			UpdatedAt:     row.ExamAttempt.UpdatedAt,
		},
		Student: &payload.User{
			Id:         row.User.Id,
			Oid:        row.User.Oid,
			Firstname:  row.User.Firstname,
			Lastname:   row.User.Lastname,
			Email:      row.User.Email,
			PictureUrl: row.User.PictureUrl,
			IsAdmin:    row.User.IsAdmin,
			CreatedAt:  row.User.CreatedAt,
			UpdatedAt:  row.User.UpdatedAt,
		},
		Exam: &payload.Exam{
			Id:           row.Exam.Id,
			ClassId:      row.Exam.ClassId,
			CollectionId: row.Exam.CollectionId,
			Name:         row.Exam.Name,
			AccessCode:   row.Exam.AccessCode,
			OpenedAt:     row.Exam.OpenedAt,
			ClosedAt:     row.Exam.ClosedAt,
			CreatedAt:    row.Exam.CreatedAt,
			UpdatedAt:    row.Exam.UpdatedAt,
		},
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}

func (r *Handler) HandleSubmissionList(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * get request body
	var req payload.SubmissionListRequest
	if err := c.BodyParser(&req); err != nil {
		return gut.Err(false, "failed to parse request body", err)
	}

	// * validate request
	if req.ExamAttemptId == nil && req.ExamQuestionId == nil {
		return gut.Err(false, "either examAttemptId or examQuestionId is required", nil)
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

	// * get submission list
	rows, err := r.database.P().SubmissionList(c.Context(), &psql.SubmissionListParams{
		ExamAttemptId:  req.ExamAttemptId,
		ExamQuestionId: req.ExamQuestionId,
	})
	if err != nil {
		return gut.Err(false, "failed to get submissions", err)
	}

	// * map submissions to payload
	submissions, er := gut.Iterate(rows, func(row psql.SubmissionListRow) (*payload.SubmissionListItem, *gut.ErrorInstance) {
		return &payload.SubmissionListItem{
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
				OpenedAt:      row.ExamAttempt.OpenedAt,
				StartedAt:     row.ExamAttempt.StartedAt,
				FinishedAt:    row.ExamAttempt.FinishedAt,
				CreatedAt:     row.ExamAttempt.CreatedAt,
				UpdatedAt:     row.ExamAttempt.UpdatedAt,
			},
			Student: &payload.User{
				Id:         row.User.Id,
				Oid:        row.User.Oid,
				Firstname:  row.User.Firstname,
				Lastname:   row.User.Lastname,
				Email:      row.User.Email,
				PictureUrl: row.User.PictureUrl,
				IsAdmin:    row.User.IsAdmin,
				CreatedAt:  row.User.CreatedAt,
				UpdatedAt:  row.User.UpdatedAt,
			},
		}, nil
	})
	if er != nil {
		return er
	}

	// * prepare response
	responsePayload := &payload.SubmissionListResponse{
		Submissions: submissions,
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}
