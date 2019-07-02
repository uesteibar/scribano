package analyzer

import (
	"fmt"
	"log"

	"github.com/uesteibar/scribano/asyncapi/repos/messagesrepo"
	"github.com/uesteibar/scribano/asyncapi/spec"
	"github.com/uesteibar/scribano/consumer"
	"github.com/uesteibar/scribano/storage/db"
)

const (
	jsonContentType        = "application/json"
	octetStreamContentType = "application/octet-stream"
)

// PayloadAnalyzer objects analyze payloads
type PayloadAnalyzer interface {
	GetPayloadSpec([]byte) (spec.PayloadSpec, error)
}

// Analyzer analyzes incoming messages and pipes the result through
type Analyzer struct {
	ChIn         chan consumer.Message
	ChOut        chan spec.MessageSpec
	messagesRepo *messagesrepo.MessagesRepo
}

// New creates a new Analyzer
func New(chIn chan consumer.Message, chOut chan spec.MessageSpec, database db.Database) *Analyzer {
	return &Analyzer{
		ChIn:         chIn,
		ChOut:        chOut,
		messagesRepo: messagesrepo.New(database),
	}
}

func getPayloadAnalyzer(msg consumer.Message) (PayloadAnalyzer, error) {
	switch msg.ContentType {
	case jsonContentType, octetStreamContentType:
		return JSONAnalyzer{}, nil
	default:
		return nil, fmt.Errorf(
			fmt.Sprintf("ContentType not supported: %s", msg.ContentType),
		)
	}
}

func (a *Analyzer) buildModel(msg consumer.Message) (spec.MessageSpec, error) {
	pa, err := getPayloadAnalyzer(msg)
	if err != nil {
		return spec.MessageSpec{}, err
	}

	p, err := pa.GetPayloadSpec(msg.Body)
	if err != nil {
		return spec.MessageSpec{}, err
	}
	return spec.MessageSpec{
		Topic:    msg.RoutingKey,
		Exchange: msg.Exchange,
		Payload:  p,
		Delivery: msg.Delivery,
	}, nil
}

func findField(target *spec.FieldSpec, fields []*spec.FieldSpec) (*spec.FieldSpec, bool) {
	for _, f := range fields {
		if f.Name == target.Name {
			return f, true
		}
	}

	return &spec.FieldSpec{}, false
}

func flagNewFields(nextFields, currentFields []*spec.FieldSpec) []*spec.FieldSpec {
	for _, f := range nextFields {
		currentField, found := findField(f, currentFields)
		if !found {
			f.Optional = true
		} else {
			f.Fields = flagNewFields(f.Fields, currentField.Fields)
			if f.Item != nil && currentField.Item != nil {
				f.Item.Fields = flagNewFields(f.Item.Fields, currentField.Item.Fields)
			}
		}
	}

	return nextFields
}

func keepMissingFields(nextFields, currentFields []*spec.FieldSpec) []*spec.FieldSpec {
	for _, f := range currentFields {
		nextField, found := findField(f, nextFields)
		if !found {
			f.Optional = true
			nextFields = append(nextFields, f)
		} else {
			nextField.Fields = keepMissingFields(nextField.Fields, f.Fields)

			if nextField.Item != nil && f.Item != nil {
				nextField.Item.Fields = keepMissingFields(nextField.Item.Fields, f.Item.Fields)
			}
		}
	}

	return nextFields
}

func merge(next *spec.MessageSpec, current spec.MessageSpec) {
	next.Payload.Fields = flagNewFields(
		next.Payload.Fields, current.Payload.Fields,
	)
	next.Payload.Fields = keepMissingFields(
		next.Payload.Fields, current.Payload.Fields,
	)
}

func (a *Analyzer) analyzeCompared(next *spec.MessageSpec) error {
	current, err := a.messagesRepo.Find(next.Topic)
	if err != nil {
		_, isNotFound := err.(*messagesrepo.ErrNotFound)
		if isNotFound {
			return nil
		}

		return err
	}

	merge(next, current)

	return nil
}

func (a *Analyzer) analyze(msg consumer.Message) (spec.MessageSpec, error) {
	model, err := a.buildModel(msg)
	if err != nil {
		return model, err
	}

	err = a.analyzeCompared(&model)
	if err != nil {
		return model, err
	}

	return model, nil
}

// Watch for incoming messages
func (a *Analyzer) Watch() {
	for msg := range a.ChIn {
		spec, err := a.analyze(msg)
		if err != nil {
			log.Printf("ERROR %s", err.Error())
		} else {
			a.ChOut <- spec
		}
	}

	log.Printf("INFO finished running analyzer")
}
