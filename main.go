package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/analyzer"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/builder"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/consumer"
)

const AMQPHost = "amqp://guest:guest@localhost"

func produce() {
	conn, _ := amqp.Dial(AMQPHost)
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	q, _ := ch.QueueDeclare("super.key", false, false, false, false, nil)
	sampleBody := `
		{
			"name": "infer type",
			"age": 27,
			"canDrive": false
		}
	`
	p := amqp.Publishing{ContentType: "application/json", Body: []byte(sampleBody)}
	_ = ch.Publish("", q.Name, false, false, p)
}

func main() {
	ch := make(chan consumer.Message)
	c := consumer.Consumer{
		Host:       AMQPHost,
		RoutingKey: "super.key",
		Ch:         ch,
	}
	go c.Consume()

	chOut := make(chan spec.MessageSpec)
	a := analyzer.Analyzer{ChIn: ch, ChOut: chOut}
	go a.Watch()

	produce()

	msg := <-chOut
	fmt.Printf("Received: %+v\n", msg)

	b := builder.SpecBuilder{}
	b.AddMessage(msg)

	res := b.Build()

	fmt.Printf("Build: %+v\n", res)
	encoded, _ := json.Marshal(res)

	fmt.Printf("Encoded: %+s\n", encoded)
}
