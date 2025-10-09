package payload

type ExamQuestionSubmitRequest struct {
	ExamAttemptId  *uint64 `json:"examAttemptId" validate:"required"`
	ExamQuestionId *uint64 `json:"examQuestionId" validate:"required"`
	Answer         *string `json:"answer" validate:"required"`
}

type ExamQuestionSubmitResponse struct {
	Submission *ExamSubmission `json:"submission"`
}

type ExamQuestionSubmissionCompare struct {
	Submitted *ExamQuestionSubmissionCompareTable `json:"submitted"`
	Expected  *ExamQuestionSubmissionCompareTable `json:"expected"`
}

type ExamQuestionSubmissionCompareTable struct {
	Header []*string `json:"header"`
	Rows   [][]any   `json:"rows"`
}

type ExamSubmissionCompareRequest struct {
	SubmissionId *uint64 `json:"submissionId" validate:"required"`
}

type ExamSubmissionCompareResponse struct {
	Comparison *ExamQuestionSubmissionCompare `json:"comparison"`
}
