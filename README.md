<div align="center">
  <img src="https://github.com/mccurdyc/neighbor/blob/master/docs/imgs/orange-background-logo.png?raw=true"><br>
</div>

[![Gitter](https://badges.gitter.im/neighborproject/community.svg)](https://gitter.im/neighborproject/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge) [![codecov](https://codecov.io/gh/mccurdyc/neighbor/branch/master/graph/badge.svg)](https://codecov.io/gh/mccurdyc/neighbor) [![Maintainability](https://api.codeclimate.com/v1/badges/8b473a645aab19597124/maintainability)](https://codeclimate.com/github/mccurdyc/neighbor/maintainability)

neighbor has importable Go packages (e.g., `builtin/*`, `sdk/*`) and an accompanying
command-line interface for searching, cloning and executing an arbitrary binary
against GitHub projects. Abstractions are in place to make doing the aforementioned
easy and efficient for projects obtained from arbitrary search retrieval methods
(i.e., not limited to GitHub Search, repositories or Git clone).

The motivation for neighbor is to provide users (e.g., developers, researchers, etc.)
with a way to search, efficiently clone and evaluate projects without having to
"roll their own". Instead users can focus on the task at hand.

neighbor uses [v3 of GitHub's REST API](https://developer.github.com/v3/).

### Why neighbor

+ Extensibility
  + Abstract interfaces for projects, search and retrieval functions which means
  that it is easy to add new "types" or projects (e.g., something other than GitHub
  repositories) and use other methods for search and retrieval in addition to
  GitHub search and Git clone, respectively.
+ Abstracting GitHub API interaction (searching, sorting and cloning)
  + Transparent pagination
  + Transparent authentication
  + Transparent rate limit handling
+ Doing the above efficiently by leveraging Go's concurrent capabilities

## Requirements

+ [Go `1.13+`](https://golang.org/dl/)
  + Why `1.13+`?
    + [Updates to error handling](https://blog.golang.org/go1.13-errors)
    + [Updates to modules](https://golang.org/doc/go1.13#modules) for dependency management
  + [Installing Go documentation](https://golang.org/doc/install)

## Getting Started

1. Installing the project

    `go get github.com/mccurdyc/neighbor`

2. Searching and Evaluating

    First, you should review the [Searching on GitHub](https://help.github.com/en/github/searching-for-information-on-github/searching-on-github) documentation.

    1. **Repository Search Example**

      ```bash
      make build
      ./bin/neighbor --query="org:neighbor-projects NOT minikube" --external_command="ls -al"
      ```

    2. **Code Search Example**

      _Note: [GitHub requires users to be logged in to search code](https://developer.github.com/v3/search/#search-code).
      Even in public repositories._ Refer to the Code search documentation [here](https://help.github.com/en/github/searching-for-information-on-github/searching-code)
      for building a query. Code searches are searched elastically and are not
      guaranteed to return exact matches. Searching code for exact matches is currently
      in beta and only work on very specific repositories, see [this section in the documentation](https://help.github.com/en/github/searching-for-information-on-github/searching-code-for-exact-matches#searching-code-for-exact-matches)

      It is critical that you read the above documentation because Code search may
      not behave as you would expect. For example,

      > You can't use the following wildcard characters as part of your search query:
      ```
      . , : ; / \ ` ' " = * ! ? # $ & + ^ | ~ < > ( ) { } [ ]
      ```

      The search will simply ignore these symbols. Additionally, I have found that
      using [`extension:EXTENSION`](https://help.github.com/en/github/searching-for-information-on-github/searching-code#search-by-file-extension)
      is more reliable and accurate than [`filename:FILENAME`](https://help.github.com/en/github/searching-for-information-on-github/searching-code#search-by-filename).

      ```bash
      make build
      ./bin/neighbor --search_type="code" --access_token="abc123" --query="pkg/errors in:file extension:mod path:/ user:mccurdyc" --external_command="ls -al"
      ```

3. Confirming

  One way to confirm that you obtained the number of projects that you expected
  is to run the following:

  ```bash
  find _external_projects -mindepth 2 -maxdepth 2 | wc -l
  ```

## Usage

```bash
Usage: neighbor (--file=<config-file> | [--search_type="repository"] [--access_token=<github-access-token>] --query=<github-query> --external_command=<command> | --search_type="code" --access_token=<github-access-token> --query=<github-query> --external_command=<command>) [--clean=<[true|false>]

  -access_token string
        Your personal GitHub access token. This is required to access private repositories and increases rate limits.
  -alsologtostderr
        log to standard error as well as files
  -clean
        Delete the directory created for each repository after running the external command against the repository. (default true)
  -external_command string
        The command to execute on each project returned from the GitHub search query.
  -file string
        Absolute filepath to the config file.
  -help
        Print this help menu.
  -log_backtrace_at value
        when logging hits line file:N, emit a stack trace
  -log_dir string
        If non-empty, write log files in this directory
  -logtostderr
        log to standard error instead of files
  -query string
        The GitHub search query to execute.
  -search_type string
        The type of GitHub search to perform. (default "repository")
  -stderrthreshold value
        logs at or above this threshold go to stderr
  -v value
        log level for V logs
  -vmodule value
        comma-separated list of pattern=N settings for file-filtered logging
```

## FAQ

### What about private repositories?

Generate a [GitHub Personal Access Token](https://github.com/settings/tokens)
neighbor uses token authentication for communicating and authenticating with GitHub.
To read more about GitHub's token authentication, visit [this site](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/).

> You can create a personal access token and use it in place of a password when performing Git operations over HTTPS with Git on the command line or the API.

Authentication is required to both increase the [GitHub API limitations](https://godoc.org/github.com/google/go-github/github#hdr-Rate_Limiting)
as well as access private content (e.g., repositories, gists, etc.).

+ Use the `--access_token` command-line argument
+ If using a config file, add the generated token to the file
  ```json
  {
    "access_token": "yourAccessToken1234567890abcdefghijklmnopqrstuvwxyz",
    ...
  }
  ```

### Executing a Cli Command/Executable Binary

neighbor allows you to specify an executable binary to be run on
a per-repository basis with **each repository as the working directory**.

Examples can be found in the [examples](./_examples).

## License
+ [GNU General Public License Version 3](./LICENSE)
