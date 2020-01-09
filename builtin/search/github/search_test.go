package github

import (
	"testing"
)

type mockURLRetriever struct {
	cloneURL string
	htmlURL  string
}

func (m mockURLRetriever) GetCloneURL() string { return m.cloneURL }
func (m mockURLRetriever) GetHTMLURL() string  { return m.htmlURL }

func Test_getCloneURL(t *testing.T) {
	type input struct {
		cloneURL string
		htmlURL  string
	}

	type want struct {
		value string
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"empty_clone_url_use_fallback": {
			input: input{
				cloneURL: "",
				htmlURL:  "htmlurl",
			},
			want: want{
				value: "htmlurl.git",
			},
		},

		"clone_url": {
			input: input{
				cloneURL: "cloneURL",
				htmlURL:  "htmlurl",
			},
			want: want{
				value: "cloneURL",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := mockURLRetriever{
				cloneURL: tt.input.cloneURL,
				htmlURL:  tt.input.htmlURL,
			}

			got := getCloneURL(r)

			if got != tt.want.value {
				t.Errorf("getCloneURL(%+v): \n\tgot: '%+v'\n\twant: '%+v'", tt.input, got, tt.want.value)
			}
		})
	}
}
