package submissionProcedure

import (
	"backend/util/orm"
	"context"
	"reflect"

	"github.com/bsthun/gut"
)

func (r *Service) ServeCheckQuery(ctx context.Context, studentAnswer string, checkQuery string, databaseName string) (bool, *gut.ErrorInstance) {
	// * connect to exam database
	gorm, er := orm.Instance(*r.config.MysqlDsn, databaseName)
	if er != nil {
		return false, gut.Err(false, "failed to connect to exam database instance", er)
	}

	// * execute student answer query
	var studentResult []map[string]any
	tx := gorm.Raw(studentAnswer).Scan(&studentResult)
	if tx.Error != nil {
		return false, gut.Err(false, "student query execution failed", tx.Error)
	}

	// * execute check query
	var checkResult []map[string]any
	tx = gorm.Raw(checkQuery).Scan(&checkResult)
	if tx.Error != nil {
		return false, gut.Err(false, "check query execution failed", tx.Error)
	}

	// * compare results using reflect.DeepEqual
	queryPassed := reflect.DeepEqual(studentResult, checkResult)

	return queryPassed, nil
}
