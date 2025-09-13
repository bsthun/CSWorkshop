package payload

import "time"

type ExamEditRequest struct {
	ExamId   int64      `json:"examId" validate:"required"`
	Name     *string    `json:"name"`
	OpenedAt *time.Time `json:"openedAt"`
	ClosedAt *time.Time `json:"closedAt"`
}
