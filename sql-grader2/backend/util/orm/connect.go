package orm

import (
	"database/sql"
	"strings"
	"time"

	"github.com/bsthun/gut"
)

func Connect(mysqlDsn string, databaseName string) (*sql.DB, *gut.ErrorInstance) {
	dsn := strings.Replace(mysqlDsn, "{{database}}", databaseName, 1)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, gut.Err(false, "failed to connect to database", err)
	}

	conn.SetMaxIdleConns(10)
	conn.SetMaxOpenConns(100)
	conn.SetConnMaxLifetime(time.Hour)

	return conn, nil
}
