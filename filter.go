package golangQL

import (
	"reflect"
	"strings"
	"sync"
)

type filterFunc func(val reflect.Value, tree *node) (interface{}, error)

func (g *golangQL) filter(v interface{}, query string) (interface{}, error) {
	if len(query) == 0 {
		return v, nil
	}

	tree, err := parse(query)
	if err != nil {
		return nil, err
	}

	val := reflect.ValueOf(v)
	filter := g.getFilter(val.Type(), tree)

	result, err := filter(val, tree)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (g *golangQL) getFilter(typ reflect.Type, tree *node) filterFunc {
	g.RLock()
	key := newCacheKey(typ, tree)
	f := g.filters[key]
	g.RUnlock()
	if f != nil {
		return f
	}

	g.Lock()
	f = g.filters[key]
	if f != nil {
		g.Unlock()
		return f
	}

	recursive := sync.WaitGroup{}
	recursive.Add(1)
	g.filters[key] = func(val reflect.Value, tree *node) (interface{}, error) {
		recursive.Wait()
		return f(val, tree)
	}
	g.Unlock()

	f = g.newFilter(typ, tree)
	g.Lock()
	g.filters[key] = f
	g.Unlock()

	recursive.Done()
	return f
}

func (g *golangQL) newFilter(typ reflect.Type, tree *node) filterFunc {
	switch typ.Kind() {
	case reflect.Ptr:
		return g.newPtrFilter(typ, tree)
	case reflect.Struct:
		return g.newStructFilter(typ, tree)
	case reflect.Slice:
		return g.newSliceFilter(typ, tree)
	default:
		return defaultFilter()
	}
}

// Filter constructors

func defaultFilter() filterFunc {
	return func(val reflect.Value, tree *node) (interface{}, error) {
		return val.Interface(), nil
	}
}

func (g *golangQL) newPtrFilter(typ reflect.Type, tree *node) filterFunc {
	elemFilter := g.getFilter(typ.Elem(), tree)

	return func(val reflect.Value, tree *node) (interface{}, error) {
		v := val.Elem()

		return elemFilter(v, tree)
	}
}

func (g *golangQL) newSliceFilter(typ reflect.Type, tree *node) filterFunc {
	elemFilter := g.getFilter(typ.Elem(), tree)

	return func(val reflect.Value, tree *node) (interface{}, error) {
		resultMap := make([]map[string]interface{}, val.Len())
		resultMapValue := reflect.ValueOf(resultMap)
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)

			filteredStruct, err := elemFilter(v, tree)
			if err != nil {
				return nil, err
			}

			filteredValue := reflect.ValueOf(filteredStruct)

			resultMapValue.Index(i).Set(filteredValue)
		}

		return resultMap, nil
	}
}

func (g *golangQL) newStructFilter(typ reflect.Type, tree *node) filterFunc {
	fields := fieldsForType(typ)
	for i := range fields {
		fields[i].child = tree.findChildByName(fields[i].tag)
		if fields[i].child != nil {
			fields[i].filter = g.getFilter(fields[i].typ, fields[i].child)
			continue
		}

		fields[i].filter = g.getFilter(fields[i].typ, tree)
	}

	return func(val reflect.Value, tree *node) (interface{}, error) {
		typ := val.Type()
		resultMap := make(map[string]interface{}, typ.NumField())

		for _, field := range fields {
			fieldValue := fieldByIndex(val, field.index)
			child := tree.findChildByName(field.tag)
			if child != nil {
				childStruct, err := field.filter(fieldValue, child)
				if err != nil {
					return nil, err
				}

				resultMap[field.tag] = childStruct
				continue
			}

			if tree.containsField(field.tag) {
				resultMap[field.tag] = fieldValue.Interface()
			}
		}

		return resultMap, nil
	}
}

func (g *golangQL) fieldParseTag(tag string) string {
	if i := strings.Index(tag, ","); i != -1 {
		name := tag[:i]
		tag = tag[i+1:]
		return name
	}
	return tag
}
