package golangQL

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type testBaseStruct struct {
	C string `json:"c"`
}

func newTestBaseStruct() *testBaseStruct {
	return &testBaseStruct{
		C: "c field",
	}
}

type testStruct struct {
	*testBaseStruct
	A int         `json:"a,omitempty"`
	B string      `json:"b"`
	E *testStruct `json:"e"`
}

func newTestStruct() *testStruct {
	return &testStruct{
		A:              12,
		B:              "b field",
		testBaseStruct: newTestBaseStruct(),
	}
}

var (
	valueA = reflect.ValueOf("a")
	valueB = reflect.ValueOf("b")
	valueC = reflect.ValueOf("c")
	valueE = reflect.ValueOf("e")
)

func TestFilterJsonFields__should_filter_struct_ptr(t *testing.T) {
	s := newTestStruct()
	s.E = newTestStruct()

	query := "{a c e { c } }"
	res, err := Filter(s, query)
	if err != nil {
		t.Fatal(err)
	}
	v := reflect.ValueOf(res)

	assert.True(t, v.MapIndex(valueA).IsValid())
	assert.True(t, v.MapIndex(valueC).IsValid())
	assert.True(t, v.MapIndex(valueE).IsValid())
	assert.True(t, v.MapIndex(valueE).IsValid())
	assert.True(t, v.MapIndex(valueE).Elem().MapIndex(valueC).IsValid())
	assert.False(t, v.MapIndex(valueB).IsValid())
}

func TestFilterJsonFields__should_not_filter_empty_query(t *testing.T) {
	s := newTestStruct()
	s.E = newTestStruct()

	query := ""
	res, err := Filter(s, query)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, reflect.ValueOf(s), reflect.ValueOf(res))
}

func TestFilterJsonFields__should_filter_slice(t *testing.T) {
	slice := []*testStruct{}

	for i := 0; i < 10; i++ {
		s := newTestStruct()
		s.E = newTestStruct()
		slice = append(slice, s)
	}

	query := "{ a e { a} }"
	res, err := Filter(slice, query)
	if err != nil {
		t.Fatal(err)
	}
	v := reflect.ValueOf(res)

	if assert.Equal(t, len(slice), v.Len()) {
		for i := 0; i < v.Len(); i++ {
			vi := v.Index(i)

			assert.True(t, vi.MapIndex(valueA).IsValid())
			assert.False(t, vi.MapIndex(valueB).IsValid())
			assert.False(t, vi.MapIndex(valueC).IsValid())
			assert.True(t, vi.MapIndex(valueE).IsValid())
			assert.True(t, vi.MapIndex(valueE).Elem().MapIndex(valueA).IsValid())
		}
	}
}
