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

func (r *Handler) HandleSubmissionList(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.SubmissionListRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate request
	if body.ExamAttemptId == nil && body.ExamQuestionId == nil {
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
		ExamAttemptId:  body.ExamAttemptId,
		ExamQuestionId: body.ExamQuestionId,
	})
	if err != nil {
		return gut.Err(false, "failed to get submissions", err)
	}

	// * get exam
	var exam *payload.Exam
	if len(rows) > 0 {
		examRow, err := r.database.P().ExamGetById(c.Context(), rows[0].ExamAttempt.ExamId)
		if err != nil {
			return gut.Err(false, "failed to get exam", err)
		}
		exam = &payload.Exam{
			Id:           examRow.Id,
			ClassId:      examRow.ClassId,
			CollectionId: examRow.CollectionId,
			Name:         examRow.Name,
			AccessCode:   examRow.AccessCode,
			OpenedAt:     examRow.OpenedAt,
			ClosedAt:     examRow.ClosedAt,
			CreatedAt:    examRow.CreatedAt,
			UpdatedAt:    examRow.UpdatedAt,
		}
	}

	// * map submissions to payload
	submissions, er := gut.Iterate(rows, func(row psql.SubmissionListRow) (*payload.SubmissionListItem, *gut.ErrorInstance) {
		return &payload.SubmissionListItem{
			Submission: &payload.ExamSubmission{
				Id:                row.ExamSubmission.Id,
				ExamQuestionId:    row.ExamSubmission.ExamQuestionId,
				ExamAttemptId:     row.ExamSubmission.ExamAttemptId,
				Answer:            row.ExamSubmission.Answer,
				Result:            row.ExamSubmission.Result,
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
			Joinee: &payload.ClassJoineeInfo{
				Id:        row.ClassJoinee.Id,
				UserId:    row.ClassJoinee.UserId,
				ClassId:   row.ClassJoinee.ClassId,
				CreatedAt: row.ClassJoinee.CreatedAt,
				UpdatedAt: row.ClassJoinee.UpdatedAt,
			},
		}, nil
	})
	if er != nil {
		return er
	}

	// * prepare response
	responsePayload := &payload.SubmissionListResponse{
		Exam:        exam,
		Submissions: submissions,
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}
