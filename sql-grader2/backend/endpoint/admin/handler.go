package adminEndpoint

import (
	entityProcedure "backend/procedure/entity"
	"backend/type/common"
)

type Handler struct {
	database common.Database
	entity   *entityProcedure.Service
}

func Handle(database common.Database, entity *entityProcedure.Service) *Handler {
	return &Handler{
		database: database,
		entity:   entity,
	}
}
