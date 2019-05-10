# neighbor

[![codecov](https://codecov.io/gh/mccurdyc/neighbor/branch/master/graph/badge.svg)](https://codecov.io/gh/mccurdyc/neighbor) [![Maintainability](https://api.codeclimate.com/v1/badges/8b473a645aab19597124/maintainability)](https://codeclimate.com/github/mccurdyc/neighbor/maintainability)

neighbor is a tool for cloning a set of repositories from GitHub specified by a
[GitHub Search Query](https://help.github.com/en/articles/searching-for-repositories)
and running a cli command or executable binary, concurrently.

## Background

neighbor aims to offload the work of cloning a set of repositories and executing
a cli command or executable binary on each of the cloned repositories, so that developers
and researchers can focus on what they are actually trying to accomplish.

### How does neighbor save developers and researchers time?
+ Abstracting GitHub API interaction (searching and cloning)
+ Abstracting concurrency

## Requirements
+ [Go](https://golang.org/dl/)

## Getting Started
1. Installing the project

    `go get -u github.com/mccurdyc/neighbor/...`

2. Usage

    ```bash
    make build
    ./bin/neighbor --query="org:neighbor-projects NOT minikube" --external_command="ls -al"
    ```

## Help Menu

```bash
Usage: neighbor (--file=<config-file> | --query=<github-query> --external_command=<command>) [--access_token=<github-access-token>] [--search_type=<repository|code>] [--clean=<true|false>]

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

## Executing a Cli Command/Executable Binary

neighbor allows you to specify an executable binary to be run on
a per-repository basis with **each repository as the working directory**.

Examples can be found in the [examples](./_examples).

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

## License
+ [GNU General Public License Version 3](./LICENSE)
