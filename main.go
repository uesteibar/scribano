package main

import (
	"github.com/uesteibar/asyncapi-watcher/asyncapi/repos/messages_repo"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
	"github.com/uesteibar/asyncapi-watcher/watcher"
)

func main() {
	repo := messages_repo.New(db.TestDB{})
	repo.Migrate()

	go watcher.Watch()
}
