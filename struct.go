package golangQL

// level

type level struct {
	name      string
	fields    []string
	sublevels []*level
	parent    *level
}

func (t *level) appendSublevel(subLevel *level) {
	t.sublevels = append(t.sublevels, subLevel)
}

func (t *level) findSublevel(fieldName string) *level {
	for _, sublevel := range t.sublevels {
		if sublevel.name == fieldName {
			return sublevel
		}
	}

	return nil
}

func (t *level) containsField(fieldName string) bool {
	for _, field := range t.fields {
		if fieldName == field {
			return true
		}
	}

	return false
}

func newLevel(name string) *level {
	return &level{
		name:      name,
		fields:    []string{},
		sublevels: []*level{},
		parent:    nil,
	}
}
