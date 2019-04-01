package analyzer

import (
	"fmt"

	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/consumer"
)

const jsonContentType = "application/json"

// PayloadAnalyzer objects analyze payloads
type PayloadAnalyzer interface {
	GetPayloadSpec([]byte) spec.PayloadSpec
}

// Analyzer analyzes incoming messages and pipes the result through
type Analyzer struct {
	ChIn  chan consumer.Message
	ChOut chan spec.MessageSpec
}

func getPayloadAnalyzer(msg consumer.Message) (PayloadAnalyzer, error) {
	switch msg.ContentType {
	case jsonContentType:
		return JSONAnalyzer{}, nil
	default:
		return nil, fmt.Errorf(
			fmt.Sprintf("ContentType not supported: %s", msg.ContentType),
		)
	}
}

func analyze(msg consumer.Message) (spec.MessageSpec, error) {
	a, err := getPayloadAnalyzer(msg)
	if err != nil {
		return spec.MessageSpec{}, fmt.Errorf(
			fmt.Sprintf("Couldn't analyze message: %+v", msg),
		)
	}
	return spec.MessageSpec{
		Topic:   msg.RoutingKey,
		Payload: a.GetPayloadSpec(msg.Body),
	}, nil
}

// Watch for incoming messages
func (a *Analyzer) Watch() {
	for msg := range a.ChIn {
		spec, err := analyze(msg)
		if err == nil {
			a.ChOut <- spec
		}
	}

}
