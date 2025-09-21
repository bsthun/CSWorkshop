package submissionProcedure

import (
	"backend/common/config"
	"backend/type/common"
	"context"

	"github.com/bsthun/gut"
	"github.com/openai/openai-go"
	"gorm.io/gorm"
)

type Server interface {
	ServeCheckQuery(ctx context.Context, studentAnswer string, checkQuery string, databaseName string) (bool, *gut.ErrorInstance)
	ServeCheckPrompt(ctx context.Context, studentAnswer string, checkPrompt string) (bool, string, *gut.ErrorInstance)
}

type Service struct {
	database     common.Database
	gorm         *gorm.DB
	config       *config.Config
	openaiClient *openai.Client
}

func Serve(
	database common.Database,
	gorm *gorm.DB,
	config *config.Config,
	openai *openai.Client,
) Server {
	return &Service{
		database:     database,
		gorm:         gorm,
		config:       config,
		openaiClient: openai,
	}
}
