package spec

type MessageSpec struct {
	Topic   string
	Payload PayloadSpec
}

type PayloadSpec struct {
	Type   string      `json:"type"`
	Fields []FieldSpec `json:"fields"`
}

type FieldSpec struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ServerSpec struct {
	Name    string
	Version string
}
