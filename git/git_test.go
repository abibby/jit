package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveRepeats(t *testing.T) {
	testCases := []struct {
		Name     string
		Str      string
		Expected string
	}{
		{
			Name:     "no remove",
			Str:      "abc123",
			Expected: "abc123",
		},
		{
			Name:     "extra spaces",
			Str:      "abc   123",
			Expected: "abc 123",
		},
		{
			Name:     "extra spaces",
			Str:      "abc   123",
			Expected: "abc 123",
		},
		{
			Name:     "extra dashes",
			Str:      "abc-----123",
			Expected: "abc-123",
		},
		{
			Name:     "extra letters",
			Str:      "aabbcc    112233",
			Expected: "aabbcc 112233",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, removeRepeatDashes(tc.Str))
		})
	}
}

func TestPrepBranchName(t *testing.T) {
	testCases := []struct {
		Name     string
		Str      string
		Expected string
	}{
		{
			Name:     "no remove",
			Str:      "abc123",
			Expected: "abc123",
		},
		{
			Name:     "extra spaces",
			Str:      "abc   123",
			Expected: "abc-123",
		},
		{
			Name:     "extra dashes",
			Str:      "abc-----123",
			Expected: "abc-123",
		},
		{
			Name:     "extra letters",
			Str:      "aabbcc    112233",
			Expected: "aabbcc-112233",
		},
		{
			Name:     "dots",
			Str:      "a.b.c",
			Expected: "a-b-c",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, PrepBranchName(tc.Str))
		})
	}
}
