package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bsthun/gut"
	"reflect"
	"regexp"
	"strings"
)

type Wrapper struct {
	db *sql.DB
}

var queryNameRegex = regexp.MustCompile(`--\s*name:\s*(\w+)`)

func PrintQueryInfo(query string, args ...any) {
	// * extract query name using regex
	matches := queryNameRegex.FindStringSubmatch(query)
	queryName := "N/A"
	if len(matches) > 1 {
		queryName = matches[1]
	}

	// * resolve and format arguments
	resolvedArgs := make([]string, 0, len(args))
	for _, arg := range args {
		resolvedArg := []rune(fmt.Sprintf("%v", resolveArg(arg)))
		if len(resolvedArg) > 64 {
			resolvedArg = append(resolvedArg[:64], []rune("...")...)
		}
		resolvedArgs = append(resolvedArgs, string(resolvedArg))
	}

	// * format the output string
	argsStr := strings.Join(resolvedArgs, ", ")
	if len(argsStr) > 0 {
		argsStr = "(" + argsStr + ")"
	}

	// * print the formatted output
	fmt.Printf("%s%sDatabase%s%s%s\n", gut.PrintTime(), gut.PrintSeparator(), gut.PrintSeparator(), queryName, argsStr)
}

func resolveArg(arg any) any {
	if arg == nil {
		return nil
	}

	val := reflect.ValueOf(arg)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		derefVal := val.Elem().Interface()
		return resolveArg(derefVal)
	}

	// * check case of json
	if rawMsg, ok := arg.(json.RawMessage); ok {
		return string(rawMsg)
	}

	return arg
}

func (r *Wrapper) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	PrintQueryInfo(query, args...)
	return r.db.ExecContext(ctx, query, args...)
}

func (r *Wrapper) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	PrintQueryInfo(query)
	return r.db.PrepareContext(ctx, query)
}

func (r *Wrapper) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	PrintQueryInfo(query, args...)
	return r.db.QueryContext(ctx, query, args...)
}

func (r *Wrapper) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	PrintQueryInfo(query, args...)
	return r.db.QueryRowContext(ctx, query, args...)
}
