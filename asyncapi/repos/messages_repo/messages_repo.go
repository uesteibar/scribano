package messages_repo

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
)

type MessageSpec struct {
	gorm.Model

	Topic   string
	Payload []byte
}

type MessagesRepo struct {
	db db.Database
}

func New(db db.Database) *MessagesRepo {
	return &MessagesRepo{db: db}
}

func (r *MessagesRepo) Migrate() {
	conn := r.db.Open()
	conn.AutoMigrate(&MessageSpec{})
}

func (r *MessagesRepo) Create(msg spec.MessageSpec) error {
	conn := r.db.Open()
	defer conn.Close()

	payload, err := json.Marshal(msg.Payload)

	if err != nil {
		return err
	}

	conn.Create(&MessageSpec{
		Topic:   msg.Topic,
		Payload: payload,
	})
	return nil
}

func (r *MessagesRepo) Find(topic string) spec.MessageSpec {
	conn := r.db.Open()
	defer conn.Close()

	var m MessageSpec
	conn.First(&m, "topic = ?", topic)

	var p spec.PayloadSpec
	json.Unmarshal(m.Payload, &p)

	messageSpec := spec.MessageSpec{
		Topic:   m.Topic,
		Payload: p,
	}

	return messageSpec
}
