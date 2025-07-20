package common

type Paginate struct {
	Limit  *int32 `json:"limit" validate:"required,gte=1,lte=100"`
	Offset *int32 `json:"offset" validate:"required,gte=0"`
}
