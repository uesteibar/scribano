package analyzer

import (
	"encoding/json"
	"github.com/uesteibar/asyncapi-watcher/consumer"
)

type JsonAnalyzer struct{}

func (a JsonAnalyzer) BuildSpec(msg consumer.Message) MessageSpec {
	var parsed map[string]interface{}
	json.Unmarshal([]byte(msg.Body), &parsed)

	var fields []FieldSpec
	for k, v := range parsed {
		fields = append(fields, FieldSpec{Name: k, Type: typeof(v)})
	}

	return MessageSpec{Fields: fields, Topic: msg.RoutingKey}
}

func typeof(v interface{}) string {
	switch v.(type) {
	case float64:
		return "float64"
	case string:
		return "string"
	case bool:
		return "boolean"
	default:
		return "unknown"
	}
}
