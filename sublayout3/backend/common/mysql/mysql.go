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
	"backend/type/table"
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
	if err := cc.DB.AutoMigrate(
		new(table.SystemGroup),
		new(table.SystemComponent),
		new(table.C01Sidebar),
		new(table.C02Recent),
		new(table.C03Track),
		new(table.C04Artist),
		new(table.C05Presearch),
		new(table.C06Search),
		new(table.C07Track),
		new(table.C07Artist),
		new(table.C08Track),
		new(table.C08Album),
		new(table.C10Concert),
		new(table.C12Track),
		new(table.C11Detail),
		new(table.C12Playlist),
		new(table.C13TrackDetail),
		new(table.C14PodcastSection),
		new(table.C14PodcastGenre),
		new(table.C16PodcastShow),
		new(table.C17PodcastEpisode),
		new(table.C18PlaylistSection),
		new(table.C18PlaylistItem),
		new(table.C19Album),
		new(table.C19Track),
	); err != nil {
		return err
	}

	return nil
}
