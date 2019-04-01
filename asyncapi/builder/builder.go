package builder

import (
	"fmt"
	"strings"

	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
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

type Ref struct {
	RefKey string `json:"$ref"`
}

type Topic struct {
	Publish  Ref    `json:"publish"`
	Exchange string `json:"x-exchange"`
}

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type AsyncAPISpec struct {
	AsyncAPI   string           `json:"asyncapi"`
	Info       Info             `json:"info"`
	Topics     map[string]Topic `json:"topics"`
	Components Components       `json:"components"`
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

func refFor(msg spec.MessageSpec) Ref {
	return Ref{RefKey: fmt.Sprintf("#/components/messages/%s", msgName(msg))}
}

func (b *SpecBuilder) AddMessage(msg spec.MessageSpec) *SpecBuilder {
	if b.Spec.Topics == nil {
		b.Spec.Topics = make(map[string]Topic)
	}
	b.Spec.Topics[msg.Topic] = Topic{Publish: refFor(msg), Exchange: msg.Exchange}

	if b.Spec.Components.Messages == nil {
		b.Spec.Components.Messages = make(map[string]Message)
	}
	b.Spec.Components.Messages[msgName(msg)] = buildMsg(msg)

	return b
}

const asyncApiVersion = "1.0.0"

func (b *SpecBuilder) AddServerInfo(info spec.ServerSpec) *SpecBuilder {
	b.Spec.Info = Info{Title: info.Name, Version: info.Version}
	b.Spec.AsyncAPI = asyncApiVersion

	return b
}

func (b *SpecBuilder) Build() AsyncAPISpec {
	return b.Spec
}
