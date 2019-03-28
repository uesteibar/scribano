package main

import (
	"flag"

	"github.com/uesteibar/asyncapi-watcher/asyncapi/repos/messages_repo"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
	"github.com/uesteibar/asyncapi-watcher/watcher"
	"github.com/uesteibar/asyncapi-watcher/watcher/config"
	yamlconfig "github.com/uesteibar/asyncapi-watcher/watcher/config/parsers/yaml_config"
	"github.com/uesteibar/asyncapi-watcher/web/api"
)

func getLoader() config.Loader {
	var url string
	flag.StringVar(&url, "u", "", "url to the config file.")
	var path string
	flag.StringVar(&path, "f", "", "path to the config file.")
	flag.Parse()

	if url != "" {
		return &config.URLLoader{Source: url}
	} else if path != "" {
		return &config.PathLoader{Source: path}
	} else {
		panic("Missing configuration")
	}
}

func main() {
	repo := messages_repo.New(db.DB{})
	repo.Migrate()

	loader := getLoader()
	configParser := yamlconfig.New(loader)

	configs, err := configParser.Parse()

	if err != nil {
		panic(err)
	}

	for _, c := range configs {
		w := watcher.New(c)
		go w.Watch()

	}

	api.Start()
}
