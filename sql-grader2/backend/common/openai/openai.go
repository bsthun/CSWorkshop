package openai

import (
	"backend/common/config"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func Init(config *config.Config) *openai.Client {
	client := openai.NewClient(
		option.WithBaseURL(*config.OpenaiBaseUrl),
		option.WithAPIKey(*config.OpenaiApiKey),
	)

	return &client
}
