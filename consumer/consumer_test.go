package consumer

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func amqpHost() string {
	return os.Getenv("RABBIT_URL")
}

func produce(topic, exchange, exchangeType, body string) {
	conn, _ := amqp.Dial(amqpHost())
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	_ = ch.ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil)
	p := amqp.Publishing{ContentType: "text/plain", Body: []byte(body)}
	_ = ch.Publish(exchange, topic, false, false, p)
}

func TestConsumer_Topic(t *testing.T) {
	topic := uuid.New().String()
	exchange := uuid.New().String()

	ch := make(chan Message)
	c := Consumer{
		Host:         amqpHost(),
		RoutingKey:   topic,
		Exchange:     exchange,
		ExchangeType: "topic",
		Ch:           ch,
	}

	go c.Consume()

	time.Sleep(time.Second)
	produce(topic, exchange, "topic", "Testing world")

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
	topic := uuid.New().String()
	exchange := uuid.New().String()
	ch := make(chan Message)
	c := Consumer{
		Host:         amqpHost(),
		RoutingKey:   topic,
		Exchange:     exchange,
		ExchangeType: "direct",
		Ch:           ch,
	}

	go c.Consume()

	time.Sleep(time.Second)
	produce(topic, exchange, "direct", "Testing direct")

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
	exchange := uuid.New().String()
	ch := make(chan Message)
	c := Consumer{
		Host:         amqpHost(),
		Exchange:     exchange,
		ExchangeType: "fanout",
		Ch:           ch,
	}

	go c.Consume()

	time.Sleep(time.Second)
	produce("any", exchange, "fanout", "Testing fanout")

	select {
	case msg, _ := <-ch:
		assert.Equal(t, "Testing fanout", string(msg.Body))
		assert.Equal(t, exchange, msg.Exchange)
		assert.Equal(t, "text/plain", string(msg.ContentType))
	case <-time.After(time.Second):
		t.Error("Expected to receive message, received nothing")
	}
}
