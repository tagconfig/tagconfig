package tagconfig

import (
	"reflect"
	"testing"
)

func Test_dotTree_Put(t *testing.T) {
	node := newDotTree("")
	node.Put("a", "{}")
	gotNode, ok := node.Get("a")
	if !ok {
		t.Log("not get a")
	}

	if gotNode.value != "{}" {
		t.Errorf("gotValue = %v, want %v", gotNode.value, "{}")
	}

	node.Put("a", "c")
	gotNode, ok = node.Get("a")
	if gotNode.value != "c" {
		t.Errorf("gotValue = %v, want %v", gotNode.value, "c")
	}

	node.Put("a.b", "d")
	gotNode, ok = node.Get("a")
	if gotNode.value != "c" {
		t.Errorf("gotValue = %v, want %v", gotNode.value, "c")
	}

	gotNode, ok = node.Get("a.b")
	if gotNode.value != "d" {
		t.Errorf("gotValue = %v, want %v", gotNode.value, "d")
	}

	gotNode, ok = node.Get("c")
	if ok {
		t.Errorf("gotOk = %v, want %v", ok, false)
	}

	if gotNode != nil {
		t.Errorf("gotValue = %v, want %v", gotNode, nil)
	}
}

func Test_dotTree_Get(t *testing.T) {
	type fields struct {
		hasNext bool
		Value   string
		nexts   map[string]*dotTree
	}

	type args struct {
		fill  string
		k     string
		value string
	}

	type usercase struct {
		name        string
		fields      fields
		args        args
		wantNodeVal interface{}
		wantNodeNil bool
		wantOk      bool
	}

	tests := []usercase{
		usercase{
			args:        args{fill: "abc.def", k: "abc", value: "a"},
			wantOk:      true,
			wantNodeVal: nil,
			wantNodeNil: false,
		},

		usercase{
			args:        args{fill: "abc.def", k: "abc.def", value: "bcd"},
			wantOk:      true,
			wantNodeVal: "bcd",
			wantNodeNil: false,
		},

		usercase{
			args:        args{fill: "abc.def", k: "abc.d"},
			wantOk:      false,
			wantNodeVal: "",
			wantNodeNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := newDotTree(nil)
			dt.Put(tt.args.fill, tt.args.value)
			gotNode, gotOk := dt.Get(tt.args.k)
			if gotOk != tt.wantOk {
				t.Errorf("dotTree.Get() gotOk = %v, want %v", gotOk, tt.wantOk)
			}

			if gotNode == nil != tt.wantNodeNil {
				t.Errorf("dotTree.Get() gotNil = %v, want %v", gotNode == nil, tt.wantNodeNil)
			}

			if gotNode == nil || tt.wantNodeNil {
				return
			}

			if !reflect.DeepEqual(gotNode.value, tt.wantNodeVal) {
				t.Errorf("dotTree.Get() gotNode = %v, want %v", gotNode.value, tt.wantNodeVal)
			}
		})
	}
}
