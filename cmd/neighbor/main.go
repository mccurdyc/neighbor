package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/golang/glog"

	"github.com/mccurdyc/neighbor/builtin/retrieval/git"
	"github.com/mccurdyc/neighbor/builtin/search/github"
	"github.com/mccurdyc/neighbor/pkg/config"
	"github.com/mccurdyc/neighbor/pkg/runner"
	"github.com/mccurdyc/neighbor/sdk/retrieval"
	"github.com/mccurdyc/neighbor/sdk/search"
)

const neighborDir = "_external_projects"

func main() {
	fp := flag.String("file", "", "Absolute filepath to the config file.")
	tkn := flag.String("access_token", "", "Your personal GitHub access token. This is required to access private repositories and increases rate limits.")
	searchType := flag.String("search_type", "repository", "The type of GitHub search to perform.")
	query := flag.String("query", "", "The GitHub search query to execute.")
	externalCmd := flag.String("external_command", "", "The command to execute on each project returned from the GitHub search query.")
	clean := flag.Bool("clean", true, "Delete the directory created for each repository after running the external command against the repository.")
	help := flag.Bool("help", false, "Print this help menu.")

	flag.Parse()

	if *help == true ||
		(*fp == "" && (*query == "" || *externalCmd == "" || *searchType == "")) {
		usage()
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// listen for signals such as SIGINT (^C, CONTROL-C)
	go func() {
		ch := make(chan os.Signal, 1)

		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(ch)

		select {
		case <-ch:
			cancel()
			os.Exit(130)
		}
	}()

	cfg := config.New(*fp)

	if len(*fp) != 0 {
		cfg.Parse()

		tkn = &cfg.Contents.AccessToken
		searchType = &cfg.Contents.SearchType
		query = &cfg.Contents.Query
		externalCmd = &cfg.Contents.ExternalCmdStr
	}

	workingDir, err := os.Getwd()
	if err != nil {
		glog.Exitf("failed to get working directory: %+v", err)
	}

	err = os.Mkdir(neighborDir, os.ModePerm)
	if err != nil {
		glog.Exitf("failed to create project directory: %+v", err)
	}

	var method uint32
	switch *searchType {
	case "repository":
		method = search.Project
	case "code":
		method = search.Code
	default:
		glog.Exit("unsupported search type")
	}

	searchConfig := search.BackendConfig{
		SearchMethod: search.Method(method),
	}

	if len(*tkn) != 0 {
		searchConfig.AuthMethod = "token"
		searchConfig.Config = map[string]string{"token": *tkn}
	}

	githubSearch, err := github.Factory(ctx, &searchConfig)
	if err != nil {
		glog.Exitf("failed to create GitHub searcher: %+v", err)
	}

	numDesiredResults := 10 // TODO: make configurable
	projects, err := githubSearch.Search(context.TODO(), *query, numDesiredResults)
	if err != nil {
		glog.Errorf("encountered error while searching GitHub for projects: %+v", err)
	}

	var retrievalConfig retrieval.BackendConfig
	if len(*tkn) != 0 {
		retrievalConfig.AuthMethod = "token"
		retrievalConfig.Config = map[string]string{"token": *tkn}
	}

	gitClone, err := git.Factory(ctx, &retrievalConfig)
	if err != nil {
		glog.Exitf("error creating Git project retriever: %+v", err)
	}

	for _, p := range projects {
		dir := filepath.Join(workingDir, neighborDir, p.Name())
		err := gitClone.Retrieve(ctx, p.SourceLocation(), dir)
		if err != nil {
			glog.Errorf("error retrieving project ('%s): %+v", p.Name(), err)
			continue
		}

		err = runBinary(dir, *externalCmd)
		if err != nil {
			glog.Errorf("failed to run binary command in '%s': %+v", dir, err)
		}
	}

	if *clean {
		err := os.RemoveAll(neighborDir)
		if err != nil {
			glog.Errorf("error cleaning up: %+v", err)
		}
	}
}

func runBinary(dir string, command string) error {
	xc := strings.Split(command, " ")

	cmd := exec.Command(xc[0])
	if len(xc) > 1 {
		cmd = exec.Command(xc[0], xc[1:]...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return runner.RunInDir(dir, cmd)
}

// usage prints the usage and the supported flags.
func usage() {
	fmt.Fprint(flag.CommandLine.Output(), "\nUsage: neighbor (--file=<config-file> | --query=<github-query> --external_command=<command>) [--access_token=<github-access-token>] [--search_type=<repository|code>] [--clean=<true|false>]\n\n")
	flag.PrintDefaults()
	fmt.Fprint(flag.CommandLine.Output(), "\n")
}
