package messages_repo

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
	"testing"
)

func TestRepo(t *testing.T) {
	repo := New(db.TestDB{})
	repo.Migrate()

	topic := uuid.New().String()
	msg := spec.MessageSpec{
		Topic: topic,
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

	m, err := repo.Find(topic)

	assert.Nil(t, err)
	assert.Equal(t, topic, m.Topic)
	assert.Equal(t, msg.Payload, m.Payload)

	notFoundMsg, err := repo.Find("wrong.topic")
	assert.Equal(t, "", notFoundMsg.Topic)
	switch err.(type) {
	case *ErrNotFound:
		assert.Equal(t, "NOT_FOUND", err.Error())
	default:
		t.Error("Expected error, got nothing")
	}

	newPayload := spec.PayloadSpec{
		Type: "object",
		Fields: []spec.FieldSpec{
			spec.FieldSpec{Name: "name", Type: "string"},
			spec.FieldSpec{Name: "age", Type: "number"},
		},
	}
	m.Payload = newPayload

	err = repo.Update(m)
	assert.Nil(t, err)
	um, _ := repo.Find(topic)
	assert.Equal(t, m, um)
}
