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
		},
		Components: Components{
			Messages: map[string]Message{
				"SomeTopic": Message{
					Payload: Payload{
						Type: "object",
						Properties: map[string]Property{
							"name": Property{
								Type: "string",
							},
							"age": Property{
								Type: "number",
							},
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
