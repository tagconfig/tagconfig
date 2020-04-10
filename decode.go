package tagconfig

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
)

type (
	// Unmarshaler  custom unmarshal function
	Unmarshaler interface {
		UnmarshalTagConfig(map[string]string) error
	}

	// ConfigProvider config provider
	ConfigProvider interface {
		Paires() ([]Paire, error)
		FieldInfo(field reflect.StructField) (namespace string, key string)
	}

	//Paire config paire
	Paire struct {
		Namespace string
		Key       string
		Value     string
	}

	// Decoder unmarshal config like json.Unmarshal
	Decoder struct {
		provider ConfigProvider
	}
)

// NewDecoder create a goconfig Decoder
func NewDecoder(provider ConfigProvider) *Decoder {
	d := &Decoder{
		provider: provider,
	}
	return d
}

// Decode works like Unmarshal, except it reads the decoder
func (d *Decoder) Decode(v interface{}) (err error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return errors.New("non-pointer passed to Unmarshal")
	}

	paires, err := d.provider.Paires()
	if err != nil {
		return
	}

	pairesTree := buildPairesTree(paires)
	if val.Elem().Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return d.unmarshalRoot(val.Elem(), pairesTree)
}

func (d *Decoder) unmarshalRoot(val reflect.Value, pairesTree map[string]*PathTrie) (err error) {
	for i := 0; i < val.NumField(); i++ {
		var (
			paires    = make(map[string]string)
			namespace string
			key       string
			useJSON   bool
		)

		namespace, key = d.provider.FieldInfo(val.Type().Field(i))

		if strings.Contains(key, ",") {
			useJSON = strings.Split(key, ",")[1] == "json"
			key = strings.Split(key, ",")[0]
		}

		node, ok := pairesTree[namespace].StartWith(key)
		if ok && node.Value != nil {
			paires["."] = node.Value.(string)
		}

		if ok {
			keys := node.FlattenChild()
			for _, k := range keys {
				kNode, _ := node.StartWith(k)
				paires[k] = kNode.Value.(string)
			}
		}

		if len(paires) == 0 {
			continue
		}

		err = d.unmarshal(val.Field(i), paires, useJSON)
		if err != nil {
			return
		}
	}

	return
}

func (d *Decoder) unmarshal(val reflect.Value, paires map[string]string, useJSON bool) (err error) {
	if val.Kind() == reflect.Interface && !val.IsNil() {
		e := val.Elem()
		if e.Kind() == reflect.Ptr && !e.IsNil() {
			val = e
		}
	}

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}

	if !val.IsValid() {
		return
	}

	if len(paires) == 0 {
		return
	}

	if val.CanInterface() && val.Type().Implements(unmarshalerType) {
		return d.unmarshalInterface(val.Interface().(Unmarshaler), paires)
	}

	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(unmarshalerType) {
			return d.unmarshalInterface(pv.Interface().(Unmarshaler), paires)
		}
	}

	if useJSON {
		err = copyJSONValue(val, []byte(paires["."]))
	}

	switch v := val; v.Kind() {
	default:
		return errors.New("unknown type " + v.Type().String())
	case reflect.Slice, reflect.Map:
		value := paires["."]
		err = copyJSONValue(val, []byte(value))
	case reflect.Interface, reflect.Bool, reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.String:
		value := paires["."]
		err = copyValue(val, []byte(value))
	case reflect.Struct:
		return d.unmarshalStruct(val, paires)
	}
	return
}

func (d *Decoder) unmarshalStruct(val reflect.Value, paires map[string]string) (err error) {
	for i := 0; i < val.NumField(); i++ {
		newPaires := make(map[string]string)
		_, key := d.provider.FieldInfo(val.Type().Field(i))
		var useJSON bool
		if strings.Contains(key, ",") {
			useJSON = strings.Split(key, ",")[1] == "json"
			key = strings.Split(key, ",")[0]
		}
		v, ok := paires[key]
		if ok {
			newPaires["."] = v
			delete(paires, key)
		}

		for k, v := range paires {
			if !strings.HasPrefix(k, key) {
				continue
			}
			newK := strings.Replace(k, key+".", "", 1)
			newPaires[newK] = v
			delete(paires, k)
		}
		err = d.unmarshal(val.Field(i), newPaires, useJSON)
		if err != nil {
			return
		}
	}
	return nil
}

func (d *Decoder) unmarshalInterface(val Unmarshaler, paires map[string]string) error {
	return val.UnmarshalTagConfig(paires)
}

func copyJSONValue(dst reflect.Value, src []byte) (err error) {

	if dst.Kind() == reflect.Ptr {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
	} else {
		dst = dst.Addr()
	}
	jsonDecoder := json.NewDecoder(bytes.NewBuffer(src))
	return jsonDecoder.Decode(dst.Interface())
}

func copyValue(dst reflect.Value, src []byte) (err error) {
	dst0 := dst
	if dst.Kind() == reflect.Ptr {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		dst = dst.Elem()
	}

	switch dst.Kind() {
	case reflect.Invalid:
	default:
		return errors.New("cannot unmarshal into " + dst0.Type().String())
	case reflect.Interface:
		if dst.NumMethod() == 0 {
			dst.Set(reflect.ValueOf(string(src)))
		} else {
			return &UnmarshalTypeError{Value: "string", Type: dst.Type()}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(src) == 0 {
			dst.SetInt(0)
			return nil
		}
		itmp, err := strconv.ParseInt(strings.TrimSpace(string(src)), 10, dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetInt(itmp)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if len(src) == 0 {
			dst.SetUint(0)
			return nil
		}
		utmp, err := strconv.ParseUint(strings.TrimSpace(string(src)), 10, dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetUint(utmp)
	case reflect.Float32, reflect.Float64:
		if len(src) == 0 {
			dst.SetFloat(0)
			return nil
		}
		ftmp, err := strconv.ParseFloat(strings.TrimSpace(string(src)), dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetFloat(ftmp)
	case reflect.Bool:
		if len(src) == 0 {
			dst.SetBool(false)
			return nil
		}
		value, err := strconv.ParseBool(strings.TrimSpace(string(src)))
		if err != nil {
			return err
		}
		dst.SetBool(value)
	case reflect.String:
		dst.SetString(string(src))
	case reflect.Slice:
		if len(src) == 0 {
			src = []byte{}
		}
		dst.SetBytes(src)
	}
	return nil
}

func buildPairesTree(paires []Paire) (root map[string]*PathTrie) {
	root = make(map[string]*PathTrie)
	for _, p := range paires {
		_, ok := root[p.Namespace]
		if !ok {
			root[p.Namespace] = newPathTrie(nil)
		}
		root[p.Namespace].Put(p.Key, p.Value)
	}
	return
}
