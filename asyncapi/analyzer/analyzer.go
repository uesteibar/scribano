package analyzer

import (
	"fmt"
	"log"

	"github.com/uesteibar/scribano/asyncapi/spec"
	"github.com/uesteibar/scribano/consumer"
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
	ChIn  chan consumer.Message
	ChOut chan spec.MessageSpec
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

func analyze(msg consumer.Message) (spec.MessageSpec, error) {
	a, err := getPayloadAnalyzer(msg)
	if err != nil {
		return spec.MessageSpec{}, err
	}

	p, err := a.GetPayloadSpec(msg.Body)
	if err != nil {
		return spec.MessageSpec{}, err
	}
	return spec.MessageSpec{
		Topic:    msg.RoutingKey,
		Exchange: msg.Exchange,
		Payload:  p,
	}, nil
}

// Watch for incoming messages
func (a *Analyzer) Watch() {
	for msg := range a.ChIn {
		spec, err := analyze(msg)
		if err != nil {
			log.Printf("ERROR %s", err.Error())
		} else {
			a.ChOut <- spec
		}

	}

}
