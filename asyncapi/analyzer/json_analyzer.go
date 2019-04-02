package analyzer

import (
	"encoding/json"

	"github.com/uesteibar/scribano/asyncapi/spec"
)

// JSONAnalyzer analyzes json payloads to build a spec
type JSONAnalyzer struct{}

const (
	arrayType   = "array"
	booleanType = "boolean"
	integerType = "integer"
	numberType  = "number"
	objectType  = "object"
	stringType  = "string"
	unknownType = "binary"
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

func typeFor(v interface{}) string {
	switch v.(type) {
	case float64:
		return inferNumberType(v)
	case string, nil:
		return stringType
	case bool:
		return booleanType
	case []interface{}:
		return arrayType
	case map[string]interface{}:
		return objectType
	default:
		return stringType
	}
}

func arrayItemFor(l []interface{}) *spec.FieldSpec {
	var item interface{}
	if len(l) > 0 {
		item = l[0]
	}
	t := typeFor(item)
	a := &spec.FieldSpec{Type: t}

	if t == "object" {
		a.Fields = fieldsFor(item.(map[string]interface{}))
	} else if t == "array" {
		a.Item = arrayItemFor(item.([]interface{}))
	}

	return a
}

func fieldFor(k string, v interface{}) spec.FieldSpec {
	switch v.(type) {
	case []interface{}:
		return spec.FieldSpec{Name: k, Type: arrayType, Item: arrayItemFor(v.([]interface{}))}
	case map[string]interface{}:
		return spec.FieldSpec{Name: k, Type: objectType, Fields: fieldsFor(v.(map[string]interface{}))}
	default:
		return spec.FieldSpec{Name: k, Type: typeFor(v)}
	}
}
