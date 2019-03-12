package consumer

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Consumer struct {
	Host       string
	RoutingKey string
	Exchange   string
	Ch         chan Message
}

type Message struct {
	Body        []byte
	ContentType string
	RoutingKey  string
}

func transformMessage(msg amqp.Delivery) Message {
	return Message{
		Body:        msg.Body,
		ContentType: msg.ContentType,
		RoutingKey:  msg.RoutingKey,
	}
}

func (c *Consumer) Consume() {
	conn, err := amqp.Dial(c.Host)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		c.Exchange, // name
		"topic",    // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		c.RoutingKey, // name
		false,        // durable
		false,        // delete when usused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
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
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received: %+v", d)
			c.Ch <- transformMessage(d)
		}
	}()

	log.Printf(" [*] Waiting for messages.")

	<-forever
}
