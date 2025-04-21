package configs

import (
	"log"
	"os"
	"runtime-manager/internals/models"
	"runtime-manager/internals/pkg"

	"gopkg.in/yaml.v3"
)

func Parser(config_file string) *models.Config {
	yamlFile, err := os.ReadFile(pkg.CONFIG_FILE_PATH)
	if err != nil {
		log.Printf("Error while reading yaml file: %v", err)
	}
	var config models.Config

	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		log.Printf("Error unmarshalling YAML: %v", err)
	}
	return &config
}
