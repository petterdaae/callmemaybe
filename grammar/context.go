package grammar

import "fmt"

type Context struct {
	values map[string]int
}

func (context *Context) Get(key string) (int, error) {
	value, ok := context.values[key]
	if ok {
		return value, nil
	}
	return 0, fmt.Errorf("uknown identifier: %s", key)
}

func (context *Context) CopyWith(key string, value int) Context {
	var copy map[string]int
	for oldKey, oldVal := range context.values {
		copy[oldKey] = oldVal
	}
	copy[key] = value
	return Context{
		values: copy,
	}
}
