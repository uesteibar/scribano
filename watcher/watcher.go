package watcher

import (
	"github.com/uesteibar/asyncapi-watcher/asyncapi/analyzer"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/persister"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/consumer"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
	"log"
)

func Watch() {
	chConsumed := make(chan consumer.Message)
	c := consumer.Consumer{
		Host:       "amqp://guest:guest@localhost",
		RoutingKey: "key.test",
		Ch:         chConsumed,
	}

	go c.Consume()

	chAnalyzed := make(chan spec.MessageSpec)

	a := analyzer.Analyzer{ChIn: chConsumed, ChOut: chAnalyzed}

	go a.Watch()

	chPersisted := make(chan spec.MessageSpec)
	p := persister.New(chAnalyzed, chPersisted, db.DB{})
	p.Watch()

	for msg := range chPersisted {
		log.Printf("INFO Persisted: %+v", msg)
	}
}
