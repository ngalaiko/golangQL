package golangQL

import (
	"sync"
)

type golangQL struct {
	sync.RWMutex
	filters map[cacheKey]filterFunc
}

var instance = &golangQL{
	filters: map[cacheKey]filterFunc{},
}

func Filter(v interface{}, query string) (interface{}, error) {
	return instance.filter(v, query)
}
