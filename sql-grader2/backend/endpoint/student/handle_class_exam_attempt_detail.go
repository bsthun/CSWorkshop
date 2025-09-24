package studentEndpoint

import (
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"regexp"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleClassExamAttemptDetail(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ExamAttemptIdRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate body
	if err := gut.Validate(body); err != nil {
		return err
	}

	// * get exam attempt with related data
	attemptData, err := r.database.P().ExamAttemptGetById(c.Context(), body.ExamAttemptId)
	if err != nil {
		return gut.Err(false, "exam attempt not found", err)
	}

	// * verify user owns this attempt
	if *attemptData.User.Id != *u.UserId {
		return gut.Err(false, "access denied", nil)
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

	// * update exam attempt started_at if null
	if attemptData.ExamAttempt.StartedAt == nil {
		_, err = r.database.P().ExamAttemptUpdateStartedAt(c.Context(), attemptData.ExamAttempt.Id)
		if err != nil {
			return gut.Err(false, "failed to update attempt started time", err)
		}
	}

	// * get exam question count
	examQuestions, err := r.database.P().ExamQuestionList(c.Context(), attemptData.Exam.Id)
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
	dialect := "MySQL"

	// * response
	return c.JSON(response.Success(c, &payload.ClassExamDetailResponse{
		Class: &payload.Class{
			Id:           attemptData.Class.Id,
			SemesterId:   attemptData.Class.SemesterId,
			Code:         attemptData.Class.Code,
			Name:         attemptData.Class.Name,
			RegisterCode: attemptData.Class.RegisterCode,
			CreatedAt:    attemptData.Class.CreatedAt,
			UpdatedAt:    attemptData.Class.UpdatedAt,
		},
		Exam: &payload.Exam{
			Id:           attemptData.Exam.Id,
			ClassId:      attemptData.Exam.ClassId,
			CollectionId: attemptData.Exam.CollectionId,
			Name:         attemptData.Exam.Name,
			AccessCode:   attemptData.Exam.AccessCode,
			OpenedAt:     attemptData.Exam.OpenedAt,
			ClosedAt:     attemptData.Exam.ClosedAt,
			CreatedAt:    attemptData.Exam.CreatedAt,
			UpdatedAt:    attemptData.Exam.UpdatedAt,
		},
		ExamQuestionCount: &examQuestionCount,
		Credential: &payload.ExamCredential{
			Dialect:      &dialect,
			Host:         &host,
			Port:         &port,
			User:         user.Metadata.Credential.Username,
			Password:     user.Metadata.Credential.Password,
			DatabaseName: attemptData.ExamAttempt.DatabaseName,
		},
	}))
}
