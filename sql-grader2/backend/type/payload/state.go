package payload

type StateResponse struct {
	UserId      *uint64 `json:"userId"`
	DisplayName *string `json:"displayName"`
	Email       *string `json:"email"`
	PictureUrl  *string `json:"pictureUrl"`
	IsAdmin     *bool   `json:"isAdmin"`
}
