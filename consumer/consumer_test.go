package consumer

import (
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

const AMQPHost = "amqp://guest:guest@localhost"

func produce(topic, exchange, exchangeType, body string) {
	conn, _ := amqp.Dial(AMQPHost)
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	_ = ch.ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil)
	p := amqp.Publishing{ContentType: "text/plain", Body: []byte(body)}
	_ = ch.Publish(exchange, topic, false, false, p)
}

func TestConsumer_Topic(t *testing.T) {
	topic := "test.key"
	exchange := "/topic-exchange"
	ch := make(chan Message)
	c := Consumer{
		Host:         AMQPHost,
		RoutingKey:   topic,
		Exchange:     exchange,
		ExchangeType: "topic",
		Ch:           ch,
	}

	produce(topic, exchange, "topic", "Testing world")

	go c.Consume()

	select {
	case msg, _ := <-ch:
		assert.Equal(t, "Testing world", string(msg.Body))
		assert.Equal(t, topic, string(msg.RoutingKey))
		assert.Equal(t, exchange, msg.Exchange)
		assert.Equal(t, "text/plain", string(msg.ContentType))
	case <-time.After(time.Second):
		t.Error("Expected to receive message, received nothing")
	}
}

func TestConsumer_Direct(t *testing.T) {
	topic := "test.key"
	exchange := "/direct-exchange"
	ch := make(chan Message)
	c := Consumer{
		Host:         AMQPHost,
		RoutingKey:   topic,
		Exchange:     exchange,
		ExchangeType: "direct",
		Ch:           ch,
	}

	produce(topic, exchange, "direct", "Testing direct")

	go c.Consume()

	select {
	case msg, _ := <-ch:
		assert.Equal(t, "Testing direct", string(msg.Body))
		assert.Equal(t, topic, string(msg.RoutingKey))
		assert.Equal(t, exchange, msg.Exchange)
		assert.Equal(t, "text/plain", string(msg.ContentType))
	case <-time.After(time.Second):
		t.Error("Expected to receive message, received nothing")
	}
}

func TestConsumer_Fanout(t *testing.T) {
	exchange := "/fanout-exchange"
	ch := make(chan Message)
	c := Consumer{
		Host:         AMQPHost,
		Exchange:     exchange,
		ExchangeType: "fanout",
		Ch:           ch,
	}

	produce("any", exchange, "fanout", "Testing fanout")

	go c.Consume()

	select {
	case msg, _ := <-ch:
		assert.Equal(t, "Testing fanout", string(msg.Body))
		assert.Equal(t, exchange, msg.Exchange)
		assert.Equal(t, "text/plain", string(msg.ContentType))
	case <-time.After(time.Second):
		t.Error("Expected to receive message, received nothing")
	}
}
