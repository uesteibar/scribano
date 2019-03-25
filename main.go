package main

import (
	"flag"

	"github.com/uesteibar/asyncapi-watcher/asyncapi/repos/messages_repo"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
	"github.com/uesteibar/asyncapi-watcher/watcher"
	yamlconfig "github.com/uesteibar/asyncapi-watcher/watcher/config/parsers/yaml_config"
	"github.com/uesteibar/asyncapi-watcher/web/api"
)

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

	configs, err := configLoader.Parse()

	if err != nil {
		panic(err)
	}

	for _, c := range configs {
		w := watcher.New(c)
		go w.Watch()

	}
	api.Start()
}
