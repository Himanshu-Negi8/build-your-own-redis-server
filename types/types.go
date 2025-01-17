package types

type RESPType string

const (
	RESPTypeSimpleString RESPType = "RESPTypeSimpleString"
	RESPTypeError        RESPType = "RESPTypeError"
	RESPTypeInteger      RESPType = "RESPTypeInteger"
	RESPTypeBulkString   RESPType = "RESPTypeBulkString"
	RESPTypeArray        RESPType = "RESPTypeArray"
)

type CustomValue struct {
	Value           string
	ValueExpiration int64
}
