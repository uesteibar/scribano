package analyzer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/scribano/asyncapi/spec"
	"github.com/uesteibar/scribano/consumer"
)

func treatTypeAsJSON(t *testing.T, contentType string) {
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
			"address": null,
			"car": {
				"brand": "mercedes"
			}
		}
	`
	chIn <- consumer.Message{
		ContentType: contentType,
		RoutingKey:  "test.routing.key",
		Body:        []byte(sampleBody),
		Exchange:    "/my-exchange",
	}
	select {
	case res, _ := <-chOut:
		expectedFields := []spec.FieldSpec{
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
				Name: "address",
				Type: "string",
			},
			spec.FieldSpec{
				Name: "car",
				Type: "object",
				Fields: []spec.FieldSpec{
					spec.FieldSpec{Name: "brand", Type: "string"},
				},
			},
		}

		assert.Equal(t, "test.routing.key", res.Topic)
		assert.Equal(t, "/my-exchange", res.Exchange)
		assert.Equal(t, "object", res.Payload.Type)

		assert.ElementsMatch(t, expectedFields, res.Payload.Fields)
	case <-time.After(3 * time.Second):
		t.Error("Expected to receive message, didn't receive any.")
	}
}

func TestAnalyze_JSON(t *testing.T) {
	treatTypeAsJSON(t, "application/json")
}

func TestAnalyze_OctetStream(t *testing.T) {
	treatTypeAsJSON(t, "application/octet-stream")
}

func TestAnalyze_JSON_InvalidContent(t *testing.T) {
	chIn := make(chan consumer.Message)
	chOut := make(chan spec.MessageSpec)

	a := Analyzer{ChIn: chIn, ChOut: chOut}
	go a.Watch()

	chIn <- consumer.Message{
		ContentType: "application/octet-stream",
		RoutingKey:  "test.routing.key",
		Body:        []byte("invalid body"),
		Exchange:    "/",
	}

	select {
	case res, _ := <-chOut:
		t.Errorf("Expected to not receive message, received: %+v", res)
	case <-time.After(time.Second):
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
