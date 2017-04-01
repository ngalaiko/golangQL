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
	filter := g.getFilter(val, tree)

	result, err := filter(val, tree)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (g *golangQL) getFilter(val reflect.Value, tree *node) filterFunc {
	key := newCacheKey(val.Type(), tree)

	g.RLock()
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

	f = g.newFilter(val.Type(), tree)
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
	elemFilter := g.newFilter(typ.Elem(), tree)

	return func(val reflect.Value, tree *node) (interface{}, error) {
		v := val.Elem()

		return elemFilter(v, tree)
	}
}

func (g *golangQL) newSliceFilter(typ reflect.Type, tree *node) filterFunc {
	return func(val reflect.Value, tree *node) (interface{}, error) {
		resultMap := make([]map[string]interface{}, val.Len())
		resultMapValue := reflect.ValueOf(resultMap)
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)

			elementType := v.Type()
			elemFilter := g.newFilter(elementType, tree)

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

	return func(val reflect.Value, tree *node) (interface{}, error) {
		typ := val.Type()
		resultMap := make(map[string]interface{}, typ.NumField())

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			tag, ok := field.Tag.Lookup("json")
			tagName := g.fieldParseTag(tag)

			if !ok {
				elemFilter := g.newFilter(field.Type, tree)

				embeded, err := elemFilter(val.Field(i), tree)
				if err != nil {
					return nil, err
				}

				if reflect.TypeOf(embeded).Kind() != reflect.Map {
					resultMap[tagName] = val
					continue
				}

				embededMapValue := reflect.ValueOf(embeded)
				for _, key := range embededMapValue.MapKeys() {
					resultMap[key.String()] = embededMapValue.MapIndex(key).Interface()
				}
				continue
			}

			child := tree.findChildByName(tagName)
			if child != nil {
				elemFilter := g.newFilter(typ, child)

				childStruct, err := elemFilter(val.Field(i).Elem(), child)
				if err != nil {
					return nil, err
				}

				resultMap[tagName] = childStruct
				continue
			}

			if tree.containsField(tagName) {
				resultMap[tagName] = val.Field(i).Interface()
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
