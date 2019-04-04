package analyzer

import (
	"encoding/json"
	"regexp"

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

	defaultFormat  = ""
	dateFormat     = "date"
	dateTimeFormat = "date-time"

	dateMatcher     = "([12]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\\d|3[01]))"
	dateTimeMatcher = "^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$"
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

func fieldsFor(raw map[string]interface{}) []*spec.FieldSpec {
	var fields []*spec.FieldSpec
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

func inferStringFormat(s string) string {
	if matches(s, dateTimeMatcher) {
		return dateTimeFormat
	} else if matches(s, dateMatcher) {
		return dateFormat
	}
	return defaultFormat
}

func matches(s, rawRegexp string) bool {
	r := regexp.MustCompile(rawRegexp)
	return r.FindString(s) != ""
}

func typeAndFormatFor(v interface{}) (string, string) {
	switch v.(type) {
	case float64:
		return inferNumberType(v), defaultFormat
	case string:
		return stringType, inferStringFormat(v.(string))
	case bool:
		return booleanType, defaultFormat
	case []interface{}:
		return arrayType, defaultFormat
	case map[string]interface{}:
		return objectType, defaultFormat
	default:
		return stringType, defaultFormat
	}
}

func arrayItemFor(l []interface{}) *spec.FieldSpec {
	var item interface{}
	if len(l) > 0 {
		item = l[0]
	}
	t, f := typeAndFormatFor(item)
	a := &spec.FieldSpec{Type: t, Format: f}

	if t == "object" {
		a.Fields = fieldsFor(item.(map[string]interface{}))
	} else if t == "array" {
		a.Item = arrayItemFor(item.([]interface{}))
	}

	return a
}

func fieldFor(k string, v interface{}) *spec.FieldSpec {
	switch v.(type) {
	case []interface{}:
		return &spec.FieldSpec{Name: k, Type: arrayType, Item: arrayItemFor(v.([]interface{}))}
	case map[string]interface{}:
		return &spec.FieldSpec{Name: k, Type: objectType, Fields: fieldsFor(v.(map[string]interface{}))}
	default:
		t, f := typeAndFormatFor(v)
		return &spec.FieldSpec{Name: k, Type: t, Format: f}
	}
}
