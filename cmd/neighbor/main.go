package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"gopkg.in/src-d/go-git.v4"
	"github.com/pkg/errors"

	"github.com/mccurdyc/neighbor/pkg/external"
	"github.com/mccurdyc/neighbor/pkg/github"
)

const projectDir = "_external_project"

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

	// if len(*fp) != 0 {
	// 	cfg := config.New(*fp)
	// 	cfg.Parse()
	//
	// 	ctx.Config = cfg
	// 	ctx.GitHub = neighbor.GitHubDetails{
	// 		AccessToken: cfg.Contents.AccessToken,
	// 		SearchType:  cfg.Contents.SearchType,
	// 		Query:       cfg.Contents.Query,
	// 	}
	//
	// 	cmd = cfg.Contents.ExternalCmdStr
	// }

	err := os.Mkdir(projectDir, os.ModePerm)
	if err != nil {
		glog.Exitf("failed to create collated project directory: %+v", err)
	}

	searcher, err := github.NewSearcher(github.Connect(ctx, *tkn), github.SearchType(*searchType))
	if err != nil {
		glog.Exitf("error creating searcher: %+v", err)
	}

	numDesiredResults := 100 // TODO: read the number of desired results from a config value
	repositories, err := github.Search(ctx, searcher, *query, github.SearchOptions().WithNumberOfResults(numDesiredResults))
	if err != nil {
		glog.V(2).Infof("error searching GitHub: %+v", err)
	}

	for _, repo := range repositories {
		err := github.Clone(ctx, projectDir, *repo, git.CloneOptions{})
		if err != nil {
			glog.V(2).Infof("failed to clone repository (%s): %+v", repo.GetName, err)
		}

		out, err := external.Run(ctx, *externalCmd, projectDir)
	}

	if *clean {
		err := os.RemoveAll(projectDir)
		if err != nil {
			glog.Errorf("error removing directory: %s", r.Directory)
		}
	}
}

// usage prints the usage and the supported flags.
func usage() {
	fmt.Fprint(flag.CommandLine.Output(), "\nUsage: neighbor (--file=<config-file> | --query=<github-query> --external_command=<command>) [--access_token=<github-access-token>] [--search_type=<repository|code>] [--clean=<true|false>]\n\n")
	flag.PrintDefaults()
	fmt.Fprint(flag.CommandLine.Output(), "\n")
}
