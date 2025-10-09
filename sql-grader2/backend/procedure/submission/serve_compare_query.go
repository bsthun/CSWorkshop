package submissionProcedure

import (
	"backend/type/payload"
	"backend/util/orm"
	"context"
	"sort"

	"github.com/bsthun/gut"
)

func (r *Service) ServeCompareQuery(ctx context.Context, studentAnswer string, checkQuery string, databaseName string) (*payload.ExamQuestionSubmissionCompare, *gut.ErrorInstance) {
	// * connect to exam database
	gorm, er := orm.Instance(*r.config.MysqlDsn, databaseName)
	if er != nil {
		return nil, gut.Err(false, "failed to connect to exam database instance", er)
	}

	// * execute student answer query
	var studentResult []map[string]any
	tx := gorm.Raw(studentAnswer).Scan(&studentResult)
	if tx.Error != nil {
		return nil, gut.Err(false, "student query execution failed", tx.Error)
	}

	// * execute check query
	var checkResult []map[string]any
	tx = gorm.Raw(checkQuery).Scan(&checkResult)
	if tx.Error != nil {
		return nil, gut.Err(false, "check query execution failed", tx.Error)
	}

	// * convert student result to table format
	studentTable := r.convertToTable(studentResult)

	// * convert check result to table format
	checkTable := r.convertToTable(checkResult)

	// * return comparison result
	return &payload.ExamQuestionSubmissionCompare{
		Submitted: studentTable,
		Expected:  checkTable,
	}, nil
}

func (r *Service) convertToTable(result []map[string]any) *payload.ExamQuestionSubmissionCompareTable {
	// * handle empty result
	if len(result) == 0 {
		return &payload.ExamQuestionSubmissionCompareTable{
			Header: []*string{},
			Rows:   [][]any{},
		}
	}

	// * extract headers from first row and sort alphabetically
	var headers []*string
	for key := range result[0] {
		headers = append(headers, gut.Ptr(key))
	}

	// * sort headers alphabetically
	sort.Slice(headers, func(i, j int) bool {
		return *headers[i] < *headers[j]
	})

	// * extract rows
	var rows [][]any
	for _, row := range result {
		var rowData []any
		for _, header := range headers {
			rowData = append(rowData, row[*header])
		}
		rows = append(rows, rowData)
	}

	return &payload.ExamQuestionSubmissionCompareTable{
		Header: headers,
		Rows:   rows,
	}
}
