package spec

type MessageSpec struct {
	Topic   string
	Payload PayloadSpec
}

type PayloadSpec struct {
	Type   string
	Fields []FieldSpec
}

type FieldSpec struct {
	Name string
	Type string
}
