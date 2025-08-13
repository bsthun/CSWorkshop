package entityProcedure

import (
	"backend/type/common"
	"backend/type/payload"
	"context"

	"github.com/bsthun/gut"
)

type Server interface {
	ServeSemesterList(ctx context.Context) (*payload.SemesterListResponse, *gut.ErrorInstance)
}

type Service struct {
	database common.Database
}

func Serve(database common.Database) Server {
	return &Service{
		database: database,
	}
}
