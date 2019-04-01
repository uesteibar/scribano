package watcher

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/repos/messages_repo"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
)

const AMQPHost = "amqp://guest:guest@localhost"
const Exchange = "/other-exchange"

func produce(topic, body string) {
	conn, _ := amqp.Dial(AMQPHost)
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	_ = ch.ExchangeDeclare(Exchange, "topic", true, false, false, false, nil)
	p := amqp.Publishing{ContentType: "application/json", Body: []byte(body)}
	_ = ch.Publish(Exchange, topic, false, false, p)
}

func TestEndToEnd(t *testing.T) {
	repo := messages_repo.New(db.DB{})
	repo.Migrate()
	body := `
		{
			"name": "infer type",
			"age": 27,
			"canDrive": false
		}
	`
	topic := uuid.New().String()
	topic = "key.test"
	produce(topic, body)

	watcher := New(Config{Host: AMQPHost, RoutingKey: "#", Exchange: Exchange})
	go watcher.Watch()

	time.Sleep(time.Millisecond * 1000)

	expected := spec.MessageSpec{
		Topic:    topic,
		Exchange: Exchange,
		Payload: spec.PayloadSpec{
			Type: "object",
			Fields: []spec.FieldSpec{
				spec.FieldSpec{Name: "name", Type: "string"},
				spec.FieldSpec{Name: "age", Type: "integer"},
				spec.FieldSpec{Name: "canDrive", Type: "boolean"},
			},
		},
	}

	m, err := repo.Find(topic)
	assert.Nil(t, err)
	assert.Equal(t, expected, m)
}
