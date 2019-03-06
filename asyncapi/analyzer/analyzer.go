package analyzer

import (
	"errors"
	"fmt"

	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/consumer"
)

const JsonContentType = "application/json"

type PayloadAnalyzer interface {
	GetPayloadSpec([]byte) spec.PayloadSpec
}

type Analyzer struct {
	ChIn  chan consumer.Message
	ChOut chan spec.MessageSpec
}

func getPayloadAnalyzer(msg consumer.Message) (PayloadAnalyzer, error) {
	switch msg.ContentType {
	case JsonContentType:
		return JsonAnalyzer{}, nil
	default:
		return nil, errors.New(
			fmt.Sprintf("ContentType not supported: %s", msg.ContentType),
		)
	}
}

func analyze(msg consumer.Message) (spec.MessageSpec, error) {
	a, err := getPayloadAnalyzer(msg)
	if err != nil {
		return spec.MessageSpec{}, errors.New(
			fmt.Sprintf("Couldn't analyze message: %+v", msg),
		)
	}
	return spec.MessageSpec{
		Topic:   msg.RoutingKey,
		Payload: a.GetPayloadSpec(msg.Body),
	}, nil
}

func (a *Analyzer) Watch() {
	for msg := range a.ChIn {
		spec, err := analyze(msg)
		if err == nil {
			a.ChOut <- spec
		}
	}

}
