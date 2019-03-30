package builder

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
)

func TestSpecBuilder(t *testing.T) {
	b := SpecBuilder{}

	b.AddServerInfo(spec.ServerSpec{})

	msg1 := spec.MessageSpec{
		Topic: "some.topic",
		Payload: spec.PayloadSpec{
			Type: "object",
			Fields: []spec.FieldSpec{
				spec.FieldSpec{Name: "name", Type: "string"},
				spec.FieldSpec{Name: "age", Type: "number"},
			},
		},
	}
	b.AddMessage(msg1)

	msg2 := spec.MessageSpec{
		Topic: "other.topic",
		Payload: spec.PayloadSpec{
			Type: "object",
			Fields: []spec.FieldSpec{
				spec.FieldSpec{Name: "uuid", Type: "string"},
			},
		},
	}
	b.AddMessage(msg2)

	res := b.Build()

	expected := AsyncAPISpec{
		AsyncAPI: "1.0.0",
		Info: Info{
			Title:   "",
			Version: "",
		},
		Topics: map[string]Topic{
			"some.topic": Topic{
				Publish: Ref{RefKey: "#/components/messages/SomeTopic"},
			},
			"other.topic": Topic{
				Publish: Ref{RefKey: "#/components/messages/OtherTopic"},
			},
		},
		Components: Components{
			Messages: map[string]Message{
				"SomeTopic": Message{
					Payload: Payload{
						Type: "object",
						Properties: map[string]Property{
							"name": Property{Type: "string"},
							"age":  Property{Type: "number"},
						},
					},
				},
				"OtherTopic": Message{
					Payload: Payload{
						Type: "object",
						Properties: map[string]Property{
							"uuid": Property{Type: "string"},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, res)

	json, _ := json.Marshal(res)

	expectedJSON := `{
		"asyncapi": "1.0.0",
		"info": {
			"title": "",
			"version": ""
		},
		"topics": {
			"other.topic": {
				"publish": {
				  "$ref": "#/components/messages/OtherTopic"
				}
			},
			"some.topic": {
				"publish": {
				  "$ref": "#/components/messages/SomeTopic"
				}
			}
		},
		"components": {
			"messages": {
				"OtherTopic": {
					"payload": {
						"type":"object",
						"properties": {
							"uuid": {
								"type":"string"
							}
						}
					}
				},
				"SomeTopic": {
					"payload": {
						"type":"object",
						"properties": {
							"age": {
								"type":"number"
							},
							"name": {
								"type":"string"
							}
						}
					}
				}
			}
		}
	}`

	expectedJSON = strings.Replace(expectedJSON, "\n", "", -1)
	expectedJSON = strings.Replace(expectedJSON, "\t", "", -1)
	expectedJSON = strings.Replace(expectedJSON, " ", "", -1)

	assert.Equal(t, expectedJSON, string(json))
}
