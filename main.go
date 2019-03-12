package main

import (
	"github.com/uesteibar/asyncapi-watcher/asyncapi/repos/messages_repo"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
	"github.com/uesteibar/asyncapi-watcher/watcher"
	"github.com/uesteibar/asyncapi-watcher/web/api"
)

const amqpHost = "amqp://guest:guest@localhost"
const exchange = "/"

func main() {
	repo := messages_repo.New(db.DB{})
	repo.Migrate()

	w := watcher.New(watcher.Config{
		Host:       amqpHost,
		RoutingKey: "#",
		Exchange:   exchange,
	})

	go w.Watch()

	api.Start()
}
