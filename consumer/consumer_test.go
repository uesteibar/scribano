package consumer

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

const AMQPHost = "amqp://guest:guest@localhost"

func produce(t *testing.T) {
	conn, _ := amqp.Dial(AMQPHost)
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	q, _ := ch.QueueDeclare("test.key", false, false, false, false, nil)
	p := amqp.Publishing{ContentType: "text/plain", Body: []byte("Testing world")}
	_ = ch.Publish("", q.Name, false, false, p)
}

func TestConsumer(t *testing.T) {
	ch := make(chan Message)
	c := Consumer{
		Host:       "amqp://guest:guest@localhost",
		RoutingKey: "test.key",
		Ch:         ch,
	}

	produce(t)

	go c.Consume()

	msg, _ := <-ch

	assert.Equal(t, "Testing world", string(msg.Body))
	assert.Equal(t, "test.key", string(msg.RoutingKey))
	assert.Equal(t, "text/plain", string(msg.ContentType))
}
