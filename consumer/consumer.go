package consumer

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Consumer consumes messages matching a given rabbitmq routing key
type Consumer struct {
	Host         string
	RoutingKey   string
	Exchange     string
	ExchangeType string
	Ch           chan Message
}

// Message is a message consumed from rabbitmq
type Message struct {
	Body        []byte
	ContentType string
	RoutingKey  string
	Exchange    string
	Delivery    amqp.Delivery
}

func (c *Consumer) transformMessage(msg amqp.Delivery) Message {
	return Message{
		Body:        msg.Body,
		ContentType: msg.ContentType,
		RoutingKey:  msg.RoutingKey,
		Exchange:    c.Exchange,
		Delivery:    msg,
	}
}

// Consume messages from rabbitmq
func (c *Consumer) Consume() {
	conn, err := amqp.Dial(c.Host)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		c.Exchange,     // name
		c.ExchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		fmt.Sprintf("scribano_watcher_%s", c.Exchange),
		false, // durable
		true,  // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,       // queue name
		c.RoutingKey, // routing key
		c.Exchange,   // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to consume from queue")

	log.Printf(" [*] Waiting for messages - queue: %s, broker: %s, exchange: %s, matcher: %s.", q.Name, c.Host, c.Exchange, c.RoutingKey)

	for d := range msgs {
		log.Printf("Received: %+v", d)
		c.Ch <- c.transformMessage(d)
	}
}
