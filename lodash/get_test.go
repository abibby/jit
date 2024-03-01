package lodash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	testCases := []struct {
		In   any
		Path string
		Out  any
	}{
		{
			In:   map[string]any{"foo": "bar"},
			Path: "foo",
			Out:  "bar",
		},
		{
			In: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			Path: "foo.bar",
			Out:  "baz",
		},
		{
			In: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			Path: "foo",
			Out: map[string]any{
				"bar": "baz",
			},
		},
		{
			In: []any{
				"foo",
				"bar",
				"baz",
			},
			Path: "1",
			Out:  "bar",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Path, func(t *testing.T) {

			result, err := Get(tc.In, tc.Path)
			assert.NoError(t, err)
			assert.Equal(t, tc.Out, result)
		})
	}

}
