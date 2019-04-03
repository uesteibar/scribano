package builder

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/scribano/asyncapi/spec"
)

func TestSpecBuilder(t *testing.T) {
	b := SpecBuilder{}

	b.AddServerInfo(spec.ServerSpec{})

	msg1 := spec.MessageSpec{
		Topic:    "some.topic",
		Exchange: "/some-exchange",
		Payload: spec.PayloadSpec{
			Type: "object",
			Fields: []spec.FieldSpec{
				spec.FieldSpec{Name: "name", Type: "string"},
				spec.FieldSpec{Name: "age", Type: "number"},
				spec.FieldSpec{Name: "emptyHash", Type: "object"},
				spec.FieldSpec{Name: "birthDate", Type: "string", Format: "date"},
				spec.FieldSpec{Name: "fines", Type: "array",
					Item: &spec.FieldSpec{
						Type: "string",
					},
				},
				spec.FieldSpec{Name: "emptyHashes", Type: "array",
					Item: &spec.FieldSpec{
						Type: "object",
					},
				},
				spec.FieldSpec{Name: "friends", Type: "array",
					Item: &spec.FieldSpec{
						Type: "object",
						Fields: []spec.FieldSpec{
							spec.FieldSpec{
								Name: "name",
								Type: "string",
							},
							spec.FieldSpec{
								Name:   "birthDate",
								Type:   "string",
								Format: "date",
							},
						},
					},
				},
				spec.FieldSpec{Name: "matrix", Type: "array",
					Item: &spec.FieldSpec{
						Type: "array",
						Item: &spec.FieldSpec{
							Type: "integer",
						},
					},
				},
				spec.FieldSpec{Name: "car", Type: "object",
					Fields: []spec.FieldSpec{
						spec.FieldSpec{Name: "brand", Type: "string"},
					},
				},
			},
		},
	}
	b.AddMessage(msg1)

	msg2 := spec.MessageSpec{
		Topic:    "other.topic",
		Exchange: "/other-exchange",
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
				Publish:  Ref{RefKey: "#/components/messages/SomeTopic"},
				Exchange: "/some-exchange",
			},
			"other.topic": Topic{
				Publish:  Ref{RefKey: "#/components/messages/OtherTopic"},
				Exchange: "/other-exchange",
			},
		},
		Components: Components{
			Messages: map[string]Message{
				"SomeTopic": Message{
					Payload: Payload{
						Type: "object",
						Properties: map[string]Property{
							"name":      Property{Type: "string"},
							"age":       Property{Type: "number"},
							"birthDate": Property{Type: "string", Format: "date"},
							"emptyHash": Property{Type: "object", Properties: map[string]Property{}},
							"fines": Property{Type: "array", Item: &Property{
								Type: "string",
							}},
							"emptyHashes": Property{Type: "array", Item: &Property{
								Type:       "object",
								Properties: map[string]Property{},
							}},
							"friends": Property{Type: "array", Item: &Property{
								Type: "object",
								Properties: map[string]Property{
									"birthDate": Property{Type: "string", Format: "date"},
									"name":      Property{Type: "string"},
								},
							}},
							"matrix": Property{Type: "array", Item: &Property{
								Type: "array",
								Item: &Property{
									Type: "integer",
								},
							}},
							"car": Property{
								Type: "object",
								Properties: map[string]Property{
									"brand": Property{Type: "string"},
								},
							},
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
				},
				"x-exchange": "/other-exchange"
			},
			"some.topic": {
				"publish": {
				  "$ref": "#/components/messages/SomeTopic"
				},
				"x-exchange": "/some-exchange"
			}
		},
		"components": {
			"messages": {
				"OtherTopic": {
					"payload": {
						"type":"object",
						"properties": {
							"uuid": {
								"type": "string"
							}
						}
					}
				},
				"SomeTopic": {
					"payload": {
						"type":"object",
						"properties": {
							"age": {
								"type": "number"
							},
							"birthDate": {
								"type": "string",
								"format": "date"
							},
							"car": {
								"type": "object",
								"properties": {
									"brand": {
										"type": "string"
									}
								}
							},
							"emptyHash": {
								"type": "object"
							},
							"emptyHashes": {
								"type": "array",
								"items": {
									"type": "object"
								}
							},
							"fines": {
								"type": "array",
								"items": {
									"type": "string"
								}
							},
							"friends": {
								"type": "array",
								"items": {
									"type": "object",
									"properties": {
										"birthDate": {
											"type": "string",
											"format": "date"
										},
										"name": {
											"type": "string"
										}
									}
								}
							},
							"matrix": {
								"type": "array",
								"items": {
									"type": "array",
									"items": {
										"type": "integer"
									}
								}
							},
							"name": {
								"type": "string"
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
