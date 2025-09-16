package tuple

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ExamSubmissionResult struct {
	ExecutionError    string `json:"executionError"`
	PromptError       string `json:"promptError"`
	PromptDescription string `json:"promptDescription"`
	Rows              []any  `json:"rows"`
}

func (r *ExamSubmissionResult) Scan(value any) error {
	if value == nil {
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(data, r)
}

func (r *ExamSubmissionResult) Value() (driver.Value, error) {
	return json.Marshal(r)
}
