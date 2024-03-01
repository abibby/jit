package lodash

import (
	"fmt"
	"strconv"
	"strings"
)

func Get(v any, path string) (any, error) {
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
		return Get(v[currentPart], newPath)
	case []any:
		i, err := strconv.Atoi(currentPart)
		if err != nil {
			return nil, fmt.Errorf("invalid array index %s: %e", currentPart, err)
		}
		return Get(v[i], newPath)
	default:
		return nil, fmt.Errorf("no value at path %s", path)
	}
}

func GetString(v any, path string) (string, error) {
	result, err := Get(v, path)
	if err != nil {
		return "", err
	}
	str, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("the value is not a string")
	}
	return str, nil
}
