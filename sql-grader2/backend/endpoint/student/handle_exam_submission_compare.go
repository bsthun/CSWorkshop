package studentEndpoint

import (
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleExamSubmissionCompare(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ExamSubmissionCompareRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate body
	if err := gut.Validate(body); err != nil {
		return err
	}

	// * get submission details
	submissionDetail, err := r.database.P().SubmissionDetail(c.Context(), body.SubmissionId)
	if err != nil {
		return gut.Err(false, "submission not found", err)
	}

	// * get exam attempt user details
	attemptDetails, err := r.database.P().ExamAttemptGetById(c.Context(), submissionDetail.ExamAttempt.Id)
	if err != nil {
		return gut.Err(false, "exam attempt not found", err)
	}

	// * verify user owns this submission
	if *attemptDetails.User.Id != *u.UserId {
		return gut.Err(false, "access denied", nil)
	}

	// * check result for execution error
	if submissionDetail.ExamSubmission.Result.ExecutionError != nil {
		return gut.Err(false, "submission has execution error", nil)
	}

	// * check if database name exists for exam attempt
	if attemptDetails.ExamAttempt.DatabaseName == nil {
		return gut.Err(false, "no database configured for exam", nil)
	}

	// * execute and compare queries
	compareResult, er := r.submissionProcedure.ServeCompareQuery(
		c.Context(),
		*submissionDetail.ExamSubmission.Answer,
		*submissionDetail.ExamQuestion.CheckQuery,
		*attemptDetails.ExamAttempt.DatabaseName,
	)
	if er != nil {
		return er
	}

	// * return success response
	return c.JSON(response.Success(c, &payload.ExamSubmissionCompareResponse{
		Comparison: compareResult,
	}))
}
