package response

import (
	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
)

type SuccessResponse struct {
	Success *bool   `json:"success"`
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
	Data    any     `json:"data,omitempty"`
}

type GenericResponse[T any] struct {
	Success *bool   `json:"success"`
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
	Data    T       `json:"data,omitempty"`
}

func Success(c *fiber.Ctx, args1 any, args2 ...any) *SuccessResponse {
	if message, ok := args1.(string); ok {
		if len(args2) == 0 {
			return &SuccessResponse{
				Success: gut.Ptr(true),
				Message: &message,
			}
		}
		if message2, ok := args2[0].(string); ok {
			return &SuccessResponse{
				Success: gut.Ptr(true),
				Code:    &message,
				Message: &message2,
			}
		} else {
			return &SuccessResponse{
				Success: gut.Ptr(true),
				Code:    &message,
				Data:    &message2,
			}
		}
	}
	return &SuccessResponse{
		Success: gut.Ptr(true),
		Data:    args1,
	}
}
