package studentEndpoint

import (
	"backend/common/config"
	entityProcedure "backend/procedure/entity"
	"backend/type/common"

	"gorm.io/gorm"
)

type Handler struct {
	config   *config.Config
	database common.Database
	entity   entityProcedure.Server
	gorm     *gorm.DB
}

func Handle(config *config.Config, database common.Database, entity entityProcedure.Server, gorm *gorm.DB) *Handler {
	return &Handler{
		config:   config,
		database: database,
		entity:   entity,
		gorm:     gorm,
	}
}
