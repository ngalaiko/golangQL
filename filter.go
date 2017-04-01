package golangQL

import (
	"reflect"
	"strings"
)

func Filter(v interface{}, fields string) (interface{}, error) {
	if len(fields) == 0 {
		return v, nil
	}

	baseLevel, err := parse(fields)
	if err != nil {
		return nil, err
	}

	return filter(v, baseLevel)
}

func filter(v interface{}, baseLevel *level) (interface{}, error) {
	typ := reflect.TypeOf(v)
	switch typ.Kind() {
	case reflect.Ptr:
		return filterPtr(v, baseLevel)
	case reflect.Struct:
		return filterStruct(v, baseLevel)
	case reflect.Slice:
		return filterSlice(v, baseLevel)
	default:
		return v, nil
	}

	return v, nil
}

func filterPtr(v interface{}, baseLevel *level) (interface{}, error) {
	rv := reflect.ValueOf(v).Elem()
	return filter(rv.Interface(), baseLevel)
}

func filterSlice(v interface{}, baseLevel *level) (interface{}, error) {
	s := reflect.ValueOf(v)

	out := make([]map[string]interface{}, s.Len())
	outV := reflect.ValueOf(out)
	for i := 0; i < s.Len(); i++ {
		v := s.Index(i).Interface()
		filteredStruc, err := filter(v, baseLevel)
		if err != nil {
			return nil, err
		}
		filteredValue := reflect.ValueOf(filteredStruc)

		outV.Index(i).Set(filteredValue)
	}

	return out, nil
}

func filterStruct(v interface{}, baseLevel *level) (interface{}, error) {
	rt := reflect.TypeOf(v)
	out := make(map[string]interface{}, rt.NumField())

	if err := filterStructFields(v, &out, baseLevel); err != nil {
		return nil, err
	}

	return out, nil
}

func filterStructFields(v interface{}, out *map[string]interface{}, level *level) error {
	rt, rv := reflect.TypeOf(v), reflect.ValueOf(v)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag, ok := field.Tag.Lookup("json")
		tagName := fieldParseTag(tag)

		if !ok {
			embeded, err := filter(rv.Field(i).Interface(), level)
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

		subfield := level.findSublevel(tagName)
		if subfield != nil {
			subStruct, err := filter(rv.Field(i).Interface(), subfield)
			if err != nil {
				return err
			}

			(*out)[tagName] = subStruct
			continue
		}

		if level.containsField(tagName) {
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
