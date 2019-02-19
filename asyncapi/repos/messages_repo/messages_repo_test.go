package messages_repo

import (
	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
	"testing"
)

func TestSave(t *testing.T) {
	repo := New(db.TestDB{})
	repo.Migrate()

	msg := spec.MessageSpec{
		Topic: "some.topic",
		Payload: spec.PayloadSpec{
			Type: "object",
			Fields: []spec.FieldSpec{
				spec.FieldSpec{Name: "name", Type: "string"},
				spec.FieldSpec{Name: "age", Type: "float"},
			},
		},
	}

	err := repo.Create(msg)
	assert.Nil(t, err)

	m := repo.Find("some.topic")

	assert.Equal(t, "some.topic", m.Topic)
	assert.Equal(t, msg.Payload, m.Payload)
}
