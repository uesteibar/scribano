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

	b.AddServerInfo(spec.ServerSpec{
		Name:    "Test_title",
		Version: "0.0.1",
	})

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
			Title:   "Test_title",
			Version: "0.0.1",
		},
		Topics: map[string]Topic{
			"some.topic": Topic{
				Subscribe: Ref{RefKey: "#/components/messages/SomeTopic"},
				Publish:   Ref{RefKey: "#/components/messages/SomeTopic"},
			},
			"other.topic": Topic{
				Subscribe: Ref{RefKey: "#/components/messages/OtherTopic"},
				Publish:   Ref{RefKey: "#/components/messages/OtherTopic"},
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

	expectedJson := `{
		"asyncapi": "1.0.0",
		"info": {
			"title": "Test_title",
			"version": "0.0.1"
		},
		"topics": {
			"other.topic": {
				"subscribe": {
				  "$ref": "#/components/messages/OtherTopic"
				},
				"publish": {
				  "$ref": "#/components/messages/OtherTopic"
				}
			},
			"some.topic": {
				"subscribe": {
				  "$ref": "#/components/messages/SomeTopic"
				},
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

	expectedJson = strings.Replace(expectedJson, "\n", "", -1)
	expectedJson = strings.Replace(expectedJson, "\t", "", -1)
	expectedJson = strings.Replace(expectedJson, " ", "", -1)

	assert.Equal(t, expectedJson, string(json))
}
