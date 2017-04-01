package golangQL

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := `{ a f1 { b f11 { c c } b f12 {e} } f2 {v f21 {q}  } a }`
		parse(query)
	}
}

func TestParse__should_parse(t *testing.T) {
	query := `{ a f1 { b f11 { c c } b f12 {e} } f2 {v f21 {q}  } a }`
	result, err := parse(query)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, result.name, baseFieldName)
	assert.Equal(t, 2, len(result.fields))
	assert.Equal(t, 2, len(result.children))
	assert.Nil(t, result.parent)

	assert.Equal(t, result.children[0].name, "f1")
	assert.Equal(t, 2, len(result.children[0].fields))
	assert.Equal(t, 2, len(result.children[0].children))
	assert.Equal(t, result.children[0].parent, result)

	assert.Equal(t, result.children[1].name, "f2")
	assert.Equal(t, 1, len(result.children[1].fields))
	assert.Equal(t, 1, len(result.children[1].children))
	assert.Equal(t, result.children[1].parent, result)

	assert.Equal(t, result.children[0].children[0].name, "f11")
	assert.Equal(t, 2, len(result.children[0].children[0].fields))
	assert.Equal(t, 0, len(result.children[0].children[0].children))
	assert.Equal(t, result.children[0].children[0].parent, result.children[0])

	assert.Equal(t, result.children[0].children[1].name, "f12")
	assert.Equal(t, 1, len(result.children[0].children[1].fields))
	assert.Equal(t, 0, len(result.children[0].children[1].children))
	assert.Equal(t, result.children[0].children[1].parent, result.children[0])

	assert.Equal(t, result.children[1].children[0].name, "f21")
	assert.Equal(t, 1, len(result.children[1].children[0].fields))
	assert.Equal(t, 0, len(result.children[1].children[0].children))
	assert.Equal(t, result.children[1].children[0].parent, result.children[1])
}

func TestParse__should_not_validate_format(t *testing.T) {
	query := "a { b } }"
	_, err := parse(query)
	assert.Equal(t, ErrWrongGQLFormat, err)
}

func TestParse__should_not_validate_format_2(t *testing.T) {
	query := "{a }{ b }"
	_, err := parse(query)
	assert.Equal(t, ErrWrongGQLFormat, err)
}

func TestParse__should_not_validate_commas(t *testing.T) {
	query := "{ e } a { e } }"
	_, err := parse(query)
	assert.Equal(t, ErrWrongGQLFormat, err)
}

func TestParse__should_not_validate_null(t *testing.T) {
	query := ""
	_, err := parse(query)
	assert.Equal(t, ErrWrongGQLFormat, err)
}
