package consumer

import (
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

const AMQPHost = "amqp://guest:guest@localhost"

func produce(topic, body string) {
	conn, _ := amqp.Dial(AMQPHost)
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	_ = ch.ExchangeDeclare("/", "topic", true, false, false, false, nil)
	p := amqp.Publishing{ContentType: "text/plain", Body: []byte(body)}
	_ = ch.Publish("/", topic, false, false, p)
}

func TestConsumer(t *testing.T) {
	ch := make(chan Message)
	c := Consumer{
		Host:       "amqp://guest:guest@localhost",
		RoutingKey: "test.key",
		Ch:         ch,
	}

	produce("test.key", "Testing world")

	go c.Consume()

	select {
	case msg, _ := <-ch:
		assert.Equal(t, "Testing world", string(msg.Body))
		assert.Equal(t, "test.key", string(msg.RoutingKey))
		assert.Equal(t, "text/plain", string(msg.ContentType))
	case <-time.After(time.Second):
		t.Error("Expected to receive message, received nothing")
	}
}
