package studentEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"time"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleClassExamList(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ClassExamListRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate body
	if err := gut.Validate(body); err != nil {
		return err
	}

	// * check if user is member of the class
	_, err := r.database.P().ClassJoineeGetByUserAndClass(c.Context(), &psql.ClassJoineeGetByUserAndClassParams{
		UserId:  u.UserId,
		ClassId: body.ClassId,
	})
	if err != nil {
		return gut.Err(false, "access denied or class not found", err)
	}

	// * get exams for the class
	examRows, err := r.database.P().ClassExamList(c.Context(), &psql.ClassExamListParams{
		ClassId: body.ClassId,
		UserId:  u.UserId,
	})
	if err != nil {
		return gut.Err(false, "failed to get exams", err)
	}

	// * map exams to items
	examItems, _ := gut.Iterate(examRows, func(examRow psql.ClassExamListRow) (*payload.ClassExamListItem, *gut.ErrorInstance) {
		item := &payload.ClassExamListItem{
			Exam: &payload.Exam{
				Id:           examRow.Exam.Id,
				ClassId:      examRow.Exam.ClassId,
				CollectionId: examRow.Exam.CollectionId,
				Name:         examRow.Exam.Name,
				AccessCode:   nil,
				OpenedAt:     examRow.Exam.OpenedAt,
				ClosedAt:     examRow.Exam.ClosedAt,
				CreatedAt:    examRow.Exam.CreatedAt,
				UpdatedAt:    examRow.Exam.UpdatedAt,
			},
			QuestionCount: examRow.ExamQuestionCount,
		}

		if examRow.ExamAttempt.Id != nil {
			item.ExamAttempt = &payload.ExamAttempt{
				Id:            examRow.ExamAttempt.Id,
				ExamId:        examRow.ExamAttempt.ExamId,
				ClassJoineeId: examRow.ExamAttempt.ClassJoineeId,
				StartedAt:     examRow.ExamAttempt.StartedAt,
				FinishedAt:    examRow.ExamAttempt.FinishedAt,
				CreatedAt:     examRow.ExamAttempt.CreatedAt,
				UpdatedAt:     examRow.ExamAttempt.UpdatedAt,
			}
		}

		now := time.Now()
		var status string

		if now.Before(*examRow.Exam.OpenedAt) {
			status = "upcoming"
		} else if now.After(*examRow.Exam.ClosedAt) {
			status = "closed"
		} else if examRow.ExamAttempt.Id == nil {
			status = "opened"
		} else if examRow.ExamAttempt.FinishedAt != nil {
			status = "finished"
		} else if examRow.ExamAttempt.Id != nil {
			status = "attempted"
		} else {
			status = "unknown"
		}
		item.Status = &status

		return item, nil
	})

	// * response
	return c.JSON(response.Success(c, &payload.ClassExamListResponse{
		Exams: examItems,
	}))
}
