package messagesrepo

import (
	"encoding/json"

	_ "github.com/jinzhu/gorm/dialects/sqlite" // required by the library
	"github.com/uesteibar/scribano/asyncapi/spec"
	"github.com/uesteibar/scribano/storage/db"
)

const gormNotFound = "record not found"

type messageSpec struct {
	Topic    string `gorm:"primary_key"`
	Exchange string
	Payload  []byte
}

// MessagesRepo is the repo for storing messages
type MessagesRepo struct {
	db db.Database
}

// ErrNotFound is returned when a message is not found
type ErrNotFound struct {
	message string
}

// NewErrNotFound returns a not found error
func NewErrNotFound() *ErrNotFound {
	return &ErrNotFound{}
}

func (e ErrNotFound) Error() string {
	return "NOT_FOUND"
}

// New creates a new messages repo
func New(db db.Database) *MessagesRepo {
	return &MessagesRepo{db: db}
}

// Migrate runs migrations for the messages repo
func (r *MessagesRepo) Migrate() {
	conn := r.db.Open()
	conn.AutoMigrate(&messageSpec{})
}

func transformMsg(msg spec.MessageSpec) (messageSpec, error) {
	payload, err := json.Marshal(msg.Payload)

	if err != nil {
		return messageSpec{}, err
	}

	return messageSpec{Topic: msg.Topic, Exchange: msg.Exchange, Payload: payload}, nil
}

func transformToMsg(msg messageSpec) spec.MessageSpec {
	var p spec.PayloadSpec
	json.Unmarshal(msg.Payload, &p)

	return spec.MessageSpec{Topic: msg.Topic, Exchange: msg.Exchange, Payload: p}
}

// Create a message spec
func (r *MessagesRepo) Create(msg spec.MessageSpec) error {
	conn := r.db.Open()
	defer conn.Close()

	m, err := transformMsg(msg)
	if err != nil {
		return err
	}

	return conn.Create(&m).Error
}

// Find a message spec by topic
func (r *MessagesRepo) Find(topic string) (spec.MessageSpec, error) {
	conn := r.db.Open()
	defer conn.Close()

	var m messageSpec
	if err := conn.First(&m, "topic = ?", topic).Error; err == nil {
		var p spec.PayloadSpec
		json.Unmarshal(m.Payload, &p)

		messageSpec := spec.MessageSpec{Topic: m.Topic, Exchange: m.Exchange, Payload: p}

		return messageSpec, nil
	} else if err.Error() == gormNotFound {
		return spec.MessageSpec{}, NewErrNotFound()
	} else {
		return spec.MessageSpec{}, err
	}
}

// FindAll returns all messages
func (r *MessagesRepo) FindAll() ([]spec.MessageSpec, error) {
	conn := r.db.Open()
	defer conn.Close()
	var msgs []messageSpec

	messageSpecs := []spec.MessageSpec{}
	err := conn.Find(&msgs).Error

	if err != nil {
		return messageSpecs, err
	}

	for _, m := range msgs {
		msg := transformToMsg(m)
		messageSpecs = append(messageSpecs, msg)
	}

	return messageSpecs, nil
}

//FindByExchange finds all messages for a exchange
func (r *MessagesRepo) FindByExchange(exchange string) ([]spec.MessageSpec, error) {
	conn := r.db.Open()
	defer conn.Close()
	var msgs []messageSpec

	err := conn.Where(&messageSpec{Exchange: exchange}).Find(&msgs).Error
	if err != nil {
		return []spec.MessageSpec{}, err
	}

	messageSpecs := []spec.MessageSpec{}

	for _, m := range msgs {
		msg := transformToMsg(m)
		messageSpecs = append(messageSpecs, msg)
	}

	return messageSpecs, nil
}

// Update a message spec
func (r *MessagesRepo) Update(msg spec.MessageSpec) error {
	conn := r.db.Open()
	defer conn.Close()

	m, err := transformMsg(msg)
	if err != nil {
		return err
	}

	return conn.Save(&m).Error
}
