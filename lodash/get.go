package lodash

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Get[T any](v any, path string) (T, error) {
	result, err := get(v, path)
	if err != nil {
		var zero T
		return zero, err
	}
	typedResult, ok := result.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("the value is %s expected %s", reflect.TypeOf(result), reflect.TypeOf(zero))
	}
	return typedResult, nil
}
func get(v any, path string) (any, error) {
	if path == "" {
		return v, nil
	}
	parts := strings.SplitN(path, ".", 2)
	currentPart := parts[0]
	newPath := ""
	if len(parts) > 1 {
		newPath = parts[1]
	}
	switch v := v.(type) {
	case map[string]any:
		return get(v[currentPart], newPath)
	case []any:
		i, err := strconv.Atoi(currentPart)
		if err != nil {
			return nil, fmt.Errorf("invalid array index %s: %e", currentPart, err)
		}
		return get(v[i], newPath)
	default:
		return nil, fmt.Errorf("no value at path %s", path)
	}
}
