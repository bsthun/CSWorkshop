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

func (r *Handler) HandleExamQuestionAdd(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ExamQuestionAddRequest)
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

	// * get exam
	exam, err := r.database.P().ExamGetById(c.Context(), body.ExamId)
	if err != nil {
		return gut.Err(false, "failed to get exam", err)
	}

	// * get collection question
	collectionQuestion, err := r.database.P().CollectionQuestionGetById(c.Context(), body.CollectionQuestionId)
	if err != nil {
		return gut.Err(false, "failed to get collection question", err)
	}

	// * check if collection_id matches
	if *collectionQuestion.Collection.Id != *exam.CollectionId {
		return gut.Err(false, "collection question does not belong to exam's collection", nil)
	}

	// * get max order_num for this exam
	maxOrderResult, err := r.database.P().ExamQuestionMaxOrderNum(c.Context(), body.ExamId)
	if err != nil {
		return gut.Err(false, "failed to get max order number", err)
	}

	// * create exam question
	examQuestion, err := r.database.P().ExamQuestionAdd(c.Context(), &psql.ExamQuestionAddParams{
		ExamId:             body.ExamId,
		OriginalQuestionId: body.CollectionQuestionId,
		OrderNum:           gut.Ptr(*maxOrderResult + 1),
		Title:              collectionQuestion.CollectionQuestion.Title,
		Description:        collectionQuestion.CollectionQuestion.Description,
		CheckQuery:         collectionQuestion.CollectionQuestion.CheckQuery,
		CheckPrompt:        collectionQuestion.CollectionQuestion.CheckPrompt,
	})
	if err != nil {
		return gut.Err(false, "failed to create exam question", err)
	}

	// * prepare response
	responsePayload := &payload.ExamQuestion{
		Id:                 examQuestion.Id,
		ExamId:             examQuestion.ExamId,
		OriginalQuestionId: examQuestion.OriginalQuestionId,
		OrderNum:           examQuestion.OrderNum,
		Title:              examQuestion.Title,
		Description:        examQuestion.Description,
		CheckQuery:         examQuestion.CheckQuery,
		CheckPrompt:        examQuestion.CheckPrompt,
		CreatedAt:          examQuestion.CreatedAt,
		UpdatedAt:          examQuestion.UpdatedAt,
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}
