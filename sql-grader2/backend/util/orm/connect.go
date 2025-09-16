package orm

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bsthun/gut"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

func Instance(mysqlDsn string, databaseName string) (*gorm.DB, *gut.ErrorInstance) {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             86400 * time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	conn, er := Connect(mysqlDsn, databaseName)
	if er != nil {
		return nil, er
	}

	dialector := mysql.New(
		mysql.Config{
			Conn: conn,
		},
	)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, gut.Err(false, "unable to open mysql connection", err)
	}

	return db, nil
}
