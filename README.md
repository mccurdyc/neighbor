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
    1. `go get -u -v github.com/mccurdyc/neighbor/...`

2. Generate a [Personal Access Token on GitHub](https://github.com/settings/tokens)
    neighbor uses token authentication for communicating and authenticating with GitHub.
    To read more about GitHub's token authentication, visit [this site](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/).

    > You can create a personal access token and use it in place of a password when performing Git operations over HTTPS with Git on the command line or the API.

    Authentication is required to both increase the [GitHub API limitations](https://godoc.org/github.com/google/go-github/github#hdr-Rate_Limiting)
    as well as access private content (e.g., repositories, gists, etc.).

    + If using a config file, add the generated token to the file
      ```json
      {
        "access_token": "yourAccessToken1234567890abcdefghijklmnopqrstuvwxyz",
        ...
      }
      ```
    + If not using a config file, use the `--access_token` command-line argument

3. Usage
```bash
$ neighbor -h
Usage of neighbor:
  -access_token string
        your personal GitHub access token.
  -external_command string
        the command to execute on each project returned from the GitHub search query.
  -file string
        absolute filepath to config [default: "$(pwd)/config.json"].
  -query string
        the GitHub search query to execute.
  -search_type string
        the type of GitHub search to perform.
```

  Example:
  ```bash
  export GITHUB_ACCESS_TOKEN="your-token-here"
  neighbor --access_token=$GITHUB_ACCESS_TOKEN --search_type="repository" --query="org:neighbor-projects NOT minikube" --external_command="ls -al"
  ```

  This will create a directory `_external-projects-wd` wherever you run `neighbor`
  with the cloned contents of the repositories.

  If you just want the output from the external command, pipe `neighbor` to a file (`neighbor ... > out.txt`).
  The logging of neighbor should be separate from the results of the command.

## License
+ [GNU General Public License Version 3](./LICENSE)
