package tuple

type ExamSubmissionResult struct {
	ExecutionError    string `json:"executionError"`
	PromptError       string `json:"promptError"`
	PromptDescription string `json:"promptDescription"`
	Rows              []any  `json:"rows"`
}
