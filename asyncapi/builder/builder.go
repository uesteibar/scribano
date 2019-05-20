package builder

import (
	"fmt"
	"strings"

	"github.com/uesteibar/scribano/asyncapi/spec"
)

type property struct {
	Type        string              `json:"type"`
	Format      string              `json:"format,omitempty"`
	Description string              `json:"description,omitempty"`
	Properties  map[string]property `json:"properties,omitempty"`
	Item        *property           `json:"items,omitempty"`
	Optional    bool                `json:"x-optional,omitempty"`
}

type payload struct {
	Type       string              `json:"type"`
	Properties map[string]property `json:"properties"`
}

type message struct {
	Payload payload `json:"payload"`
}

type components struct {
	Messages map[string]message `json:"messages"`
}

type ref struct {
	RefKey string `json:"$ref"`
}

type topic struct {
	Publish  ref    `json:"publish"`
	Exchange string `json:"x-exchange"`
}

type information struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

// AsyncAPISpec represent the asyncapi data structure
type AsyncAPISpec struct {
	AsyncAPI   string           `json:"asyncapi"`
	Info       information      `json:"info"`
	Topics     map[string]topic `json:"topics"`
	Components components       `json:"components"`
}

// SpecBuilder builds a spec using the builder pattern
type SpecBuilder struct {
	Spec AsyncAPISpec
}

func buildArrayItem(f *spec.FieldSpec) *property {
	p := &property{Type: f.Type}
	if f.Type == "object" {
		p.Properties = buildProperties(f.Fields)
	} else if f.Type == "array" {
		p.Item = buildArrayItem(f.Item)
	}

	return p
}

const optionalDescription = "Optional field"

func buildProperties(fields []*spec.FieldSpec) map[string]property {
	properties := make(map[string]property)

	for _, f := range fields {
		p := property{Type: f.Type, Format: f.Format, Optional: f.Optional}
		if p.Optional {
			p.Description = optionalDescription
		}
		if f.Type == "object" {
			p.Properties = buildProperties(f.Fields)
		} else if f.Type == "array" {
			p.Item = buildArrayItem(f.Item)
		}

		properties[f.Name] = p
	}

	return properties
}

func buildMsg(msg spec.MessageSpec) message {
	m := message{Payload: payload{Type: msg.Payload.Type}}

	m.Payload.Properties = buildProperties(msg.Payload.Fields)

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

func refFor(msg spec.MessageSpec) ref {
	return ref{RefKey: fmt.Sprintf("#/components/messages/%s", msgName(msg))}
}

// AddMessage inserts the information for a message in the asyncapi spec
func (b *SpecBuilder) AddMessage(msg spec.MessageSpec) *SpecBuilder {
	if b.Spec.Topics == nil {
		b.Spec.Topics = make(map[string]topic)
	}
	b.Spec.Topics[msg.Topic] = topic{Publish: refFor(msg), Exchange: msg.Exchange}

	if b.Spec.Components.Messages == nil {
		b.Spec.Components.Messages = make(map[string]message)
	}
	b.Spec.Components.Messages[msgName(msg)] = buildMsg(msg)

	return b
}

const asyncAPIVersion = "1.0.0"

// AddServerInfo adds the server information to the asyncapi spec
func (b *SpecBuilder) AddServerInfo(info spec.ServerSpec) *SpecBuilder {
	b.Spec.Info = information{Title: info.Name, Version: info.Version}
	b.Spec.AsyncAPI = asyncAPIVersion

	return b
}

// Build builds the final asyncapi spec
func (b *SpecBuilder) Build() AsyncAPISpec {
	if b.Spec.Topics == nil {
		b.Spec.Topics = make(map[string]topic)
	}
	if b.Spec.Components.Messages == nil {
		b.Spec.Components.Messages = make(map[string]message)
	}

	return b.Spec
}
