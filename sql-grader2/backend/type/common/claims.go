package common

import (
	"encoding/json"
	"github.com/bsthun/gut"
	"time"
)

type OidcClaims struct {
	Id        *string `json:"sub"`
	FirstName *string `json:"given_name"`
	Lastname  *string `json:"family_name"`
	Picture   *string `json:"picture"`
	Email     *string `json:"email"`
}

type LoginClaims struct {
	UserId    *uint64    `json:"userId"`
	ExpiredAt *time.Time `json:"exp"`
}

func (r *LoginClaims) Valid() error {
	return nil
}

func (r *LoginClaims) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"userId": gut.EncodeId(*r.UserId),
	})
}

func (r *LoginClaims) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	userId, err := gut.Decode(raw["userId"].(string))
	if err != nil {
		return err
	}
	r.UserId = &userId
	return nil
}
