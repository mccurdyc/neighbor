package github

import (
	"testing"

	"github.com/mccurdyc/neighbor/sdk/project"
	"github.com/mccurdyc/neighbor/sdk/retrieval"
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

type mockProject struct {
	name           string
	version        string
	retrievalFunc  retrieval.Backend
	sourceLocation string
	localLocation  string
}

func (m *mockProject) Name() string                     { return m.name }
func (m *mockProject) Version() string                  { return m.version }
func (m *mockProject) RetrievalFunc() retrieval.Backend { return m.retrievalFunc }
func (m *mockProject) SourceLocation() string           { return m.sourceLocation }
func (m *mockProject) LocalLocation() string            { return m.localLocation }
func (m *mockProject) SetLocalLocation(l string) project.Backend {
	m.localLocation = l
	return m
}

func Test_contains(t *testing.T) {
	type input struct {
		projects []project.Backend
		p        project.Backend
	}

	type want struct {
		value bool
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"project_contained": {
			input: input{
				projects: []project.Backend{
					&mockProject{name: "one/one"},
					&mockProject{name: "two/two"},
					&mockProject{name: "three/three"},
				},
				p: &mockProject{name: "two/two"},
			},
			want: want{value: true},
		},

		"project_not_contained": {
			input: input{
				projects: []project.Backend{
					&mockProject{name: "one/one"},
					&mockProject{name: "two/two"},
					&mockProject{name: "three/three"},
				},
				p: &mockProject{name: "four/four"},
			},
			want: want{value: false},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := contains(tt.input.projects, tt.input.p)

			if got != tt.want.value {
				t.Errorf("contains(%+v): \n\tgot: '%+v'\n\twant: '%+v'", tt.input, got, tt.want.value)
			}
		})
	}
}
