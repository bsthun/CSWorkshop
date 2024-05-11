package response

type ErrorInstance struct {
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
	Err     error  `json:"error,omitempty"`
}

func (v *ErrorInstance) Error() string {
	return v.Message
}

func Error(critical bool, message string, args2 ...any) *ErrorInstance {
	if len(args2) == 1 {
		if code, ok := args2[0].(string); ok {
			return &ErrorInstance{
				Message: message,
				Code:    code,
				Err:     nil,
			}
		}
		if err, ok := args2[0].(error); ok {
			return &ErrorInstance{
				Message: message,
				Code:    "",
				Err:     err,
			}
		}
	}

	if len(args2) == 2 {
		if code, ok := args2[0].(string); ok {
			if err, ok := args2[1].(error); ok {
				return &ErrorInstance{
					Message: message,
					Code:    code,
					Err:     err,
				}
			}
		}
	}

	return &ErrorInstance{
		Message: message,
		Err:     nil,
	}
}
