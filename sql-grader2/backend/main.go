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
	"embed"

	"go.uber.org/fx"
)

//go:embed database/postgres/migration/*.sql
var embedMigrations embed.FS

func main() {
	fx.New(
		fx.Supply(
			embedMigrations,
		),
		fx.Provide(
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
