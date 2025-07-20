package database

import (
	"backend/generate/psql"
	"backend/type/common"
	"context"
	"database/sql"
	"github.com/bsthun/gut"
)

type Database struct {
	PQuerier *psql.Queries
	PConn    *sql.DB
	CConn    *sql.DB
}

func (r *Database) P() psql.PQuerier {
	return r.PQuerier
}

func (r *Database) Ptx(context context.Context, opts *sql.TxOptions) (common.DatabaseTx, psql.PQuerier) {
	tx, err := r.PConn.BeginTx(context, opts)
	querier := r.PQuerier.WithTx(tx)
	if err != nil {
		gut.Fatal("failed to begin transaction", err)
	}
	return tx, querier
}
