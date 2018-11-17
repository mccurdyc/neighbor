package neighbor

import (
	"reflect"
	"testing"
)

const (
	cmdFailFormatString = "\n\tACTUAL: %#v\tACTUAL LEN: %d\n\tEXPECTED: %#v\tEXPECTED LEN: %d\n"
	failFormatString    = "\n\tACTUAL: %+v\n\tEXPECTED: %+v\n"
)

func TestSetExternalCmd(t *testing.T) {
}

func TestParseCmd(t *testing.T) {
	suite := []struct {
		name        string
		input       string
		expected    []string
		expectedErr error
	}{
		{
			name:        "single word command",
			input:       "a",
			expected:    []string{"a"},
			expectedErr: nil,
		},
		{
			name:        "space delimited command",
			input:       "a b",
			expected:    []string{"a", "b"},
			expectedErr: nil,
		},
		{
			name:        "with flag",
			input:       "a b -c",
			expected:    []string{"a", "b", "-c"},
			expectedErr: nil,
		},
		{
			name:        "with double quote",
			input:       "a b \"c\"",
			expected:    []string{"a", "b", "c"},
			expectedErr: nil,
		},
		{
			name:        "with double quote contents space delimited",
			input:       "a b \"c d e f\"",
			expected:    []string{"a", "b", "c d e f"},
			expectedErr: nil,
		},
		{
			name:        "with double quote content after quotes",
			input:       "a b \"c\" d",
			expected:    []string{"a", "b", "c", "d"},
			expectedErr: nil,
		},
		{
			name:        "starts with double quote",
			input:       "\"a\"",
			expected:    nil,
			expectedErr: errInvalidCmd,
		},
	}

	for _, c := range suite {
		t.Run(c.name, func(t *testing.T) {
			actual, actualErr := parseCmd(c.input)

			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf(cmdFailFormatString, actual, len(actual), c.expected, len(c.expected))
			}

			if (actualErr == nil && c.expectedErr != nil) || (actualErr != nil && c.expectedErr == nil) {
				t.Errorf(failFormatString, actualErr, c.expectedErr)
			}
		})
	}
}

func TestSplitCmd(t *testing.T) {
	suite := []struct {
		name       string
		inputS     string
		inputChars []rune
		expected   []string
	}{
		{
			name:       "no split quote char",
			inputS:     "a",
			inputChars: []rune{34},
			expected:   []string{"a"},
		},
		{
			name:       "split quote char",
			inputS:     "a \"b\" c",
			inputChars: []rune{34},
			expected:   []string{"a", "b", "c"},
		},
		{
			name:       "split quote char multiple following",
			inputS:     "a \"b\" c d",
			inputChars: []rune{34},
			expected:   []string{"a", "b", "c d"},
		},
		{
			name:       "split quote char multiple contained",
			inputS:     "a \"b c d\" e f",
			inputChars: []rune{34},
			expected:   []string{"a", "b c d", "e f"},
		},
	}

	for _, c := range suite {
		t.Run(c.name, func(t *testing.T) {
			actual := splitCmd(c.inputS, c.inputChars...)

			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf(cmdFailFormatString, actual, len(actual), c.expected, len(c.expected))
			}
		})
	}
}
