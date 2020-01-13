package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/golang/glog"

	"github.com/mccurdyc/neighbor/builtin/retrieval/git"
	"github.com/mccurdyc/neighbor/builtin/run/binary"
	"github.com/mccurdyc/neighbor/builtin/search/github"
	"github.com/mccurdyc/neighbor/sdk/retrieval"
	"github.com/mccurdyc/neighbor/sdk/run"
	"github.com/mccurdyc/neighbor/sdk/search"
)

func main() {
	fp := flag.String("file", "", "Absolute filepath to the config file.")
	tkn := flag.String("access_token", "", "Your personal GitHub access token. This is required to access private repositories and increases rate limits.")
	searchType := flag.String("search_type", "project", "The type of search to perform.")
	query := flag.String("query", "", "The search query to execute.")
	command := flag.String("command", "", "The command to execute on each project returned from a search query.")
	projectsDir := flag.String("projects_directory", "_external_projects", "Where the projects should be stored locally and found for evalutation.")
	plainRetrieve := flag.Bool("plain_retrieve", false, "Whether projects should just be retrieved and not evaluated.")
	clean := flag.Bool("clean", true, "Delete the result directory after running the command against every project.")
	help := flag.Bool("help", false, "Print this help menu.")

	flag.Parse()

	if *help == true ||
		(*fp == "" && (*query == "" || *searchType == "")) ||
		(*plainRetrieve == false && *command == "") {
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

	cfg := NewCfg(*fp)

	if len(*fp) != 0 {
		cfg.Parse()

		tkn = &cfg.Contents.AccessToken
		searchType = &cfg.Contents.SearchType
		query = &cfg.Contents.Query
		command = &cfg.Contents.Command
	}

	workingDir, err := os.Getwd()
	if err != nil {
		glog.Exitf("failed to get working directory: %+v", err)
	}

	err = os.Mkdir(*projectsDir, os.ModePerm)
	if err != nil {
		glog.Exitf("failed to create project directory: %+v", err)
	}

	if *clean {
		defer func() {
			err := os.RemoveAll(*projectsDir)
			if err != nil {
				glog.Errorf("error cleaning up: %+v", err)
			}
		}()
	}

	var method uint32
	switch *searchType {
	case "project":
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

	cmd, err := binary.Factory(ctx, &run.BackendConfig{
		Cmd:    *command,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		glog.Exitf("failed to handle command: %+v", err)
	}

	for _, p := range projects {
		dir := filepath.Join(workingDir, *projectsDir, p.Name())
		err := gitClone.Retrieve(ctx, p.SourceLocation(), dir)
		if err != nil {
			glog.Errorf("error retrieving project ('%s): %+v", p.Name(), err)
			continue
		}

		err = cmd.Run(ctx, dir)
		if err != nil {
			glog.Errorf("failed to run binary command in '%s': %+v", dir, err)
		}
	}
}

// usage prints the usage and the supported flags.
func usage() {
	fmt.Fprint(flag.CommandLine.Output(), "\nUsage: neighbor (--file=<config-file> | --query=<github-query> (--command=<command> | --plain_retrieve)) [--access_token=<github-access-token>] [--search_type=<repository|code>] [--clean=<true|false>]\n\n")
	flag.PrintDefaults()
	fmt.Fprint(flag.CommandLine.Output(), "\n")
}
