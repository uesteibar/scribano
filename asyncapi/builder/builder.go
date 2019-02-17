package builder

import (
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"strings"
)

type Property struct {
	Type string `json:"type"`
}

type Payload struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
}

type Message struct {
	Payload Payload `json:"payload"`
}

type Components struct {
	Messages map[string]Message `json:"messages"`
}

type AsyncAPISpec struct {
	Components Components `json:"components"`
}

type SpecBuilder struct {
	Spec AsyncAPISpec
}

func buildMsg(msg spec.MessageSpec) Message {
	m := Message{
		Payload: Payload{
			Type: msg.Payload.Type,
		},
	}
	m.Payload.Properties = make(map[string]Property)
	for _, f := range msg.Payload.Fields {
		m.Payload.Properties[f.Name] = Property{
			Type: f.Type,
		}
	}

	return m
}

func msgName(msg spec.MessageSpec) string {
	split := strings.Split(msg.Topic, ".")

	var pieces []string
	for _, p := range split {
		pieces = append(pieces, strings.Title(p))
	}

	return strings.Join(pieces, "")
}

func (b *SpecBuilder) AddMessage(msg spec.MessageSpec) *SpecBuilder {
	if b.Spec.Components.Messages == nil {
		b.Spec.Components.Messages = make(map[string]Message)
	}

	b.Spec.Components.Messages[msgName(msg)] = buildMsg(msg)
	return b
}

func (b *SpecBuilder) Build() AsyncAPISpec {
	return b.Spec
}