package gorm

import (
	"backend/common/config"
	"backend/util/orm"
	"log"
	"os"
	"time"

	"github.com/bsthun/gut"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(config *config.Config) *gorm.DB {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             86400 * time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	conn, er := orm.Connect(*config.MysqlDsn, "mysql")
	if er != nil {
		gut.Fatal("unable to create mysql connection", err)
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
		gut.Fatal("unable to open mysql connection", err)
	}

	return db
}
