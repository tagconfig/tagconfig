package tagconfig

import (
	"sort"
	"strings"
)

type PathTrie struct {
	Value    interface{}
	Has      bool
	Children map[string]*PathTrie
}

func newPathTrie(value interface{}) *PathTrie {
	return &PathTrie{
		Value: value,
	}
}

//Put inserts a path into the trie.
func (t *PathTrie) Put(path string, value interface{}) {
	parts := strings.Split(path, ".")
	var current = t
	for i := 0; i < len(parts); i++ {
		if current.Children == nil {
			current.Children = make(map[string]*PathTrie)
		}
		var ok bool
		_, ok = current.Children[parts[i]]
		if !ok {
			current.Children[parts[i]] = new(PathTrie)
			current = current.Children[parts[i]]
			continue
		}
		current = current.Children[parts[i]]
	}
	current.Has = true
	current.Value = value
	return
}

//Search returns if the word is in the trie
func (t *PathTrie) Search(path string) (node *PathTrie, ok bool) {
	var current = t
	parts := strings.Split(path, ".")
	for i := 0; i < len(parts); i++ {
		var ok bool
		current, ok = current.Children[parts[i]]
		if !ok {
			return nil, false
		}
	}
	if current.Has {
		return current, true
	}
	return nil, false
}

//StartWith returns if there is any word in the trie that starts with the given prefix.
func (t *PathTrie) StartWith(prefix string) (node *PathTrie, ok bool) {
	var current = t
	var parts = strings.Split(prefix, ".")
	for i := 0; i < len(parts); i++ {
		var ok bool
		current, ok = current.Children[parts[i]]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

func (t *PathTrie) FlattenChild() (ret []string) {
	node := t
	if node == nil {
		return
	}

	if node.Children != nil {
		ret = falt(node, "")
		for i := 0; i < len(ret); i++ {
			ret[i] = strings.TrimLeft(ret[i], ".")
		}
	}

	sort.Strings(ret)
	return ret
}

func falt(root *PathTrie, last string) (ret []string) {
	if root == nil {
		return
	}

	if root.Children == nil {
		return []string{last}
	}

	for k := range root.Children {
		v := root.Children[k]
		list := falt(v, last+"."+k)
		ret = append(ret, list...)
	}
	return
}
