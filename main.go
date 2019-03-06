package main

import (
	"github.com/uesteibar/asyncapi-watcher/asyncapi/repos/messages_repo"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
	"github.com/uesteibar/asyncapi-watcher/watcher"
	"github.com/uesteibar/asyncapi-watcher/web/api"
)

func main() {
	repo := messages_repo.New(db.DB{})
	repo.Migrate()

	go watcher.Watch()

	api.Start()
}
