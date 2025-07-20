package middleware

import (
	"backend/common/config"
	"backend/type/common"
)

type Middleware struct {
	config   *config.Config
	database common.Database
}

func Init(config *config.Config, database common.Database) *Middleware {
	return &Middleware{
		config:   config,
		database: database,
	}
}
