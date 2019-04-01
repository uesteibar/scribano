package analyzer

import (
	"encoding/json"

	"github.com/uesteibar/scribano/asyncapi/spec"
)

// JSONAnalyzer analyzes json payloads to build a spec
type JSONAnalyzer struct{}

const (
	payloadType = "object"
	unknownType = "unknown"
)

// GetPayloadSpec analyzes a payload and returns the spec
func (a JSONAnalyzer) GetPayloadSpec(payload []byte) spec.PayloadSpec {
	var parsed map[string]interface{}
	json.Unmarshal([]byte(payload), &parsed)

	var fields []spec.FieldSpec
	for k, v := range parsed {
		fields = append(fields, spec.FieldSpec{Name: k, Type: typeof(v)})
	}

	return spec.PayloadSpec{Fields: fields, Type: payloadType}
}

func isRound(n float64) bool {
	return n == float64(int64(n))
}

func inferNumberType(n interface{}) string {
	if isRound(n.(float64)) {
		return "integer"
	}

	return "number"
}

func typeof(v interface{}) string {
	switch v.(type) {
	case float64:
		return inferNumberType(v)
	case string:
		return "string"
	case bool:
		return "boolean"
	default:
		return unknownType
	}
}
