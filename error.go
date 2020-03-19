package tagconfig

import (
	"reflect"
)

type UnmarshalTypeError struct {
	Value  string
	Type   reflect.Type // type of Go value it could not be assigned to
	Struct string       // name of the struct type containing the field
	Field  string       // the full path from root node to the field
}

func (e *UnmarshalTypeError) Error() string {
	if e.Struct != "" || e.Field != "" {
		return "tagconfig: cannot unmarshal " + e.Value + " into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String()
	}
	return "tagconfig: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}
