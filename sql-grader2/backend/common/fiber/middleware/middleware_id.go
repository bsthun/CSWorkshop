package middleware

import (
	"encoding/json"
	"github.com/bsthun/gut"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (r *Middleware) Id() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if len(c.Response().Body()) > 0 {
				var responseBody map[string]any
				if err := json.Unmarshal(c.Response().Body(), &responseBody); err != nil {
					return
				}

				IdProcessResponseId(responseBody)

				modifiedResponse, err := json.Marshal(responseBody)
				if err != nil {
					return
				}
				c.Response().SetBody(modifiedResponse)
			}
		}()

		if !strings.Contains(c.Get("Content-Type"), "application/json") {
			return c.Next()
		}

		var requestBody map[string]any

		if len(c.Body()) > 0 {
			if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
				return gut.Err(false, "unable to parse body", err)
			}

			IdProcessRequestPayload(requestBody)

			modifiedBody, err := json.Marshal(requestBody)
			if err != nil {
				return gut.Err(false, "failed to marshal modified request", err)
			}
			c.Request().SetBody(modifiedBody)
		}

		return c.Next()
	}
}

func IdProcessRequestPayload(payload any) {
	if payload == nil {
		return
	}

	switch v := payload.(type) {
	case map[string]any:
		for key, value := range v {
			if IdField(key) {
				if strId, ok := value.(string); ok {
					decodedId, err := gut.Decode(strId)
					if err == nil {
						v[key] = decodedId
					}
				}
			} else {
				IdProcessRequestPayload(value)
			}
		}
	case []any:
		for _, item := range v {
			IdProcessRequestPayload(item)
		}
	}
}

func IdProcessResponseId(data any) {
	if data == nil {
		return
	}

	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			if IdField(key) {
				switch idVal := value.(type) {
				case float64:
					uint64Id := uint64(idVal)
					encodedId := gut.EncodeId(uint64Id)
					v[key] = encodedId
				case uint64:
					encodedId := gut.EncodeId(idVal)
					v[key] = encodedId
				case json.Number:
					if numValue, err := idVal.Int64(); err == nil {
						encodedId := gut.EncodeId(uint64(numValue))
						v[key] = encodedId
					}
				case nil:
				default:
					rv := reflect.ValueOf(idVal)
					if rv.Kind() == reflect.Ptr && !rv.IsNil() {
						if rv.Elem().Kind() == reflect.Uint64 {
							encodedId := gut.EncodeId(rv.Elem().Uint())
							v[key] = encodedId
						}
					}
				}
			} else {
				IdProcessResponseId(value)
			}
		}
	case []any:
		for _, item := range v {
			IdProcessResponseId(item)
		}
	}
}

func IdField(fieldName string) bool {
	return fieldName == "id" || strings.HasSuffix(fieldName, "Id")
}
