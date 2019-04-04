package analyzer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/scribano/asyncapi/repos/messages_repo"
	"github.com/uesteibar/scribano/asyncapi/spec"
	"github.com/uesteibar/scribano/consumer"
	"github.com/uesteibar/scribano/storage/db"
)

func treatTypeAsJSON(t *testing.T, contentType string) {
	chIn := make(chan consumer.Message)
	chOut := make(chan spec.MessageSpec)
	database := db.GetUniqueDB()
	repo := messages_repo.New(database)
	repo.Migrate()

	a := New(chIn, chOut, database)
	go a.Watch()

	sampleBody := `
		{
			"name": "infer type",
			"age": 27,
			"grade": 9.5,
			"canDrive": false,
			"birthDate": "1991-08-29",
			"lastLogin": "2015-06-10T13:23:30-08:00",
			"address": null,
			"emptyHash": {},
			"fines": [],
			"emptyHashes": [{}],
			"matrix": [
				[1, 2, 3],
				[3, 2, 1]
			],
			"friends": [
			  { "name": "pepe" },
			  { "name": "gotera" }
			],
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
		expectedFields := []*spec.FieldSpec{
			&spec.FieldSpec{
				Name: "name",
				Type: "string",
			},
			&spec.FieldSpec{
				Name: "age",
				Type: "integer",
			},
			&spec.FieldSpec{
				Name: "grade",
				Type: "number",
			},
			&spec.FieldSpec{
				Name: "canDrive",
				Type: "boolean",
			},
			&spec.FieldSpec{
				Name:   "birthDate",
				Type:   "string",
				Format: "date",
			},
			&spec.FieldSpec{
				Name:   "lastLogin",
				Type:   "string",
				Format: "date-time",
			},
			&spec.FieldSpec{
				Name: "address",
				Type: "string",
			},
			&spec.FieldSpec{
				Name: "emptyHash",
				Type: "object",
			},
			&spec.FieldSpec{
				Name: "fines",
				Type: "array",
				Item: &spec.FieldSpec{
					Type: "string",
				},
			},
			&spec.FieldSpec{
				Name: "emptyHashes",
				Type: "array",
				Item: &spec.FieldSpec{
					Type: "object",
				},
			},
			&spec.FieldSpec{
				Name: "matrix",
				Type: "array",
				Item: &spec.FieldSpec{
					Type: "array",
					Item: &spec.FieldSpec{
						Type: "integer",
					},
				},
			},
			&spec.FieldSpec{
				Name: "friends",
				Type: "array",
				Item: &spec.FieldSpec{
					Type: "object",
					Fields: []*spec.FieldSpec{
						&spec.FieldSpec{
							Name: "name",
							Type: "string",
						},
					},
				},
			},
			&spec.FieldSpec{
				Name: "car",
				Type: "object",
				Fields: []*spec.FieldSpec{
					&spec.FieldSpec{Name: "brand", Type: "string"},
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
	database := db.GetUniqueDB()
	repo := messages_repo.New(database)
	repo.Migrate()

	a := New(chIn, chOut, database)
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
	database := db.GetUniqueDB()
	repo := messages_repo.New(database)
	repo.Migrate()

	a := New(chIn, chOut, database)
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

func TestAnalyze_ChangingShape(t *testing.T) {
	chIn := make(chan consumer.Message)
	chOut := make(chan spec.MessageSpec)
	database := db.GetUniqueDB()
	repo := messages_repo.New(database)
	repo.Migrate()

	repo.Create(spec.MessageSpec{
		Topic:    "test.changing.shapes",
		Exchange: "/shapes",
		Payload: spec.PayloadSpec{
			Type: "object",
			Fields: []*spec.FieldSpec{
				&spec.FieldSpec{Name: "name", Type: "string"},
				&spec.FieldSpec{
					Name: "friends",
					Type: "array",
					Item: &spec.FieldSpec{
						Type: "object",
						Fields: []*spec.FieldSpec{
							&spec.FieldSpec{Name: "name", Type: "string"},
						},
					},
				},
				&spec.FieldSpec{
					Name:   "car",
					Type:   "object",
					Fields: []*spec.FieldSpec{&spec.FieldSpec{Name: "brand", Type: "string"}},
				},
			},
		},
	})

	a := New(chIn, chOut, database)

	go a.Watch()

	changedBody := `
		{
			"age": 27,
			"car": {
				"doors": 3
			},
			"friends": [
				{ "city": "barcelona" }
			]
		}
	`
	chIn <- consumer.Message{
		ContentType: "application/json",
		RoutingKey:  "test.changing.shapes",
		Body:        []byte(changedBody),
		Exchange:    "/shapes",
	}

	select {
	case res, _ := <-chOut:
		expectedFields := []*spec.FieldSpec{
			&spec.FieldSpec{Name: "name", Type: "string", Optional: true},
			&spec.FieldSpec{Name: "age", Type: "integer", Optional: true},
			&spec.FieldSpec{
				Name: "friends",
				Type: "array",
				Item: &spec.FieldSpec{
					Type: "object",
					Fields: []*spec.FieldSpec{
						&spec.FieldSpec{Name: "city", Type: "string", Optional: true},
						&spec.FieldSpec{Name: "name", Type: "string", Optional: true},
					},
				},
			},
			&spec.FieldSpec{
				Name: "car",
				Type: "object",
				Fields: []*spec.FieldSpec{
					&spec.FieldSpec{Name: "doors", Type: "integer", Optional: true},
					&spec.FieldSpec{Name: "brand", Type: "string", Optional: true},
				},
			},
		}

		assert.Equal(t, "test.changing.shapes", res.Topic)
		assert.Equal(t, "/shapes", res.Exchange)
		assert.Equal(t, "object", res.Payload.Type)

		assert.ElementsMatch(t, expectedFields, res.Payload.Fields)
	case <-time.After(3 * time.Second):
		t.Error("Expected to receive message, didn't receive any.")
	}
}
