<div align="center">
  <img src="https://github.com/mccurdyc/neighbor/blob/master/docs/imgs/orange-background-logo.png?raw=true"><br>
</div>

[![Build Status][build-badge]][build-url]
[![GolangCI][golint-badge]][golint-url]
[![GoDoc][godoc-badge]][godoc-url]
[![License][license-badge]][license-url]
[![codecov][codecov-badge]][codecov-url]
[![Discord chat][discord-badge]][discord-url]
[![Gitter][gitter-badge]][gitter-url]
[![Release][release-badge]][release-url]

[build-badge]: https://circleci.com/gh/mccurdyc/neighbor/tree/master.svg?style=svg
[build-url]: https://circleci.com/gh/mccurdyc/neighbor/tree/master
[golint-badge]: https://golangci.com/badges/github.com/mccurdyc/neighbor.svg
[golint-url]: https://golangci.com
[godoc-badge]: https://godoc.org/github.com/mccurdyc/neighbor?status.svg
[godoc-url]: https://pkg.go.dev/github.com/mccurdyc/neighbor?tab=overview
[license-badge]: https://img.shields.io/github/license/mccurdyc/neighbor
[license-url]: LICENSE
[codecov-badge]: https://codecov.io/gh/mccurdyc/neighbor/branch/master/graph/badge.svg
[codecov-url]: https://codecov.io/gh/mccurdyc/neighbor
[discord-badge]: https://img.shields.io/discord/666244141784498177?logo=discord&label=discord&logoColor=white
[discord-url]: https://discord.gg/qq9sA7
[gitter-badge]: https://badges.gitter.im/neighborproject/community.svg
[gitter-url]: https://gitter.im/neighborproject/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge
[release-badge]: https://img.shields.io/github/release/mccurdyc/neighbor.svg
[release-url]: https://github.com/mccurdyc/neighbor/releases/latest

neighbor has importable Go packages (e.g., `builtin/*`, `sdk/*`) and an accompanying
command-line interface for searching, cloning and executing an arbitrary binary
against GitHub projects. Abstractions are in place to make doing the aforementioned
easy and efficient for projects obtained from arbitrary search and retrieval methods
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

    `GOPROXY=https://proxy.golang.org go get github.com/mccurdyc/neighbor@latest`

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

    3. **Multi-Line Command Example**

      Multi-line commands work, but **pipes (i.e., `|`) do not**. In order to use pipes,
      you should create a custom binary that handles piping the output from one command
      to the next (e.g., ["How to pipe several comands in Go?" on StackOverflow](https://stackoverflow.com/questions/10781516/how-to-pipe-several-commands-in-go))

      ```bash
      make build
      ./bin/neighbor --search_type="code" --access_token="abc123" --query="pkg/errors in:file extension:mod path:/ user:mccurdyc" --external_command="ls \
      -al"
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

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fmccurdyc%2Fneighbor.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fmccurdyc%2Fneighbor?ref=badge_large)
