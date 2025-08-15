package database

import (
	"backend/common/config"
	"backend/generate/psql"
	"backend/type/common"
	"database/sql"
	"strings"

	"github.com/bsthun/gut"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func Init(config *config.Config, migration common.MigrationFS) common.Database {
	// * initialize postgres database
	postgres, err := sql.Open("postgres", *config.PostgresDsn)
	if err != nil {
		gut.Fatal("unable to connect to postgres database", err)
	}

	if err = postgres.Ping(); err != nil {
		gut.Fatal("unable to ping database", err)
	}

	// * wrap the database with query logging
	postgresQuerier := psql.New(&Wrapper{
		db: postgres,
	})

	// * construct struct
	database := &Database{
		PQuerier: postgresQuerier,
		PConn:    postgres,
	}

	// * run migrations
	goose.SetBaseFS(migration)
	goose.SetTableName("_gooses")
	if err := goose.Up(postgres, "database/postgres/migration"); err != nil {
		if !strings.HasSuffix(err.Error(), "directory does not exist") {
			gut.Fatal("failed to run migrations", err)
		}
	}

	return database
}
