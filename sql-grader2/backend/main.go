package main

import (
	"backend/common/config"
	"backend/common/database"
	"backend/common/fiber"
	"backend/common/fiber/middleware"
	"backend/common/gorm"
	"backend/endpoint"
	adminEndpoint "backend/endpoint/admin"
	publicEndpoint "backend/endpoint/public"
	stateEndpoint "backend/endpoint/state"
	entityProcedure "backend/procedure/entity"
	"backend/type/common"
	"embed"

	"go.uber.org/fx"
)

//go:embed database/postgres/migration/*.sql
var migration embed.FS

//go:embed .local/dist/*
var frontend embed.FS

func main() {
	fx.New(
		fx.Provide(
			func() common.MigrationFS {
				return migration
			},
			func() common.FrontendFS {
				return frontend
			},
			config.Init,
			database.Init,
			gorm.Init,
			fiber.Init,
			middleware.Init,
			entityProcedure.Serve,
			publicEndpoint.Handle,
			stateEndpoint.Handle,
			adminEndpoint.Handle,
		),
		fx.Invoke(
			endpoint.Bind,
		),
	).Run()
}
