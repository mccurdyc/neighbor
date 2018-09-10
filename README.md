# neighbor
---

A neighborhood watch for testing on GitHub.

## Requirements
+ [Go](https://golang.org/dl/)
    1. `mkir $HOME/go`
+ [dep](https://github.com/golang/dep)
    1. `cd $HOME/go`
    2. `go get -u -v github.com/golang/dep`
    3. `go install -v github.com/golang/dep/cmd/dep`

## Getting Started
1. Installing the project
    1. `cd $HOME/go`
    2. `go get -u -v github.com/mccurdyc/neighbor`
2. Installing dependencies
    ```bash
    dep ensure -v
    ```
3. Generate a [Personal Access Token on GitHub](https://github.com/settings/tokens)
    + Add generated token token to config file
      ```json
      {
        "access_token": "1234567890abcdefghijklmnopqrstuvwxyz",
        ...
      }
      ```
4. Run neighbor
    ```bash
    make run
    ```

## License

+ [GNU General Public License Version 3](./LICENSE)
