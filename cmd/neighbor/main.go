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

	"github.com/mccurdyc/neighbor/pkg/config"
	"github.com/mccurdyc/neighbor/pkg/github"
	"github.com/mccurdyc/neighbor/pkg/runner"
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

	err = os.Mkdir(projectDir, os.ModePerm)
	if err != nil {
		glog.Exitf("failed to create collated project directory: %+v", err)
	}

	searcher, err := github.NewSearcher(github.Connect(ctx, *tkn), github.SearchType(*searchType))
	if err != nil {
		glog.Exitf("error creating searcher: %+v", err)
	}

	numDesiredResults := 10 // TODO: read the number of desired results from a config value
	repositories, err := github.Search(ctx, searcher, *query, github.SearchOptions().WithNumberOfResults(numDesiredResults))
	if err != nil {
		glog.Errorf("error searching GitHub: %+v", err)
	}

	dir := filepath.Join(workingDir, projectDir)
	doneCh := github.CloneRepositories(ctx, dir, repositories, &github.PlainCloner{}, github.NewCloneConfig().WithTokenAuth(*tkn))

	for info := range doneCh {
		if info.Error != nil {
			glog.Errorf("error cloning repository: %+v", info.Error)
			continue
		}

		// Right now, commands are sequentially run. The only part that is done concurrently
		// is cloning.
		err = runBinary(info.Meta.ClonedDir, *externalCmd)
		if err != nil {
			glog.Errorf("failed to run binary command in '%s': %+v", info.Meta.ClonedDir, err)
		}
	}

	if *clean {
		err := os.RemoveAll(projectDir)
		if err != nil {
			glog.Errorf("error cleaning up: %w", err)
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
