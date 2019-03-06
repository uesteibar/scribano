package messages_repo

import (
	"encoding/json"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
)

const gormNotFound = "record not found"

type MessageSpec struct {
	Topic   string `gorm:"primary_key"`
	Payload []byte
}

type MessagesRepo struct {
	db db.Database
}

type ErrNotFound struct {
	message string
}

func NewErrNotFound() *ErrNotFound {
	return &ErrNotFound{}
}
func (e *ErrNotFound) Error() string {
	return "NOT_FOUND"
}

func New(db db.Database) *MessagesRepo {
	return &MessagesRepo{db: db}
}

func (r *MessagesRepo) Migrate() {
	conn := r.db.Open()
	conn.AutoMigrate(&MessageSpec{})
}

func transformMsg(msg spec.MessageSpec) (MessageSpec, error) {
	payload, err := json.Marshal(msg.Payload)

	if err != nil {
		return MessageSpec{}, err
	}

	return MessageSpec{Topic: msg.Topic, Payload: payload}, nil
}

func transformToMsg(msg MessageSpec) spec.MessageSpec {
	var p spec.PayloadSpec
	json.Unmarshal(msg.Payload, &p)

	return spec.MessageSpec{Topic: msg.Topic, Payload: p}
}

func (r *MessagesRepo) Create(msg spec.MessageSpec) error {
	conn := r.db.Open()
	defer conn.Close()

	if m, err := transformMsg(msg); err != nil {
		return err
	} else {
		return conn.Create(&m).Error
	}
}

func (r *MessagesRepo) Find(topic string) (spec.MessageSpec, error) {
	conn := r.db.Open()
	defer conn.Close()

	var m MessageSpec
	if err := conn.First(&m, "topic = ?", topic).Error; err == nil {
		var p spec.PayloadSpec
		json.Unmarshal(m.Payload, &p)

		messageSpec := spec.MessageSpec{Topic: m.Topic, Payload: p}

		return messageSpec, nil
	} else if err.Error() == gormNotFound {
		return spec.MessageSpec{}, NewErrNotFound()
	} else {
		return spec.MessageSpec{}, err
	}
}

func (r *MessagesRepo) FindAll() ([]spec.MessageSpec, error) {
	conn := r.db.Open()
	defer conn.Close()
	var msgs []MessageSpec

	if err := conn.Find(&msgs).Error; err == nil {
		var messageSpecs []spec.MessageSpec

		for _, m := range msgs {
			msg := transformToMsg(m)
			messageSpecs = append(messageSpecs, msg)
		}

		return messageSpecs, nil
	} else if err.Error() == gormNotFound {
		return []spec.MessageSpec{}, NewErrNotFound()
	} else {
		return []spec.MessageSpec{}, err
	}
}

func (r *MessagesRepo) Update(msg spec.MessageSpec) error {
	conn := r.db.Open()
	defer conn.Close()

	if m, err := transformMsg(msg); err != nil {
		return err
	} else {
		return conn.Save(&m).Error
	}
}
