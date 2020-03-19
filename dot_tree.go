package tagconfig

import (
	"sort"
	"strings"
)

type dotTree struct {
	hasNext bool
	value   interface{}
	nexts   map[string]*dotTree
}

func newDotTree(value interface{}) *dotTree {
	return &dotTree{
		value: value,
	}
}

func (t *dotTree) Put(path string, value interface{}) {
	var (
		parts   = strings.Split(path, ".")
		current = t
		length  = len(parts)
	)

	for i := 0; i < length; i++ {
		part := parts[i]
		v, ok := current.nexts[part]
		if ok {
			current = v
			if length-1 == i {
				v.value = value
				return
			}
			continue
		}

		if !current.hasNext {
			current.nexts = make(map[string]*dotTree)
			current.hasNext = true
		}

		node := newDotTree(nil)
		if i == length-1 {
			node = newDotTree(value)
		}
		current.nexts[part] = node
		current = current.nexts[part]
	}
}

func (t *dotTree) Get(path string) (node *dotTree, ok bool) {
	parts := strings.Split(path, ".")
	length := len(parts)
	current := t
	i := 0
	for i < length {
		part := parts[i]
		next, ok := current.nexts[part]
		if ok {
			current = next
			i++
			continue
		}
		break
	}
	if i != length {
		return nil, false
	}
	return current, true
}

func (t *dotTree) FlattenChild() (ret []string) {
	node := t
	if node == nil {
		return
	}

	if node.hasNext {
		ret = falt(node, "")
		for i := 0; i < len(ret); i++ {
			ret[i] = strings.TrimLeft(ret[i], ".")
		}
	}

	sort.Strings(ret)
	return ret
}

func falt(root *dotTree, last string) (ret []string) {
	if root == nil {
		return
	}

	if !root.hasNext {
		return []string{last}
	}

	for k := range root.nexts {
		v := root.nexts[k]
		list := falt(v, last+"."+k)
		ret = append(ret, list...)
	}
	return
}
