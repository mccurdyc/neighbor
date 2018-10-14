# neighbor
---
A neighborhood watch tool for evaluating neighbors' test suite adequacy in the neighborhood, GitHub projects.

## Requirements
+ [Go](https://golang.org/dl/)
    1. `mkdir $HOME/go`

## Getting Started
1. Installing the project
    1. `cd $HOME/go`
    2. `go get -u -v github.com/mccurdyc/neighbor`
2. Prepare neighbor for Execution
    ```bash
    make setup
    ```

    The `setup` `make` target will do the following:
    1. Install `dep`
    2. Create a `config.json` file from the `sample.config.json` file

      **NOTE: You still need to update the access token in the config file to use your personal access token.**

      **NOTE: The setup target will check to see if you have already copied the sample.config.json to
      config.json for execution. If you have, the setup will not overwrite the config.json file.**
    3. Backup your installed version of `go` (`$ which go`) to the value returned
      from `$ which go` in your shell with a file extension `.bak` appended.

      **NOTE: The setup target will check to see if you have already backed up a go command. If you have,
      the setup will not overwrite the backup.**
    4. Move `./bin/go-cover` to `which go` to be used as the system-wide `go` command
      + You can verify this by running `go version` after running `make setup`
          + If you see something similar to the following, then you are still
            running an officially-released version of `go`.
              ```bash
              go version go1.11 linux/amd64
              ```
          + If you see something like the following, then you are running neighbor's `go` version.
              ```bash
              go version devel +2afdd17e3f Mon Oct 8 19:13:38 2018 +0000 linux/amd64
              ```
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
4. Performing a neighborhood analysis with neighbor
    ```bash
    make run
    ```

    The `run` target will first build neighbor and then invoke neighbor pointed
    at the config.json file in the root of the project.

    neighbor will use the GitHub query specified in the config file to find projects
    on GitHub. neighbor will then clone and test each of these projects sequentially
    using the test command specified in the config. neighbor will use a custom go binary
    instead of the default go binary. This custom go binary always enables the
    `-coverprofile` flag during `go test` and writes to a easy-to-find location
    at the root of each project in the `_ext-results/` directory.

## License
+ [GNU General Public License Version 3](./LICENSE)
