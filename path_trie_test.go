package tagconfig

import (
	"reflect"
	"testing"
)

func Test_PathTrie_Put(t *testing.T) {
	node := newPathTrie("")
	node.Put("a", "{}")
	gotNode, ok := node.StartWith("a")
	if !ok {
		t.Log("not get a")
	}

	if gotNode.Value != "{}" {
		t.Errorf("gotValue = %v, want %v", gotNode.Value, "{}")
	}

	node.Put("a", "c")
	gotNode, ok = node.StartWith("a")
	if gotNode.Value != "c" {
		t.Errorf("gotValue = %v, want %v", gotNode.Value, "c")
	}

	node.Put("a.b", "d")
	gotNode, ok = node.StartWith("a")
	if gotNode.Value != "c" {
		t.Errorf("gotValue = %v, want %v", gotNode.Value, "c")
	}

	gotNode, ok = node.StartWith("a.b")
	if gotNode.Value != "d" {
		t.Errorf("gotValue = %v, want %v", gotNode.Value, "d")
	}

	gotNode, ok = node.StartWith("c")
	if ok {
		t.Errorf("gotOk = %v, want %v", ok, false)
	}

	if gotNode != nil {
		t.Errorf("gotValue = %v, want %v", gotNode, nil)
	}
}

func Test_PathTrie_StartWith(t *testing.T) {
	type fields struct {
		hasNext bool
		Value   string
		nexts   map[string]*PathTrie
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
			dt := newPathTrie(nil)
			dt.Put(tt.args.fill, tt.args.value)
			gotNode, gotOk := dt.StartWith(tt.args.k)
			if gotOk != tt.wantOk {
				t.Errorf("PathTrie.StartWith() args:%+v gotOk = %v, want %v", tt, gotOk, tt.wantOk)
			}

			if gotNode == nil != tt.wantNodeNil {
				t.Errorf("PathTrie.StartWith() args:%+v gotNil = %v, want %v", tt, gotNode == nil, tt.wantNodeNil)
			}

			if gotNode == nil || tt.wantNodeNil {
				return
			}

			if !reflect.DeepEqual(gotNode.Value, tt.wantNodeVal) {
				t.Errorf("PathTrie.StartWith() args:%+v gotNode = %v, want %v", tt, gotNode.Value, tt.wantNodeVal)
			}
		})
	}
}
