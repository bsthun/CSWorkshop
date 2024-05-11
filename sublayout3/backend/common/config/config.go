package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"backend/common"
	"backend/type/common"
	"backend/util/text"
)

func Init() {
	// * Declare struct
	config := new(common.Config)

	// * Load configurations to struct
	yml, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Unable to read configuration file", err)
	}
	if err := yaml.Unmarshal(yml, config); err != nil {
		log.Fatal("Unable to parse configuration file", err)
	}

	// * Validate configurations
	if err := text.Validator.Struct(config); err != nil {
		log.Fatal("Invalid configuration file", err)
	}

	// * Set global config
	cc.Config = config
}
