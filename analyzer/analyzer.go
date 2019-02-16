package analyzer

import (
	"errors"
	"fmt"
	"github.com/uesteibar/asyncapi-watcher/consumer"
)

const JsonContentType = "application/json"

type MessageAnalyzer interface {
	BuildSpec(consumer.Message) MessageSpec
}

type MessageSpec struct {
	Topic  string
	Fields []FieldSpec
}

type FieldSpec struct {
	Name string
	Type string
}

type Analyzer struct {
	ChIn  chan consumer.Message
	ChOut chan MessageSpec
}

func getAnalyzer(msg consumer.Message) (MessageAnalyzer, error) {
	switch msg.ContentType {
	case JsonContentType:
		return JsonAnalyzer{}, nil
	default:
		return nil, errors.New(
			fmt.Sprintf("ContentType not supported: %s", msg.ContentType),
		)
	}
}

func analyze(msg consumer.Message) (MessageSpec, error) {
	a, err := getAnalyzer(msg)
	if err != nil {
		return MessageSpec{}, errors.New(
			fmt.Sprintf("Couldn't analyze message: %+v", msg),
		)
	}
	return a.BuildSpec(msg), nil
}

func (a *Analyzer) Watch() {
	for msg := range a.ChIn {
		spec, err := analyze(msg)
		if err == nil {
			a.ChOut <- spec
		}
	}

}
