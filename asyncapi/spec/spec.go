package spec

type MessageSpec struct {
	Topic    string
	Exchange string
	Payload  PayloadSpec
}

type PayloadSpec struct {
	Type   string       `json:"type"`
	Fields []*FieldSpec `json:"fields"`
}

type FieldSpec struct {
	Name     string       `json:"name"`
	Type     string       `json:"type"`
	Format   string       `json:"format"`
	Fields   []*FieldSpec `json:"fields"`
	Item     *FieldSpec   `json:"Item"`
	Optional bool         `json:"optional"`
}

type ServerSpec struct {
	Name    string
	Version string
}
