package watcher

import (
	"log"

	"github.com/uesteibar/asyncapi-watcher/asyncapi/analyzer"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/persister"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/consumer"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
)

type Config struct {
	Host       string
	RoutingKey string
	Exchange   string
}

type Watcher struct {
	Config Config
}

func New(c Config) Watcher {
	return Watcher{Config: c}
}

// Watch the amqp server for incoming messages and store the spec
func (w Watcher) Watch() {
	chConsumed := make(chan consumer.Message)
	c := consumer.Consumer{
		Host:       w.Config.Host,
		RoutingKey: w.Config.RoutingKey,
		Exchange:   w.Config.Exchange,
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
