package submissionProcedure

import (
	"backend/util/orm"
	"context"
	"reflect"
	"sort"

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

	// * sort both results by headers before comparison
	studentSorted := r.sortResultByHeaders(studentResult)
	checkSorted := r.sortResultByHeaders(checkResult)

	// * compare results using reflect.DeepEqual
	queryPassed := reflect.DeepEqual(studentSorted, checkSorted)

	return queryPassed, nil
}

func (r *Service) sortResultByHeaders(result []map[string]any) []map[string]any {
	if len(result) == 0 {
		return result
	}

	// * get sorted headers from first row
	var headers []string
	for key := range result[0] {
		headers = append(headers, key)
	}
	sort.Strings(headers)

	// * create new result with sorted column order
	var sortedResult []map[string]any
	for _, row := range result {
		sortedRow := make(map[string]any)
		for _, header := range headers {
			sortedRow[header] = row[header]
		}
		sortedResult = append(sortedResult, sortedRow)
	}

	return sortedResult
}
