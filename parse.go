package golangQL

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

const (
	openComma     = '{'
	closeComma    = '}'
	baseFieldName = "data"
	space         = ' '
)

var (
	ErrWrongGQLFormat = errors.New("errors.GQLWrongFormat")
)

func parse(str string) (*level, error) {
	str, err := formatString(str)
	if err != nil {
		return nil, err
	}

	baseLevel := newLevel(baseFieldName)
	currentLevel := baseLevel

	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case openComma: // level down
			newLevel := newLevel(buf.String())
			newLevel.parent = currentLevel
			currentLevel.appendSublevel(newLevel)

			currentLevel = newLevel

			buf.Reset()
		case closeComma: // level up
			if currentLevel.parent == nil {
				return nil, ErrWrongGQLFormat
			}

			currentLevel = currentLevel.parent
		case space: // append buffer to current level
			currentLevel.fields = append(currentLevel.fields, buf.String())
			buf.Reset()
		default: // append buffer
			buf.WriteRune(char)
		}
	}

	return baseLevel, nil
}

func formatString(str string) (string, error) {
	if len(str) == 0 {
		return "", ErrWrongGQLFormat
	}

	str = strings.Replace(str, "}", " }", -1) // add spaces

	str, err := deleteDoubleSpaces(str)
	if err != nil {
		return "", err
	}

	if str[0] != openComma || str[len(str)-1] != closeComma {
		return "", ErrWrongGQLFormat
	}

	str = strings.Replace(str, "{ ", "{", -1) // get rid of '{ '
	str = strings.Replace(str, " {", "{", -1) // get rid of ' {'
	str = strings.Replace(str, "} ", "}", -1) // get rid of '} '
	str = strings.Trim(str, "{}")

	return str, nil
}

func deleteDoubleSpaces(str string) (string, error) {
	str = strings.TrimSpace(str)
	r, err := regexp.Compile("( )+")
	if err != nil {
		return "", err
	}
	str = r.ReplaceAllString(str, " ")

	return str, nil
}
