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

type TestConfigProvider struct {
	configMap map[string]string
}

func (t *TestConfigProvider) Paires() (paires []Paire, err error) {
	paires = make([]Paire, 0)
	for k, v := range t.configMap {
		paires = append(paires, Paire{
			Namespace: testDefaultNS,
			Key:       k,
			Value:     v,
		})
	}
	return
}

func (t *TestConfigProvider) FieldInfo(field reflect.StructField) (namespace string, key string) {
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
		Field1 string    `test:"field1"`
		Field2 float64   `test:"field2"`
		Field3 []float64 `test:"field3"`
	}

	tests := []Case{
		{
			fields: fields{provider: new(TestConfigProvider)},
			args: args{
				val:    reflect.ValueOf(new(Dst)).Elem(),
				paires: map[string]string{"field1": "1a", "Field2": "2b"},
			},
			wantErr: false,
		},

		{
			fields: fields{provider: new(TestConfigProvider)},
			args: args{
				val:    reflect.ValueOf(new(Dst)).Elem(),
				paires: map[string]string{"field1": "hello", "field2": "2b"},
			},
			wantErr: true,
		},
		{
			fields: fields{provider: new(TestConfigProvider)},
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
			if err := d.unmarshal(tt.args.val, tt.args.paires); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
