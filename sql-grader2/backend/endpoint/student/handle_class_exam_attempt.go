package studentEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"fmt"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleClassExamAttempt(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ClassExamAttemptRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate body
	if err := gut.Validate(body); err != nil {
		return err
	}

	// * get exam
	exam, err := r.database.P().ExamGetById(c.Context(), body.ExamId)
	if err != nil {
		return gut.Err(false, "exam not found", err)
	}

	// * check access code
	if exam.AccessCode == nil || *exam.AccessCode != *body.AccessCode {
		return gut.Err(false, "invalid access code", nil)
	}

	// * check if user is member of the class
	joinee, err := r.database.P().ClassJoineeGetByUserAndClass(c.Context(), &psql.ClassJoineeGetByUserAndClassParams{
		UserId:  u.UserId,
		ClassId: exam.ClassId,
	})
	if err != nil {
		return gut.Err(false, "access denied or class not found", err)
	}

	// * check if user already has an attempt for this exam
	_, err = r.database.P().ExamAttemptGetByExamAndJoinee(c.Context(), &psql.ExamAttemptGetByExamAndJoineeParams{
		ExamId:        body.ExamId,
		ClassJoineeId: joinee.Id,
	})
	if err == nil {
		return gut.Err(false, "user already enrolled in this exam", nil)
	}

	// * create database name for this attempt
	databaseName := fmt.Sprintf("exam_%s_%s", gut.Base62(*body.ExamId), gut.Base62(*joinee.Id))

	// * create exam attempt
	_, err = r.database.P().ExamAttemptCreate(c.Context(), &psql.ExamAttemptCreateParams{
		ExamId:        body.ExamId,
		ClassJoineeId: joinee.Id,
		DatabaseName:  &databaseName,
	})
	if err != nil {
		return gut.Err(false, "failed to create exam attempt", err)
	}

	// * return success response
	return c.JSON(response.Success(c, "successfully enrolled in exam"))
}
