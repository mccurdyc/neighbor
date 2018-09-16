package config

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"
)

func compareContentsFields(a, b Contents) bool {
	return a.AccessToken != b.AccessToken ||
		a.SearchType != b.SearchType ||
		!bytes.Equal(a.Query, b.Query)
}

func Test_parse(t *testing.T) {
	cases := []struct {
		name        string
		r           io.Reader
		d           interface{}
		expected    interface{}
		expectedErr error
	}{
		{
			name: "base",
			r: strings.NewReader(`{
															"access_token": "123abc",
															"search_type": "abc",
															"query": "abc"
														}`),
			d: &Contents{},
			expected: Contents{
				AccessToken: "123abc",
				SearchType:  "abc",
				Query:       json.RawMessage(`"abc"`),
			},
			expectedErr: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actualErr := parse(c.r, c.d)

			switch c.d.(type) {
			case *Contents:
				if actualErr != c.expectedErr || compareContentsFields(*c.d.(*Contents), c.expected.(Contents)) {
					t.Errorf("\n\tACTUAL: %+v\n\tEXPECTED: %+v\n\tACTUAL ERROR: %+v\n\tEXPECTED ERROR: %+v\n", c.d, c.expected, actualErr, c.expectedErr)
				}
			}
		})
	}
}
