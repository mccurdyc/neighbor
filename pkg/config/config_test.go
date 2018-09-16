package config

import (
	"io"
	"strings"
	"testing"
)

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
				AccessToken: "1123abc",
				SearchType:  "abc",
				Query:       []byte("abc"),
			},
			expectedErr: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actualErr := parse(c.r, c.d)

			// we should switch on type, cast to type and compare values
			if actualErr != c.expectedErr {
				t.Errorf("\tACTUAL: %+v\n\tEXPECTED: %+v\n\tACTUAL ERROR: %+v\n\tEXPECTED ERROR: %+v\n", c.d, c.expected, actualErr, c.expectedErr)
			}
		})
	}
}
