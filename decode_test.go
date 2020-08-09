package tagconfig

import (
	"reflect"
	"strings"
	"testing"
)

var (
	testDefaultNS = "test"
	testTag       = "tagconfig"
)

type TestConfigProvider struct {
	configs [][]string //[[key,value,namespace]]
}

func (t *TestConfigProvider) Paires() (paires []Paire, err error) {
	paires = make([]Paire, 0)
	for _, v := range t.configs {
		p := Paire{
			Key:       v[0],
			Value:     v[1],
			Namespace: testDefaultNS,
		}

		if len(v) == 3 {
			p.Namespace = v[2]
		}
		paires = append(paires, p)
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
		Field1 string           `tagconfig:"field1"`
		Field2 float64          `tagconfig:"field2"`
		Field3 []float64        `tagconfig:"field3"`
		Field4 map[int64]string `tagconfig:"field4"`
		Field5 map[string]struct {
			Hello string
		} `tagconfig:"field5"`
	}

	tests := []Case{
		{
			fields: fields{provider: new(TestConfigProvider)},
			args: args{
				val:    reflect.ValueOf(new(Dst)).Elem(),
				paires: map[string]string{"field1": "1a", "Field2": "2b", "field4.1": "a", "field4.2": "b"},
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
				val:    reflect.ValueOf(new(Dst)).Elem(),
				paires: map[string]string{"field3": "[1,2]"},
			},
			wantErr: false,
		},

		{
			fields: fields{provider: new(TestConfigProvider)},
			args: args{
				val:    reflect.ValueOf(new(Dst)).Elem(),
				paires: map[string]string{"field4.a": "b", "field4.2": "d"},
			},
			wantErr: true,
		},

		{
			fields: fields{provider: new(TestConfigProvider)},
			args: args{
				val:    reflect.ValueOf(new(Dst)).Elem(),
				paires: map[string]string{"field5.a.Hello": "b", "field5.b.Hello": "x"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decoder{
				provider: tt.fields.provider,
			}
			if err := d.unmarshal(tt.args.val, tt.args.paires, nil); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}
