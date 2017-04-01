package golangQL

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type A struct {
	A string
	C string
}

type B struct {
	A
	B string
}

type C struct {
	*B
	C string
}

func TestFieldsForType(t *testing.T) {
	typ := reflect.TypeOf(new(C)).Elem()
	fields := fieldsForType(typ)

	assert.Equal(t, 3, len(fields))

	assert.Equal(t, []int{0, 0, 0}, fields[0].index)
	assert.Equal(t, []int{0, 1}, fields[1].index)
	assert.Equal(t, []int{1}, fields[2].index)

	assert.Equal(t, "A", fields[0].name)
	assert.Equal(t, "B", fields[1].name)
	assert.Equal(t, "C", fields[2].name)
}
