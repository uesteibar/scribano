package analyzer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/scribano/asyncapi/spec"
	"github.com/uesteibar/scribano/consumer"
)

func TestAnalyze_JSON(t *testing.T) {
	chIn := make(chan consumer.Message)
	chOut := make(chan spec.MessageSpec)

	a := Analyzer{ChIn: chIn, ChOut: chOut}
	go a.Watch()

	sampleBody := `
		{
			"name": "infer type",
			"age": 27,
			"grade": 9.5,
			"canDrive": false,
			"car": {
				"brand": "mercedes",
				"seats": 5
			}
		}
	`
	chIn <- consumer.Message{
		ContentType: "application/json",
		RoutingKey:  "test.routing.key",
		Body:        []byte(sampleBody),
		Exchange:    "/my-exchange",
	}
	select {
	case res, _ := <-chOut:
		expected := spec.MessageSpec{
			Topic:    "test.routing.key",
			Exchange: "/my-exchange",
			Payload: spec.PayloadSpec{
				Type: "object",
				Fields: []spec.FieldSpec{
					spec.FieldSpec{
						Name: "name",
						Type: "string",
					},
					spec.FieldSpec{
						Name: "age",
						Type: "integer",
					},
					spec.FieldSpec{
						Name: "grade",
						Type: "number",
					},
					spec.FieldSpec{
						Name: "canDrive",
						Type: "boolean",
					},
					spec.FieldSpec{
						Name: "car",
						Type: "object",
						Fields: []spec.FieldSpec{
							spec.FieldSpec{
								Name: "brand",
								Type: "string",
							},
							spec.FieldSpec{
								Name: "seats",
								Type: "integer",
							},
						},
					},
				},
			},
		}

		assert.Equal(t, expected, res)
	case <-time.After(3 * time.Second):
		t.Error("Expected to receive message, didn't receive any.")
	}

}

func TestAnalyze_UnknownFormat(t *testing.T) {
	chIn := make(chan consumer.Message)
	chOut := make(chan spec.MessageSpec)

	a := Analyzer{ChIn: chIn, ChOut: chOut}
	go a.Watch()

	chIn <- consumer.Message{
		ContentType: "plain/text",
		RoutingKey:  "test.routing.key",
		Body:        []byte("plain body"),
		Exchange:    "/",
	}

	select {
	case res, _ := <-chOut:
		t.Errorf("Expected to not receive message, received: %+v", res)
	case <-time.After(time.Second):
	}
}
