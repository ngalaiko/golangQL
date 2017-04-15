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
	spacesRegEx = regexp.MustCompile("( )+")
)

func parse(query string) (*node, error) {
	query, err := formatQuery(query)
	if err != nil {
		return nil, err
	}

	tree := newNode(baseFieldName)
	currentNode := tree

	var buf bytes.Buffer
	for _, char := range query {
		switch char {
		case openComma:
			newNode := newNode(buf.String())
			newNode.parent = currentNode
			currentNode.appendChild(newNode)

			currentNode = newNode

			buf.Reset()
		case closeComma:
			if currentNode.parent == nil {
				return nil, ErrWrongGQLFormat
			}

			currentNode = currentNode.parent
		case space:
			currentNode.fields = append(currentNode.fields, buf.String())
			buf.Reset()
		default:
			buf.WriteRune(char)
		}
	}

	return tree, nil
}

func formatQuery(query string) (string, error) {
	if len(query) == 0 {
		return "", ErrWrongGQLFormat
	}

	query = strings.Replace(query, "}", " }", -1)

	query, err := deleteDoubleSpaces(query)
	if err != nil {
		return "", err
	}

	if query[0] != openComma || query[len(query)-1] != closeComma {
		return "", ErrWrongGQLFormat
	}

	query = strings.Replace(query, "{ ", "{", -1)
	query = strings.Replace(query, " {", "{", -1)
	query = strings.Replace(query, "} ", "}", -1)
	query = strings.Trim(query, "{}")

	return query, nil
}

func deleteDoubleSpaces(str string) (string, error) {
	str = strings.TrimSpace(str)
	str = spacesRegEx.ReplaceAllString(str, " ")

	return str, nil
}
