package tagconfig

import (
	"reflect"
	"strings"
	"testing"
)

var (
	testDefaultNS = "test"
	testTag       = "test"
)

type testConfigProvider struct {
}

func (t *testConfigProvider) Paires() ([]Paire, error) {
	return nil, nil
}

func (t *testConfigProvider) FieldInfo(field reflect.StructField) (namespace string, key string) {
	namespace = testDefaultNS
	tag, hasTag := field.Tag.Lookup(testTag)
	fieldName := field.Name
	if !hasTag {
		return namespace, fieldName
	}

	parts := strings.Split(tag, ":")
	switch len(parts) {
	case 1:
		key = parts[0]
	case 2:
		namespace, key = parts[0], parts[1]
	default:
	}
	return
}

func TestDecoder_unmarshal(t *testing.T) {
	type fields struct {
		provider ConfigProvider
	}
	type args struct {
		val    reflect.Value
		paires map[string]string
	}
	type Case struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}

	type Dst struct {
		Field1 string `test:"field1"`
		Field2 int64  `test:"field2"`
		Field3 []int  `test:"field3"`
	}
	tests := []Case{
		Case{
			fields: fields{provider: new(testConfigProvider)},
			args: args{
				val:    reflect.ValueOf(new(Dst)).Elem(),
				paires: map[string]string{"field1": "1a", "Field2": "2b"},
			},
			wantErr: false,
		},

		Case{
			fields: fields{provider: new(testConfigProvider)},
			args: args{
				val:    reflect.ValueOf(new(Dst)).Elem(),
				paires: map[string]string{"field1": "hello", "field2": "2b"},
			},
			wantErr: true,
		},
		Case{
			fields: fields{provider: new(testConfigProvider)},
			args: args{
				val:    reflect.ValueOf(new(Dst)).Elem().FieldByName("field3"),
				paires: map[string]string{".": "[1,2]"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decoder{
				provider: tt.fields.provider,
			}
			if err := d.unmarshal(tt.args.val, tt.args.paires, false); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
