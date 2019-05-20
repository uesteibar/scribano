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
			Fields: []*spec.FieldSpec{
				&spec.FieldSpec{Name: "name", Type: "string"},
				&spec.FieldSpec{Name: "age", Type: "number"},
				&spec.FieldSpec{Name: "emptyHash", Type: "object"},
				&spec.FieldSpec{Name: "birthDate", Type: "string", Format: "date"},
				&spec.FieldSpec{Name: "fines", Type: "array",
					Item: &spec.FieldSpec{
						Type: "string",
					},
				},
				&spec.FieldSpec{Name: "emptyHashes", Type: "array",
					Item: &spec.FieldSpec{
						Type: "object",
					},
				},
				&spec.FieldSpec{Name: "friends", Type: "array",
					Item: &spec.FieldSpec{
						Type: "object",
						Fields: []*spec.FieldSpec{
							&spec.FieldSpec{
								Name: "name",
								Type: "string",
							},
							&spec.FieldSpec{
								Name:     "birthDate",
								Type:     "string",
								Format:   "date",
								Optional: true,
							},
						},
					},
				},
				&spec.FieldSpec{Name: "matrix", Type: "array",
					Item: &spec.FieldSpec{
						Type: "array",
						Item: &spec.FieldSpec{
							Type: "integer",
						},
					},
				},
				&spec.FieldSpec{Name: "car", Type: "object",
					Fields: []*spec.FieldSpec{
						&spec.FieldSpec{Name: "brand", Type: "string"},
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
			Fields: []*spec.FieldSpec{
				&spec.FieldSpec{Name: "uuid", Type: "string"},
			},
		},
	}
	b.AddMessage(msg2)

	res := b.Build()

	expected := AsyncAPISpec{
		AsyncAPI: "1.0.0",
		Info: information{
			Title:   "",
			Version: "",
		},
		Topics: map[string]topic{
			"some.topic": topic{
				Publish:  ref{RefKey: "#/components/messages/SomeTopic"},
				Exchange: "/some-exchange",
			},
			"other.topic": topic{
				Publish:  ref{RefKey: "#/components/messages/OtherTopic"},
				Exchange: "/other-exchange",
			},
		},
		Components: components{
			Messages: map[string]message{
				"SomeTopic": message{
					Payload: payload{
						Type: "object",
						Properties: map[string]property{
							"name":      property{Type: "string"},
							"age":       property{Type: "number"},
							"birthDate": property{Type: "string", Format: "date"},
							"emptyHash": property{Type: "object", Properties: map[string]property{}},
							"fines": property{Type: "array", Item: &property{
								Type: "string",
							}},
							"emptyHashes": property{Type: "array", Item: &property{
								Type:       "object",
								Properties: map[string]property{},
							}},
							"friends": property{Type: "array", Item: &property{
								Type: "object",
								Properties: map[string]property{
									"birthDate": property{Type: "string", Format: "date", Optional: true, Description: "Optional field"},
									"name":      property{Type: "string"},
								},
							}},
							"matrix": property{Type: "array", Item: &property{
								Type: "array",
								Item: &property{
									Type: "integer",
								},
							}},
							"car": property{
								Type: "object",
								Properties: map[string]property{
									"brand": property{Type: "string"},
								},
							},
						},
					},
				},
				"OtherTopic": message{
					Payload: payload{
						Type: "object",
						Properties: map[string]property{
							"uuid": property{Type: "string"},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, res)

	json, _ := json.Marshal(res)

	expectedJSON := `{
		"asyncapi":"1.0.0",
		"info":{
			"title":"",
			"version":""
		},
		"topics":{
			"other.topic":{
				"publish":{
					"$ref":"#/components/messages/OtherTopic"
				},
				"x-exchange":"/other-exchange"
			},
			"some.topic":{
				"publish":{
					"$ref":"#/components/messages/SomeTopic"
				},
				"x-exchange":"/some-exchange"
			}
		},
		"components":{
			"messages":{
				"OtherTopic":{
					"payload":{
						"type":"object",
						"properties":{
							"uuid":{
								"type":"string"
							}
						}
					}
				},
				"SomeTopic":{
					"payload":{
						"type":"object",
						"properties":{
							"age":{
								"type":"number"
							},
							"birthDate":{
								"type":"string",
								"format":"date"
							},
							"car":{
								"type":"object",
								"properties":{
									"brand":{
										"type":"string"
									}
								}
							},
							"emptyHash":{
								"type":"object"
							},
							"emptyHashes":{
								"type":"array",
								"items":{
									"type":"object"
								}
							},
							"fines":{
								"type":"array",
								"items":{
									"type":"string"
								}
							},
							"friends":{
								"type":"array",
								"items":{
									"type":"object",
									"properties":{
										"birthDate":{
											"type":"string",
											"format":"date",
											"description":"Optional field",
											"x-optional":true
										},
										"name":{
											"type":"string"
										}
									}
								}
							},
							"matrix":{
								"type":"array",
								"items":{
									"type":"array",
									"items":{
										"type":"integer"
									}
								}
							},
							"name":{
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

	assert.Equal(t, expectedJSON, string(json))
}

func TestSpecBuilder_NoMessages(t *testing.T) {
	b := SpecBuilder{}

	b.AddServerInfo(spec.ServerSpec{})

	res := b.Build()

	expected := AsyncAPISpec{
		AsyncAPI: "1.0.0",
		Info: information{
			Title:   "",
			Version: "",
		},
		Topics: map[string]topic{},
		Components: components{
			Messages: map[string]message{},
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
		"topics": {},
		"components": {
			"messages": {}
		}
	}`

	expectedJSON = strings.Replace(expectedJSON, "\n", "", -1)
	expectedJSON = strings.Replace(expectedJSON, "\t", "", -1)
	expectedJSON = strings.Replace(expectedJSON, " ", "", -1)

	assert.Equal(t, expectedJSON, string(json))
}
