package studentEndpoint

import (
	"backend/common/config"
	entityProcedure "backend/procedure/entity"
	schemaProcedure "backend/procedure/schema"
	"backend/type/common"

	"gorm.io/gorm"
)

type Handler struct {
	config          *config.Config
	database        common.Database
	entity          entityProcedure.Server
	schemaProcedure schemaProcedure.Server
	gorm            *gorm.DB
}

func Handle(config *config.Config, database common.Database, entity entityProcedure.Server, schemaProcedure schemaProcedure.Server, gorm *gorm.DB) *Handler {
	return &Handler{
		config:          config,
		database:        database,
		entity:          entity,
		schemaProcedure: schemaProcedure,
		gorm:            gorm,
	}
}
