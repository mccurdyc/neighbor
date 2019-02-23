# neighbor
---

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
+ [Go](https://golang.org/dl/) >= 1.11 (in order to guarantee reproducible builds)

## Executing a Cli Command/Executable Binary
neighbor allows you to specify a cli command or executable binary to be run on
a per-repository basis with **each repository as the working directory**.

Sample custom binaries can be found in the [examples](./_examples).

## Getting Started
1. Installing the project
    1. `go get -u github.com/mccurdyc/neighbor`

2. Generate a [Personal Access Token on GitHub](https://github.com/settings/tokens)
    neighbor uses token authentication for communicating and authenticating with GitHub.
    To read more about GitHub's token authentication, visit [this site](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/).

    > You can create a personal access token and use it in place of a password when performing Git operations over HTTPS with Git on the command line or the API.

    Authentication is required to both increase the [GitHub API limitations](https://godoc.org/github.com/google/go-github/github#hdr-Rate_Limiting)
    as well as access private content (e.g., repositories, gists, etc.).

    + Add the generated token to the configuration file (`config.json`).
      ```json
      {
        "access_token": "yourAccessToken1234567890abcdefghijklmnopqrstuvwxyz",
        ...
      }
      ```
3. Prepare neighbor for Execution
    ```bash
    make
    ```

    This will do the following:
    + Check that you have the appropriate Go version
    + Create a `config.json` file from the `sample.config.json` file
    + You still need to update the access token in the config file to use your personal access token.
    + Enable [Go modules](https://github.com/golang/go/wiki/Modules) by setting `GO111MODULE=on`

4. Executing an external command on each of the GitHub projects returned from the query
    ```bash
    make run
    ```

    The `run` target will first build neighbor and then invoke neighbor pointed
    at the config.json file in the root of the project.

    neighbor will use the GitHub query specified in the config file to find projects
    on GitHub. neighbor will then clone and run the external command in each of the
    projects' directory, sequentially using the command specified in the config, `external_command`.

## License
+ [GNU General Public License Version 3](./LICENSE)
