package main

import (
	"flag"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/repos/messages_repo"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
	"github.com/uesteibar/asyncapi-watcher/watcher"
	yamlconfig "github.com/uesteibar/asyncapi-watcher/watcher/config/parsers/yaml_config"
	"github.com/uesteibar/asyncapi-watcher/web/api"
)

const configFile = "./fixtures/test/yaml_config.yml"

func configFilePath() string {
	var path string
	flag.StringVar(&path, "f", "", "path to the config file.")
	flag.Parse()

	return path
}

func main() {
	repo := messages_repo.New(db.DB{})
	repo.Migrate()

	configLoader := yamlconfig.New(configFilePath())

	if config, err := configLoader.Parse(); err == nil {
		w := watcher.New(config)

		go w.Watch()

		api.Start()
	} else {
		panic(err)
	}
}
