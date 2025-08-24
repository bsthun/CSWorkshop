package main

import (
	"backend/common/config"
	"backend/common/database"
	"backend/common/fiber"
	"backend/common/fiber/middleware"
	"backend/common/gorm"
	"backend/common/openai"
	"backend/endpoint"
	adminEndpoint "backend/endpoint/admin"
	publicEndpoint "backend/endpoint/public"
	stateEndpoint "backend/endpoint/state"
	studentEndpoint "backend/endpoint/student"
	entityProcedure "backend/procedure/entity"
	schemaProcedure "backend/procedure/schema"
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
			openai.Init,
			entityProcedure.Serve,
			schemaProcedure.Serve,
			publicEndpoint.Handle,
			stateEndpoint.Handle,
			adminEndpoint.Handle,
			studentEndpoint.Handle,
		),
		fx.Invoke(
			endpoint.Bind,
		),
	).Run()
}
