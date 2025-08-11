package config

import (
	"backend/type/enum"
	"os"

	"github.com/bsthun/gut"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment       *enum.Environment `yaml:"environment" validate:"required"`
	WebRoot           *string           `yaml:"webRoot" validate:"omitempty"`
	WebListen         [2]*string        `yaml:"webListen" validate:"required"`
	FrontendUrl       *string           `yaml:"frontendUrl" validate:"required"`
	Secret            *string           `yaml:"secret" validate:"required"`
	PostgresDsn       *string           `yaml:"postgresDsn" validate:"required"`
	MysqlDsn          *string           `yaml:"mysqlDsn" validate:"required"`
	OauthClientId     *string           `yaml:"oauthClientId" validate:"required"`
	OauthClientSecret *string           `yaml:"oauthClientSecret" validate:"required"`
	OauthEndpoint     *string           `yaml:"oauthEndpoint" validate:"required"`
	OpenaiBaseUrl     *string           `yaml:"openaiBaseUrl" validate:"required"`
	OpenaiModel       *string           `yaml:"openaiModel" validate:"required"`
	OpenaiApiKey      *string           `yaml:"openaiApiKey" validate:"required"`
}

func Init() *Config {
	// * Parse arguments
	path := os.Getenv("BACKEND_CONFIG_PATH")
	if path == "" {
		path = "config.yml"
	}

	// * Declare struct
	config := new(Config)

	// * Read config
	yml, err := os.ReadFile(path)
	if err != nil {
		gut.Fatal("Unable to read configuration file", err)
	}

	// * Parse config
	if err := yaml.Unmarshal(yml, config); err != nil {
		gut.Fatal("Unable to parse configuration file", err)
	}

	// * Validate config
	if err := gut.Validate(config); err != nil {
		gut.Fatal("Invalid configuration", err)
	}

	// * apply secret key
	var bytes = []byte(*config.Secret)
	if len(bytes) < 16 {
		for i := len(bytes); i < 16; i++ {
			bytes = append(bytes, 0)
		}
	}
	if err := gut.SetIdEncoderKey(bytes[:16]); err != nil {
		gut.Fatal("unable to set secret key", err)
	}

	return config
}
