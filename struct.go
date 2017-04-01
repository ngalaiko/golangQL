package golangQL

type node struct {
	name     string
	fields   []string
	children []*node
	parent   *node
}

func (t *node) appendChild(child *node) {
	t.children = append(t.children, child)
}

func (t *node) findChild(childName string) *node {
	for _, child := range t.children {
		if child.name == childName {
			return child
		}
	}

	return nil
}

func (t *node) containsField(childName string) bool {
	for _, field := range t.fields {
		if childName == field {
			return true
		}
	}

	return false
}

func newNode(name string) *node {
	return &node{
		name:      name,
		fields:    []string{},
		children: []*node{},
		parent:    nil,
	}
}
