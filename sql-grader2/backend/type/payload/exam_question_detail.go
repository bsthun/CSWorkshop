package payload

type ExamQuestionDetailRequest struct {
	ExamAttemptId  *uint64 `json:"examAttemptId" validate:"required"`
	ExamQuestionId *uint64 `json:"examQuestionId" validate:"required"`
}

type ExamSubmissionWithStatus struct {
	ExamSubmission *ExamSubmission `json:"examSubmission"`
	Status         *string         `json:"status" validate:"required,oneof=passed rejected invalid unsubmitted"`
}

type StudentExamQuestionDetailResponse struct {
	ExamQuestion *ExamQuestion               `json:"examQuestion"`
	Submissions  []*ExamSubmissionWithStatus `json:"submissions"`
}
