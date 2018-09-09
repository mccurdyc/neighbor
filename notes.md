1. collaborator [invite link](https://github.com/mccurdyc/neighbor/invitations) to private repo github.com/mccurdc/neighbor
2. golang modules are new, [here](https://github.com/golang/go/wiki/Modules#installing-and-activating-module-support) is a guide
  + these allow for you to work outside of the infamous GOPATH
  + way to version dependencies for reproducibility
  + [other go 1.11 release notes](https://blog.golang.org/go1.11)
3. project directory structure
  + we will adopt some, *but not all*, of the recommendations made in [this github repo](https://github.com/golang-standards/project-layout)
      + for example, we will not use an `internal` directory, unless there is a strong reason to do so
