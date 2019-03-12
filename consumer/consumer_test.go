package consumer

import (
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

const AMQPHost = "amqp://guest:guest@localhost"
const Exchange = "/"

func produce(topic, body string) {
	conn, _ := amqp.Dial(AMQPHost)
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	_ = ch.ExchangeDeclare(Exchange, "topic", true, false, false, false, nil)
	p := amqp.Publishing{ContentType: "text/plain", Body: []byte(body)}
	_ = ch.Publish(Exchange, topic, false, false, p)
}

func TestConsumer(t *testing.T) {
	topic := "test.key"
	ch := make(chan Message)
	c := Consumer{
		Host:       AMQPHost,
		RoutingKey: topic,
		Exchange:   Exchange,
		Ch:         ch,
	}

	produce(topic, "Testing world")

	go c.Consume()

	select {
	case msg, _ := <-ch:
		assert.Equal(t, "Testing world", string(msg.Body))
		assert.Equal(t, topic, string(msg.RoutingKey))
		assert.Equal(t, "text/plain", string(msg.ContentType))
	case <-time.After(time.Second):
		t.Error("Expected to receive message, received nothing")
	}
}
