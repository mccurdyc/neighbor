package main

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parse(t *testing.T) {
	type input struct {
		reader  io.Reader
		content *Contents
	}

	type want struct {
		content Contents
		err     error
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"all_fields": {
			input: input{
				reader: strings.NewReader(`{
															"auth_token": "123abc",
															"search_type": "type",
															"query": "query",
															"command": "hello",
															"plain_retrieve": true,
															"clean": false,
															"projects_directory": "/hello/there",
															"num_projects": 11
														}`),
				content: &Contents{},
			},
			want: want{
				content: Contents{
					AuthToken:     "123abc",
					SearchType:    "type",
					Query:         "query",
					Command:       "hello",
					PlainRetrieve: true,
					Clean:         false,
					ProjectsDir:   "/hello/there",
					NumProjects:   11,
				},
				err: nil,
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			gotErr := parse(tt.input.reader, tt.input.content)

			if diff := cmp.Diff(tt.want.content, *tt.input.content); diff != "" {
				t.Errorf("parse() mismatch (-want +got):\n%s", diff)
			}

			// https://github.com/google/go-cmp/issues/24
			errorCmp := func(x, y error) bool {
				if x == nil || y == nil {
					return x == nil && y == nil
				}
				return x.Error() == y.Error()
			}

			if ok := errorCmp(gotErr, tt.want.err); !ok {
				t.Errorf("parse() \n\tgotErr: '%+v'\n\twantErr: '%+v'", gotErr, tt.want.err)
			}
		})
	}
}
