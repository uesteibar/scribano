package builder

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"testing"
)

func TestSpecBuilder(t *testing.T) {
	b := SpecBuilder{}

	msg1 := spec.MessageSpec{
		Topic: "some.topic",
		Payload: spec.PayloadSpec{
			Type: "object",
			Fields: []spec.FieldSpec{
				spec.FieldSpec{Name: "name", Type: "string"},
				spec.FieldSpec{Name: "age", Type: "float64"},
			},
		},
	}
	b.AddMessage(msg1)

	res := b.Build()

	expected := AsyncAPISpec{
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
								Type: "float64",
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, res)

	json, _ := json.Marshal(res)

	expectedJson := `{"components":{"messages":{"SomeTopic":{"payload":{"type":"object","properties":{"age":{"type":"float64"},"name":{"type":"string"}}}}}}}`

	assert.Equal(t, expectedJson, string(json))
}
