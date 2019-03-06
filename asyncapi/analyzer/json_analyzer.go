package analyzer

import (
	"encoding/json"

	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
)

type JsonAnalyzer struct{}

const PayloadType = "object"

func (a JsonAnalyzer) GetPayloadSpec(payload []byte) spec.PayloadSpec {
	var parsed map[string]interface{}
	json.Unmarshal([]byte(payload), &parsed)

	var fields []spec.FieldSpec
	for k, v := range parsed {
		fields = append(fields, spec.FieldSpec{Name: k, Type: typeof(v)})
	}

	return spec.PayloadSpec{Fields: fields, Type: PayloadType}
}

func typeof(v interface{}) string {
	switch v.(type) {
	case float64:
		return "number"
	case string:
		return "string"
	case bool:
		return "boolean"
	default:
		return "unknown"
	}
}
