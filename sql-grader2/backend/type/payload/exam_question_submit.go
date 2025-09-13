package payload

type ExamQuestionSubmitRequest struct {
	ExamAttemptId  *uint64 `json:"examAttemptId" validate:"required"`
	ExamQuestionId *uint64 `json:"examQuestionId" validate:"required"`
	Answer         *string `json:"answer" validate:"required"`
}

type ExamQuestionSubmitResponse struct {
	Submission *ExamSubmission `json:"submission"`
}
