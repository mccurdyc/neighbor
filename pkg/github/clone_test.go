package github

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var cloneCfgFixture = CloneConfig{
	url: "cloneurl",
	auth: &http.BasicAuth{
		Username: "username",
		Password: "password",
	},
}

func Test_NewCloneConfig(t *testing.T) {
	got := NewCloneConfig()
	want := CloneConfig{}

	compareCloneConfig(t, "NewCloneConfig", got, want)
}

func Test_WithRepoURL(t *testing.T) {
	repo := &github.Repository{
		CloneURL: ptrToString("new"),
	}

	type input struct {
		cf   CloneConfig
		repo *github.Repository
	}

	type want struct {
		cf CloneConfig
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"empty_config": {
			input: input{
				cf:   CloneConfig{},
				repo: repo,
			},
			want: want{
				cf: CloneConfig{url: "new"},
			},
		},

		"overwrite_url": {
			input: input{
				cf:   cloneCfgFixture,
				repo: repo,
			},
			want: want{
				cf: CloneConfig{url: "new", auth: &http.BasicAuth{Username: "username", Password: "password"}},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.input.cf.WithRepoURL(tt.input.repo)

			compareCloneConfig(t, "WithRepoURL", got, tt.want.cf)
		})
	}
}

func Test_WithTokenAuth(t *testing.T) {
	type input struct {
		cf    CloneConfig
		token string
	}

	type want struct {
		cf CloneConfig
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"empty_config": {
			input: input{
				cf:    CloneConfig{},
				token: "tkn",
			},
			want: want{
				cf: CloneConfig{auth: &http.BasicAuth{Username: "null", Password: "tkn"}},
			},
		},

		"overwrite_auth": {
			input: input{
				cf:    cloneCfgFixture,
				token: "tkn",
			},
			want: want{
				cf: CloneConfig{url: "cloneurl", auth: &http.BasicAuth{Username: "null", Password: "tkn"}},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.input.cf.WithTokenAuth(tt.input.token)

			compareCloneConfig(t, "WithTokenAuth", got, tt.want.cf)
		})
	}
}

func Test_WithBasicAuth(t *testing.T) {
	type input struct {
		cf       CloneConfig
		username string
		password string
	}

	type want struct {
		cf CloneConfig
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"empty_config": {
			input: input{
				cf:       CloneConfig{},
				username: "newusername",
				password: "newpassword",
			},
			want: want{
				cf: CloneConfig{auth: &http.BasicAuth{Username: "newusername", Password: "newpassword"}},
			},
		},

		"overwrite_auth": {
			input: input{
				cf:       cloneCfgFixture,
				username: "newusername",
				password: "newpassword",
			},
			want: want{
				cf: CloneConfig{url: "cloneurl", auth: &http.BasicAuth{Username: "newusername", Password: "newpassword"}},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.input.cf.WithBasicAuth(tt.input.username, tt.input.password)

			compareCloneConfig(t, "WithBasicAuth", got, tt.want.cf)
		})
	}
}

func compareCloneConfig(t *testing.T, fnName string, got, want CloneConfig) {
	if diff := cmp.Diff(want.auth, got.auth, cmp.AllowUnexported()); diff != "" {
		t.Errorf("%s() mismatch (-want +got):\n%s", fnName, diff)
	}

	if diff := cmp.Diff(want.url, got.url, cmp.AllowUnexported()); diff != "" {
		t.Errorf("%s() mismatch (-want +got):\n%s", fnName, diff)
	}
}
