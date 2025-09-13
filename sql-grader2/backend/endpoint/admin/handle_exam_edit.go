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

func (r *Handler) HandleExamEdit(c *fiber.Ctx) error {
	// * parse body
	body := new(payload.ExamEditRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate body
	if err := gut.Validate(body); err != nil {
		return err
	}

	// * get user claims
	_ = c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * update exam
	examId := uint64(body.ExamId)
	exam, err := r.database.P().ExamEdit(c.Context(), &psql.ExamEditParams{
		ExamId:   &examId,
		Name:     body.Name,
		OpenedAt: body.OpenedAt,
		ClosedAt: body.ClosedAt,
	})
	if err != nil {
		return gut.Err(false, "failed to update exam", err)
	}

	// * return success response
	return c.JSON(response.Success(c, exam))
}
