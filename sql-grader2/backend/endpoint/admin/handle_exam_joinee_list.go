package adminEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"strconv"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleExamJoineeList(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * get exam id from params
	examIdParam := c.Params("id")
	examId, err := strconv.ParseUint(examIdParam, 10, 64)
	if err != nil {
		return gut.Err(false, "invalid exam id", err)
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

	// * get exam joinee list with user info
	joineeRows, err := r.database.P().ExamJoineeList(c.Context(), &examId)
	if err != nil {
		return gut.Err(false, "failed to get exam joinees", err)
	}

	// * map joinees to payload
	joinees, er := gut.Iterate(joineeRows, func(row psql.ExamJoineeListRow) (*payload.ExamJoineeListItem, *gut.ErrorInstance) {
		return &payload.ExamJoineeListItem{
			Id:            row.ExamAttempt.Id,
			ExamId:        row.ExamAttempt.ExamId,
			ClassJoineeId: row.ExamAttempt.ClassJoineeId,
			OpenedAt:      row.ExamAttempt.OpenedAt,
			StartedAt:     row.ExamAttempt.StartedAt,
			FinishedAt:    row.ExamAttempt.FinishedAt,
			CreatedAt:     row.ExamAttempt.CreatedAt,
			UpdatedAt:     row.ExamAttempt.UpdatedAt,
			Joinee: &payload.ExamJoinee{
				Id: row.ClassJoinee.Id,
				User: &payload.User{
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
			},
			Score: &payload.ExamScore{
				Passed:      row.PassedCount,
				Rejected:    row.RejectedCount,
				Invalid:     row.InvalidCount,
				Unsubmitted: row.UnsubmittedCount,
			},
		}, nil
	})
	if er != nil {
		return er
	}

	// * prepare response
	responsePayload := &payload.ExamJoineeListResponse{
		Joinees: joinees,
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}
