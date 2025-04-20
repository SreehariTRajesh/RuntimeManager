package utils

import (
	"log"
	"runtime-manager/configs"
	"runtime-manager/internals/pkg"
	"sync"
)

// following module initialiases the docker network
var (
	config_manager_instance *configs.Config
	once                    sync.Once
)

func GetConfig() *configs.Config {
	once.Do(func() {
		config_manager_instance = configs.Parser(pkg.CONFIG_FILE_PATH)
		if config_manager_instance == nil {
			log.Fatalf("failed to load configuration from %s", pkg.CONFIG_FILE_PATH)
		}
	})
	return config_manager_instance
}
