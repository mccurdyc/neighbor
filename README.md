# neighbor
---

That friendly neighbor that makes the neighborhood a great place by helping everyone
get things done.

## Requirements
+ [Go](https://golang.org/dl/) >= 1.11
    1. `mkdir $HOME/go`

## Creating an External Command
neighbor allows you to specify an arbitrary command to be run on a per-repository basis.
The command should be executable from the command-line.

Some sample external commands can be found in the [examples](./_examples).

## Getting Started
1. Installing the project
    1. `cd $HOME/go`
    2. `go get -u -v github.com/mccurdyc/neighbor`

2. Prepare neighbor for Execution
    ```bash
    make setup
    ```

    The `setup` `make` target will do the following:
    + Create a `config.json` file from the `sample.config.json` file
      + You still need to update the access token in the config file to use your personal access token.
      + The setup target will check to see if you have already copied the sample.config.json to
        config.json for execution. If you have, the setup will not overwrite the config.json file.

3. Generate a [Personal Access Token on GitHub](https://github.com/settings/tokens)
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
