package mysql

import "C"
import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"backend/common"
	ilog "backend/util/log"
)

func Init() {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             100 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open SQL connection
	dialector := mysql.New(
		mysql.Config{
			DSN:               *cc.Config.MySqlDsn,
			DefaultStringSize: 255,
		},
	)

	// * Open main database connection
	if db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
	}); err != nil {
		ilog.Fatal("Unable to connect to MySQL database", err)
	} else {
		cc.DB = db
	}

	// Initialize model migrations
	if *cc.Config.MySqlMigrate {
		if err := migrate(); err != nil {
			ilog.Fatal("Unable to migrate GORM model", err)
		}
	}

	ilog.Debug("MySQL database connected")
}

func migrate() error {
	// * Migrate model
	if err := cc.DB.AutoMigrate(); err != nil {
		return err
	}

	return nil
}
