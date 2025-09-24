package submissionProcedure

import (
	"context"
	"time"

	"github.com/bsthun/gut"
)

func (r *Service) ServeOpened(ctx context.Context, userId *uint64, examAttemptId *uint64) *gut.ErrorInstance {
	// * get exam attempt with exam details
	attemptDetails, err := r.database.P().ExamAttemptGetById(ctx, examAttemptId)
	if err != nil {
		return gut.Err(false, "exam attempt not found", err)
	}

	// * verify user attempt ownership
	if *attemptDetails.User.Id != *userId {
		return gut.Err(false, "access denied", nil)
	}

	// * check if attempt is already finished
	if attemptDetails.ExamAttempt.FinishedAt != nil {
		return gut.Err(false, "exam attempt is already finished", nil)
	}

	// * check if current time is between exam opened_at and closed_at
	now := time.Now()
	if now.Before(*attemptDetails.Exam.OpenedAt) {
		return gut.Err(false, "exam is not opened yet", nil)
	}
	if now.After(*attemptDetails.Exam.ClosedAt) {
		return gut.Err(false, "exam is already closed", nil)
	}

	// * check if exam attempt has been started
	if attemptDetails.ExamAttempt.StartedAt == nil {
		return gut.Err(false, "exam attempt has not been started", nil)
	}

	// * update updated_at to track user activity
	_, err = r.database.P().ExamAttemptUpdateUpdatedAt(ctx, examAttemptId)
	if err != nil {
		return gut.Err(false, "failed to update exam attempt activity", err)
	}

	return nil
}
