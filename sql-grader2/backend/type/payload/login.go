package payload

type OauthRedirect struct {
	RedirectUrl *string `json:"redirectUrl" validate:"required"`
}

type LoginInput struct {
	Username *string `json:"username" validate:"required"`
}

type OauthCallback struct {
	Code *string `json:"code" validate:"required"`
}
