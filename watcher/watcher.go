package watcher

import (
	"log"

	"github.com/uesteibar/scribano/asyncapi/analyzer"
	"github.com/uesteibar/scribano/asyncapi/persister"
	"github.com/uesteibar/scribano/asyncapi/spec"
	"github.com/uesteibar/scribano/consumer"
	"github.com/uesteibar/scribano/storage/db"
)

// Config for the Watcher
type Config struct {
	Host         string
	RoutingKey   string
	Exchange     string
	ExchangeType string
}

// Watcher orchestrates consuming, analyzing and persisting amqp messages
type Watcher struct {
	Config Config
}

// New creates a new watcher for the given config
func New(c Config) Watcher {
	return Watcher{Config: c}
}

// Watch the amqp server for incoming messages and store the spec
func (w Watcher) Watch() {
	dbConn := db.DB{}

	chConsumed := make(chan consumer.Message)
	c := consumer.Consumer{
		Host:         w.Config.Host,
		RoutingKey:   w.Config.RoutingKey,
		Exchange:     w.Config.Exchange,
		ExchangeType: w.Config.ExchangeType,
		Ch:           chConsumed,
	}

	go c.Consume()

	chAnalyzed := make(chan spec.MessageSpec)

	a := analyzer.New(chConsumed, chAnalyzed, dbConn)

	go a.Watch()

	chPersisted := make(chan spec.MessageSpec)
	p := persister.New(chAnalyzed, chPersisted, dbConn)
	go p.Watch()

	for msg := range chPersisted {
		log.Printf("INFO Persisted: %+v", msg)
	}

	log.Printf("INFO finished running watcher")
}
