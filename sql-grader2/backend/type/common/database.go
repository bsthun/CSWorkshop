package common

import (
	"backend/generate/psql"
	"context"
	"database/sql"
)

type Database interface {
	P() psql.PQuerier
	Ptx(context context.Context, opts *sql.TxOptions) (DatabaseTx, psql.PQuerier)
}

type DatabaseTx interface {
	Commit() error
	Rollback() error
}
