package spec

import (
	"github.com/streadway/amqp"
)

// MessageSpec represents a message in the asyncapi spec
type MessageSpec struct {
	Topic    string
	Exchange string
	Payload  PayloadSpec
	Delivery amqp.Delivery
}

// PayloadSpec represents a payload (message body) in the asyncapi spec
type PayloadSpec struct {
	Type   string       `json:"type"`
	Fields []*FieldSpec `json:"fields"`
}

// FieldSpec represents a field of a payload in the asyncapi spec
type FieldSpec struct {
	Name     string       `json:"name"`
	Type     string       `json:"type"`
	Format   string       `json:"format"`
	Fields   []*FieldSpec `json:"fields"`
	Item     *FieldSpec   `json:"Item"`
	Optional bool         `json:"optional"`
}

// ServerSpec represents a server in the asyncapi spec
type ServerSpec struct {
	Name    string
	Version string
}
