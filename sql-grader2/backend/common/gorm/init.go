package gorm

import (
	"backend/common/config"
	"backend/util/orm"

	"github.com/bsthun/gut"
	"gorm.io/gorm"
)

func Init(config *config.Config) *gorm.DB {
	db, er := orm.Instance(*config.MysqlDsn, "mysql")
	if er != nil {
		gut.Fatal(er.Error(), er.Errors[0].Err)
	}

	return db
}
