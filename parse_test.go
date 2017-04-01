package golangQL

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str := `{ a f1 { b f11 { c c } b f12 {e} } f2 {v f21 {q}  } a }`
		parse(str)
	}
}

func TestParse__should_parse(t *testing.T) {
	str := `{ a f1 { b f11 { c c } b f12 {e} } f2 {v f21 {q}  } a }`
	result, err := parse(str)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, result.name, baseFieldName)
	assert.Equal(t, 2, len(result.fields))
	assert.Equal(t, 2, len(result.sublevels))
	assert.Nil(t, result.parent)

	assert.Equal(t, result.sublevels[0].name, "f1")
	assert.Equal(t, 2, len(result.sublevels[0].fields))
	assert.Equal(t, 2, len(result.sublevels[0].sublevels))
	assert.Equal(t, result.sublevels[0].parent, result)

	assert.Equal(t, result.sublevels[1].name, "f2")
	assert.Equal(t, 1, len(result.sublevels[1].fields))
	assert.Equal(t, 1, len(result.sublevels[1].sublevels))
	assert.Equal(t, result.sublevels[1].parent, result)

	assert.Equal(t, result.sublevels[0].sublevels[0].name, "f11")
	assert.Equal(t, 2, len(result.sublevels[0].sublevels[0].fields))
	assert.Equal(t, 0, len(result.sublevels[0].sublevels[0].sublevels))
	assert.Equal(t, result.sublevels[0].sublevels[0].parent, result.sublevels[0])

	assert.Equal(t, result.sublevels[0].sublevels[1].name, "f12")
	assert.Equal(t, 1, len(result.sublevels[0].sublevels[1].fields))
	assert.Equal(t, 0, len(result.sublevels[0].sublevels[1].sublevels))
	assert.Equal(t, result.sublevels[0].sublevels[1].parent, result.sublevels[0])

	assert.Equal(t, result.sublevels[1].sublevels[0].name, "f21")
	assert.Equal(t, 1, len(result.sublevels[1].sublevels[0].fields))
	assert.Equal(t, 0, len(result.sublevels[1].sublevels[0].sublevels))
	assert.Equal(t, result.sublevels[1].sublevels[0].parent, result.sublevels[1])
}

func TestParse__should_not_validate_format(t *testing.T) {
	str := "a { b } }"
	_, err := parse(str)
	assert.Equal(t, ErrWrongGQLFormat, err)
}

func TestParse__should_not_validate_format_2(t *testing.T) {
	str := "{a }{ b }"
	_, err := parse(str)
	assert.Equal(t, ErrWrongGQLFormat, err)
}

func TestParse__should_not_validate_commas(t *testing.T) {
	str := "{ e } a { e } }"
	_, err := parse(str)
	assert.Equal(t, ErrWrongGQLFormat, err)
}

func TestParse__should_not_validate_null(t *testing.T) {
	str := ""
	_, err := parse(str)
	assert.Equal(t, ErrWrongGQLFormat, err)
}
