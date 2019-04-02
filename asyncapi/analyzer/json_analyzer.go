package analyzer

import (
	"encoding/json"

	"github.com/uesteibar/scribano/asyncapi/spec"
)

// JSONAnalyzer analyzes json payloads to build a spec
type JSONAnalyzer struct{}

const (
	booleanType = "boolean"
	integerType = "integer"
	numberType  = "number"
	objectType  = "object"
	stringType  = "string"
	unknownType = "unknown"
)

// GetPayloadSpec analyzes a payload and returns the spec
func (a JSONAnalyzer) GetPayloadSpec(payload []byte) (spec.PayloadSpec, error) {
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(payload), &parsed)

	if err != nil {
		return spec.PayloadSpec{}, err
	}

	fields := fieldsFor(parsed)

	return spec.PayloadSpec{Fields: fields, Type: objectType}, nil
}

func fieldsFor(raw map[string]interface{}) []spec.FieldSpec {
	var fields []spec.FieldSpec
	for k, v := range raw {
		fields = append(fields, fieldFor(k, v))
	}

	return fields
}

func isRound(n float64) bool {
	return n == float64(int64(n))
}

func inferNumberType(n interface{}) string {
	if isRound(n.(float64)) {
		return integerType
	}

	return numberType
}

func fieldFor(k string, v interface{}) spec.FieldSpec {
	switch v.(type) {
	case float64:
		return spec.FieldSpec{Name: k, Type: inferNumberType(v)}
	case string:
		return spec.FieldSpec{Name: k, Type: stringType}
	case bool:
		return spec.FieldSpec{Name: k, Type: booleanType}
	default:
		return spec.FieldSpec{Name: k, Type: objectType, Fields: fieldsFor(v.(map[string]interface{}))}
	}
}
