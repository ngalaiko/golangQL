package golangQL

import (
	"reflect"
	"strings"
)

/*
   part of TheQuestionru/jsonx package
*/

type field struct {
	name  string
	tag   string
	typ   reflect.Type
	index []int

	filter filterFunc
	child  *node
}

// Returns an invalid value for a nil embedded field, does not panic.
func fieldByIndex(v reflect.Value, index []int) reflect.Value {
	for _, i := range index {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return reflect.Value{}
			}
			v = v.Elem()
		}
		v = v.Field(i)
	}
	return v
}

func fieldsForType(typ reflect.Type) []field {
	fields := fieldsFlatten(typ, []int{})
	return fieldsUnique(fields)
}

// Flattens type fields and embedded type fields.
func fieldsFlatten(typ reflect.Type, parentIndex []int) []field {
	fields := []field{}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		// Skip unexported fields.
		if f.PkgPath != "" {
			continue
		}

		// Skip hidden fields.
		tag := f.Tag.Get("json")
		if tag == "-" {
			continue
		}

		// Calculate a field index.
		findex := make([]int, len(parentIndex))
		copy(findex, parentIndex)
		findex = append(findex, i)

		// Traverse embedded fields.
		if f.Anonymous {
			etype := f.Type
			if etype.Kind() == reflect.Ptr { // Follow the embedded struct pointer.
				etype = etype.Elem()
			}

			embedded := fieldsFlatten(etype, findex)
			fields = append(fields, embedded...)
			continue
		}

		// Add a normal field.
		name := f.Name
		tagName := fieldParseTag(tag)
		if tagName != "" {
			name = tagName
		}

		fields = append(fields, field{
			name:  name,
			typ:   f.Type,
			index: findex,
			tag:   tagName,
		})
	}

	return fields
}

// Fields with shorter indexes (top) hide fields with longer indexes with the same name.
func fieldsUnique(fields []field) []field {
	nameToIndex := map[string][]int{}
	for _, f := range fields {
		existing, ok := nameToIndex[f.name]
		if ok && len(existing) < len(f.index) {
			continue
		}

		nameToIndex[f.name] = f.index
	}

	result := []field{}
	for _, f := range fields {
		index := nameToIndex[f.name]
		if reflect.DeepEqual(index, f.index) {
			result = append(result, f)
		}
	}

	return result
}

func fieldParseTag(tag string) string {
	if i := strings.Index(tag, ","); i != -1 {
		name := tag[:i]
		tag = tag[i+1:]
		return name
	}
	return tag
}
