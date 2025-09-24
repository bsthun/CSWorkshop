package studentEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

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

	// * check if exam is open
	now := time.Now()
	if now.Before(*exam.OpenedAt) {
		return gut.Err(false, "exam is not opened yet", nil)
	}
	if now.After(*exam.ClosedAt) {
		return gut.Err(false, "exam is already closed", nil)
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
	if !errors.Is(err, sql.ErrNoRows) {
		return gut.Err(false, "unable to check existing exam attempt", err)
	}

	// * generate unique database name for this attempt
	var databaseName string
	for {
		databaseName = fmt.Sprintf("exam_%s", *gut.Random(gut.RandomSet.Num, 12))
		_, err := r.database.P().ExamAttemptGetByDatabaseName(c.Context(), &databaseName)
		if errors.Is(err, sql.ErrNoRows) {
			break
		}
		if err != nil {
			return gut.Err(false, "failed to check database name uniqueness", err)
		}
	}

	// * create exam attempt
	attempt, err := r.database.P().ExamAttemptCreate(c.Context(), &psql.ExamAttemptCreateParams{
		ExamId:        body.ExamId,
		ClassJoineeId: joinee.Id,
		DatabaseName:  &databaseName,
	})
	if err != nil {
		return gut.Err(false, "failed to create exam attempt", err)
	}

	// * get user details with metadata
	user, err := r.database.P().UserGetById(c.Context(), u.UserId)
	if err != nil {
		return gut.Err(false, "failed to get user", err)
	}

	// * check if user has database credential
	if user.Metadata == nil || user.Metadata.Credential == nil || user.Metadata.Credential.Username == nil || user.Metadata.Credential.Password == nil {
		return gut.Err(false, "user database credential not found", nil)
	}

	// * get collection for schema
	_, err = r.database.P().CollectionGetById(c.Context(), exam.CollectionId)
	if err != nil {
		return gut.Err(false, "failed to get collection", err)
	}

	// * read schema file
	schemaPath := fmt.Sprintf(".local/collection/%s/schema.sql", gut.Base62(*exam.CollectionId))
	schemaContent, err := os.ReadFile(schemaPath)
	if err != nil {
		return gut.Err(false, "failed to read schema file", err)
	}

	// * migrate database
	er := r.schemaProcedure.ServeDatabaseMigrate(c.Context(), *r.config.MysqlDsn, databaseName, string(schemaContent))
	if er != nil {
		return er
	}

	// * grant select access to user
	grantSQL := fmt.Sprintf("GRANT SELECT ON `%s`.* TO '%s'@'%%'", databaseName, *user.Metadata.Credential.Username)
	tx := r.gorm.Exec(grantSQL)
	if tx.Error != nil {
		return gut.Err(false, "failed to grant database access", tx.Error)
	}

	// * flush privileges
	tx = r.gorm.Exec("FLUSH PRIVILEGES")
	if tx.Error != nil {
		return gut.Err(false, "failed to flush privileges", tx.Error)
	}

	// * response
	return c.JSON(response.Success(c, &payload.ClassExamAttemptResponse{
		ExamAttempt: &payload.ExamAttempt{
			Id:            attempt.Id,
			ExamId:        attempt.ExamId,
			ClassJoineeId: attempt.ClassJoineeId,
			StartedAt:     attempt.StartedAt,
			FinishedAt:    attempt.FinishedAt,
			CreatedAt:     attempt.CreatedAt,
			UpdatedAt:     attempt.UpdatedAt,
		},
	}))
}
