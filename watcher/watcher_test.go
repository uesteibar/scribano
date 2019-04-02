package watcher

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/scribano/asyncapi/repos/messages_repo"
	"github.com/uesteibar/scribano/asyncapi/spec"
	"github.com/uesteibar/scribano/storage/db"
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
			"canDrive": false,
			"friends": [
				{ "name": "jose" },
				{ "name": "maria" }
			],
			"car": {
				"brand": "ford"
			}
		}
	`
	topic := uuid.New().String()
	topic = "key.test"
	produce(topic, body)

	watcher := New(Config{Host: AMQPHost, RoutingKey: "#", Exchange: Exchange})
	go watcher.Watch()

	time.Sleep(time.Millisecond * 1000)

	expectedFields := []spec.FieldSpec{
		spec.FieldSpec{Name: "name", Type: "string"},
		spec.FieldSpec{Name: "age", Type: "integer"},
		spec.FieldSpec{Name: "canDrive", Type: "boolean"},
		spec.FieldSpec{Name: "friends", Type: "array",
			Item: &spec.FieldSpec{
				Type: "object",
				Fields: []spec.FieldSpec{
					spec.FieldSpec{Name: "name", Type: "string"},
				},
			},
		},
		spec.FieldSpec{
			Name: "car",
			Type: "object",
			Fields: []spec.FieldSpec{
				spec.FieldSpec{Name: "brand", Type: "string"},
			},
		},
	}

	m, err := repo.Find(topic)
	assert.Nil(t, err)
	assert.Equal(t, topic, m.Topic)
	assert.Equal(t, Exchange, m.Exchange)
	assert.Equal(t, "object", m.Payload.Type)
	assert.ElementsMatch(t, expectedFields, m.Payload.Fields)
}
