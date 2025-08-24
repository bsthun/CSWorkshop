package schemaProcedure

import (
	"backend/common/config"
	"backend/type/common"
	"context"

	"github.com/bsthun/gut"
	"gorm.io/gorm"
)

type Server interface {
	ServeDatabaseMigrate(ctx context.Context, mysqlDsn string, databaseName string, schemaContent string) *gut.ErrorInstance
}

type Service struct {
	database common.Database
	gorm     *gorm.DB
	config   *config.Config
}

func Serve(database common.Database, gorm *gorm.DB, config *config.Config) Server {
	return &Service{
		database: database,
		gorm:     gorm,
		config:   config,
	}
}