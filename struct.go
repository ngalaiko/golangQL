package golangQL

import (
	"bytes"
	"reflect"
)

// cache key

type cacheKey struct {
	typ    reflect.Type
	fields string
}

func newCacheKey(typ reflect.Type, tree *node) cacheKey {
	return cacheKey{
		typ:    typ,
		fields: concatStrings(" ", tree.fields...),
	}
}

func concatStrings(delimer string, strings ...string) string {
	var result bytes.Buffer
	for _, s := range strings {
		result.WriteString(delimer)
		result.WriteString(s)
	}

	return result.String()
}

// node

type node struct {
	name     string
	fields   []string
	children []*node
	parent   *node
}

func (t *node) appendChild(children *node) {
	t.children = append(t.children, children)
}

func (t *node) findChildByName(childName string) *node {
	for _, child := range t.children {
		if child.name == childName {
			return child
		}
	}

	return nil
}

func (t *node) containsField(fieldName string) bool {
	for _, field := range t.fields {
		if fieldName == field {
			return true
		}
	}

	return false
}

func newNode(name string) *node {
	return &node{
		name:     name,
		fields:   []string{},
		children: []*node{},
		parent:   nil,
	}
}
