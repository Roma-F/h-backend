package storage

type OptValueType int

const (
	OVT_UINT    OptValueType = 0
	OVT_STRING  OptValueType = 1
	OVT_BOOLEAN OptValueType = 2
	OVT_BLOB    OptValueType = 3
)

var OptValueTypeNames = map[OptValueType]string{
	OVT_UINT:    "UINT",
	OVT_STRING:  "STRING",
	OVT_BOOLEAN: "BOOLEAN",
	OVT_BLOB:    "BLOB",
}

type DbOptType struct {
	OptType   int          `db:"opt_type"`
	Name      string       `db:"name"`
	ValueType OptValueType `db:"value_type"`
}
