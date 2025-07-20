package adminEndpoint

import (
	"backend/type/common"
)

type Handler struct {
	database common.Database
}

func Handle(database common.Database) *Handler {
	return &Handler{
		database: database,
	}
}
