package neighbor

import (
	// stdlib
	"reflect"
	"testing"
	// external
	// internal
)

const (
	cmdFailFormatString = "\n\tACTUAL: %#v\tACTUAL LEN: %d\n\tEXPECTED: %#v\tEXPECTED LEN: %d\n"
	failFormatString    = "\n\tACTUAL: %+v\n\tEXPECTED: %+v\n"
)

func TestSetExternalCmd(t *testing.T) {
}

func TestParseCmd(t *testing.T) {
	suite := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single word command",
			input:    "a",
			expected: []string{"a"},
		},
		{
			name:     "space delimited command",
			input:    "a b",
			expected: []string{"a", "b"},
		},
		{
			name:     "with flag",
			input:    "a b -c",
			expected: []string{"a", "b", "-c"},
		},
		{
			name:     "with double quote",
			input:    "a b \"c\"",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with double quote contents space delimited",
			input:    "a b \"c d e f\"",
			expected: []string{"a", "b", "c d e f"},
		},
		{
			name:     "with double quote content after quotes",
			input:    "a b \"c\" d",
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "starts with double quote",
			input:    "\"a\"",
			expected: []string{"a"},
		},
		{
			name:     "starts with double quote multi word",
			input:    "\"a\" b",
			expected: []string{"a", "b"},
		},
		{
			name:     "starts with multi workd double quote",
			input:    "\"a b\" c",
			expected: []string{"a b", "c"},
		},
	}

	for _, c := range suite {
		t.Run(c.name, func(t *testing.T) {
			actual := parseCmd(c.input)

			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf(cmdFailFormatString, actual, len(actual), c.expected, len(c.expected))
			}
		})
	}
}

func TestCleanWords(t *testing.T) {
	suite := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "regex no match",
			input:    []string{"hello there", "world"},
			expected: []string{"hello there", "world"},
		},
		{
			name:     "regex match double quote",
			input:    []string{"\"hello there\"", "\"world\""},
			expected: []string{"hello there", "world"},
		},
		{
			name:     "regex match single quote",
			input:    []string{"'hello there'", "'world'"},
			expected: []string{"hello there", "world"},
		},
	}

	for _, c := range suite {
		t.Run(c.name, func(t *testing.T) {
			actual := cleanWords(c.input)

			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf(cmdFailFormatString, actual, len(actual), c.expected, len(c.expected))
			}
		})
	}
}
