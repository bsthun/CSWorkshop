package studentEndpoint

import (
	"backend/generate/psql"
	"backend/procedure/schema"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"backend/type/tuple"
	"fmt"
	"os"
	"regexp"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleClassExamDetail(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ClassExamDetailRequest)
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

	// * check if user is member of the class
	joinee, err := r.database.P().ClassJoineeGetByUserAndClass(c.Context(), &psql.ClassJoineeGetByUserAndClassParams{
		UserId:  u.UserId,
		ClassId: exam.ClassId,
	})
	if err != nil {
		return gut.Err(false, "access denied or class not found", err)
	}

	// * get class details
	classDetail, err := r.database.P().ClassGetById(c.Context(), exam.ClassId)
	if err != nil {
		return gut.Err(false, "failed to get class", err)
	}

	// * get exam attempt
	attempt, err := r.database.P().ExamAttemptGetByExamAndJoinee(c.Context(), &psql.ExamAttemptGetByExamAndJoineeParams{
		ExamId:        body.ExamId,
		ClassJoineeId: joinee.Id,
	})
	if err != nil {
		return gut.Err(false, "exam attempt not found", err)
	}

	// * get user details with metadata
	user, err := r.database.P().UserGetById(c.Context(), u.UserId)
	if err != nil {
		return gut.Err(false, "failed to get user", err)
	}

	// * check if user has mysql username in metadata
	var username, password string
	if user.Metadata != nil && user.Metadata.Credential != nil && user.Metadata.Credential.Username != nil {
		username = *user.Metadata.Credential.Username
		password = *user.Metadata.Credential.Password
	} else {
		// * create mysql username and password
		username = fmt.Sprintf("user_%s", gut.Base62(*u.UserId))
		password = *gut.Random(gut.RandomSet.MixedAlphaNum, 16)

		// * create mysql user
		createUserSQL := fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY '%s'", username, password)
		tx := r.gorm.Exec(createUserSQL)
		if tx.Error != nil {
			return gut.Err(false, "failed to create mysql user", tx.Error)
		}

		// * update user metadata
		metadata := &tuple.UserMetadata{
			Credential: &tuple.UserMetadataCredential{
				Username: &username,
				Password: &password,
			},
		}
		_, err = r.database.P().UserUpdateMetadata(c.Context(), &psql.UserUpdateMetadataParams{
			Id:       u.UserId,
			Metadata: metadata,
		})
		if err != nil {
			return gut.Err(false, "failed to update user metadata", err)
		}
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

	// * create schema procedure service
	schemaProcedure := schemaProcedure.Serve(r.database, r.gorm, r.config)

	// * migrate database
	er := schemaProcedure.ServeDatabaseMigrate(c.Context(), *r.config.MysqlDsn, *attempt.DatabaseName, string(schemaContent))
	if er != nil {
		return er
	}

	// * grant select access to user
	grantSQL := fmt.Sprintf("GRANT SELECT ON `%s`.* TO '%s'@'%%'", *attempt.DatabaseName, username)
	tx := r.gorm.Exec(grantSQL)
	if tx.Error != nil {
		return gut.Err(false, "failed to grant database access", tx.Error)
	}

	// * flush privileges
	tx = r.gorm.Exec("FLUSH PRIVILEGES")
	if tx.Error != nil {
		return gut.Err(false, "failed to flush privileges", tx.Error)
	}

	// * update exam attempt opened_at if null
	if attempt.OpenedAt == nil {
		_, err = r.database.P().ExamAttemptUpdateOpenedAt(c.Context(), attempt.Id)
		if err != nil {
			return gut.Err(false, "failed to update attempt opened time", err)
		}
	}

	// * get exam question count
	examQuestions, err := r.database.P().ExamQuestionList(c.Context(), body.ExamId)
	if err != nil {
		return gut.Err(false, "failed to get exam questions", err)
	}
	examQuestionCount := uint64(len(examQuestions))

	// * parse mysql dsn for host and port
	re := regexp.MustCompile(`^([^:]+):([^@]+)@tcp\(([^:]+):(\d+)\)/`)
	matches := re.FindStringSubmatch(*r.config.MysqlDsn)
	if len(matches) != 5 {
		return gut.Err(false, "invalid mysql dsn format", nil)
	}

	host := matches[3]
	port := matches[4]
	dialect := "mysql"

	// * response
	return c.JSON(response.Success(c, &payload.ClassExamDetailResponse{
		Class: &payload.Class{
			Id:           classDetail.Class.Id,
			SemesterId:   classDetail.Class.SemesterId,
			Code:         classDetail.Class.Code,
			Name:         classDetail.Class.Name,
			RegisterCode: classDetail.Class.RegisterCode,
			CreatedAt:    classDetail.Class.CreatedAt,
			UpdatedAt:    classDetail.Class.UpdatedAt,
		},
		Exam: &payload.Exam{
			Id:           exam.Id,
			ClassId:      exam.ClassId,
			CollectionId: exam.CollectionId,
			Name:         exam.Name,
			AccessCode:   exam.AccessCode,
			OpenedAt:     exam.OpenedAt,
			ClosedAt:     exam.ClosedAt,
			CreatedAt:    exam.CreatedAt,
			UpdatedAt:    exam.UpdatedAt,
		},
		ExamQuestionCount: &examQuestionCount,
		Credential: &payload.ExamCredential{
			Dialect:      &dialect,
			Host:         &host,
			Port:         &port,
			User:         &username,
			Password:     &password,
			DatabaseName: attempt.DatabaseName,
		},
	}))
}