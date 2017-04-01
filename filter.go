package golangQL

import (
	"reflect"
	"strings"
)

func filter(v interface{}, tree *node) (interface{}, error) {
	typ := reflect.TypeOf(v)
	switch typ.Kind() {
	case reflect.Ptr:
		return filterPtr(v, tree)
	case reflect.Struct:
		return filterStruct(v, tree)
	case reflect.Slice:
		return filterSlice(v, tree)
	default:
		return v, nil
	}

	return v, nil
}

func filterPtr(v interface{}, tree *node) (interface{}, error) {
	rv := reflect.ValueOf(v).Elem()
	return filter(rv.Interface(), tree)
}

func filterSlice(v interface{}, tree *node) (interface{}, error) {
	s := reflect.ValueOf(v)

	out := make([]map[string]interface{}, s.Len())
	outV := reflect.ValueOf(out)
	for i := 0; i < s.Len(); i++ {
		v := s.Index(i).Interface()
		filteredStruc, err := filter(v, tree)
		if err != nil {
			return nil, err
		}
		filteredValue := reflect.ValueOf(filteredStruc)

		outV.Index(i).Set(filteredValue)
	}

	return out, nil
}

func filterStruct(v interface{}, tree *node) (interface{}, error) {
	rt := reflect.TypeOf(v)
	out := make(map[string]interface{}, rt.NumField())

	if err := filterStructFields(v, &out, tree); err != nil {
		return nil, err
	}

	return out, nil
}

func filterStructFields(v interface{}, out *map[string]interface{}, tree *node) error {
	rt, rv := reflect.TypeOf(v), reflect.ValueOf(v)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag, ok := field.Tag.Lookup("json")
		tagName := fieldParseTag(tag)

		if !ok {
			embeded, err := filter(rv.Field(i).Interface(), tree)
			if err != nil {
				return err
			}

			if reflect.TypeOf(embeded).Kind() != reflect.Map {
				(*out)[tagName] = rv
				continue
			}

			embededMapValue := reflect.ValueOf(embeded)
			for _, key := range embededMapValue.MapKeys() {
				(*out)[key.String()] = embededMapValue.MapIndex(key).Interface()
			}
			continue
		}

		subfield := tree.findChild(tagName)
		if subfield != nil {
			subStruct, err := filter(rv.Field(i).Interface(), subfield)
			if err != nil {
				return err
			}

			(*out)[tagName] = subStruct
			continue
		}

		if tree.containsField(tagName) {
			(*out)[tagName] = rv.Field(i).Interface()
		}
	}

	return nil
}

func fieldParseTag(tag string) string {
	if i := strings.Index(tag, ","); i != -1 {
		name := tag[:i]
		tag = tag[i+1:]
		return name
	}
	return tag
}
