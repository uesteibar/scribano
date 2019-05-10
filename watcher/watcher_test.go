package watcher

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/scribano/asyncapi/repos/messages_repo"
	"github.com/uesteibar/scribano/asyncapi/spec"
	"github.com/uesteibar/scribano/storage/db"
)

func amqpHost() string {
	return os.Getenv("RABBIT_URL")
}

func produce(exchange, topic, body string) {
	conn, _ := amqp.Dial(amqpHost())
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	_ = ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	p := amqp.Publishing{ContentType: "application/json", Body: []byte(body)}
	_ = ch.Publish(exchange, topic, false, false, p)
}

func TestEndToEnd(t *testing.T) {
	repo := messages_repo.New(db.DB{})
	repo.Migrate()
	body := `
		{
			"name": "infer type",
			"age": 27,
			"canDrive": false,
			"birthDate": "1991-08-29",
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
	exchange := uuid.New().String()

	watcher := New(Config{Host: amqpHost(), RoutingKey: "#", Exchange: exchange, ExchangeType: "topic"})
	go watcher.Watch()

	time.Sleep(time.Second)
	produce(exchange, topic, body)
	time.Sleep(time.Second)

	expectedFields := []*spec.FieldSpec{
		&spec.FieldSpec{Name: "name", Type: "string"},
		&spec.FieldSpec{Name: "age", Type: "integer"},
		&spec.FieldSpec{Name: "canDrive", Type: "boolean"},
		&spec.FieldSpec{Name: "birthDate", Type: "string", Format: "date"},
		&spec.FieldSpec{Name: "friends", Type: "array",
			Item: &spec.FieldSpec{
				Type: "object",
				Fields: []*spec.FieldSpec{
					&spec.FieldSpec{Name: "name", Type: "string"},
				},
			},
		},
		&spec.FieldSpec{
			Name: "car",
			Type: "object",
			Fields: []*spec.FieldSpec{
				&spec.FieldSpec{Name: "brand", Type: "string"},
			},
		},
	}

	m, err := repo.Find(topic)
	assert.Nil(t, err)
	assert.Equal(t, topic, m.Topic)
	assert.Equal(t, exchange, m.Exchange)
	assert.Equal(t, "object", m.Payload.Type)
	assert.ElementsMatch(t, expectedFields, m.Payload.Fields)
}
