# neighbor
---

A neighborhood watch for testing on GitHub.

## Requirements
+ [Go](https://golang.org/dl/)
    1. `mkdir $HOME/go`

## Getting Started
1. Installing the project
    1. `cd $HOME/go`
    2. `go get -u -v github.com/mccurdyc/neighbor`
2. Start SSH Agent and Add Keys to Agent
    1. Start SSH Agent
    ```
    eval `ssh-agent -s`
    ```
    2. Add SSH Keys (your setup may differ, but generally you can do the following)
    ```
    ssh-add $HOME/.ssh/id_rsa
    ```
3. Run the Setup
    ```bash
    make setup
    ```

    The `setup` `make` target will do the following:
    1. Install `dep`
    2. Backup YOUR installed version of `go` (`which go`) to the value returned
      from `which go` with a file extension `.bak` appended
    3. Move `./bin/go-cover` to `which go` to be used as the system-wide `go` command
      + you can verify this by running `go version` after running `make setup`
          + If you see something similar to the following, then you are still
            running an officially-released version of `go`.
              ```bash
              go version go1.11 linux/amd64
              ```
          + If you see something like the following, then you are running neighbor's `go` version.
              ```bash
              go version devel +2afdd17e3f Mon Oct 8 19:13:38 2018 +0000 linux/amd64
              ```
4. Generate a [Personal Access Token on GitHub](https://github.com/settings/tokens)
    + Add generated token to the GitHub configuration file
      ```json
      {
        "access_token": "1234567890abcdefghijklmnopqrstuvwxyz",
        ...
      }
      ```
5. Run neighbor
    ```bash
    make run
    ```

## License
+ [GNU General Public License Version 3](./LICENSE)
