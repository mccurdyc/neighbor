package main

import (
	// stdlib
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	// external
	"github.com/golang/glog"
	"github.com/pkg/errors"

	// internal
	"github.com/mccurdyc/neighbor/pkg/config"
	"github.com/mccurdyc/neighbor/pkg/external"
	"github.com/mccurdyc/neighbor/pkg/github"
	"github.com/mccurdyc/neighbor/pkg/neighbor"
)

func main() {
	fp := flag.String("file", "", "absolute filepath to config [default: \"$(pwd)/config.json\"].")
	tkn := flag.String("access_token", "", "your personal GitHub access token.")
	searchType := flag.String("search_type", "", "the type of GitHub search to perform.")
	query := flag.String("query", "", "the GitHub search query to execute.")
	externalCmd := flag.String("external_command", "", "the command to execute on each project returned from the GitHub search query.")

	flag.Parse()

	// #TODO - would be nice to be able to override
	wd, err := os.Getwd()
	if err != nil {
		glog.Exitf("error getting current directory: %+v", err)
	}

	c, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create a context object that will be used for the life of the program and passed around
	ctx := &neighbor.Ctx{
		Context: c,
		GitHub: neighbor.GitHubDetails{
			// by default, use the cli args
			// if the config file is specified, these will be overwritten
			AccessToken: *tkn,
			SearchType:  *searchType,
			Query:       *query,
		},
		NeighborDir:  wd,
		ExtResultDir: fmt.Sprintf("%s/%s", wd, "_external-projects-wd"), // go tools handle directories prepended with '_' differently; often they ignore those directories
	}

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

	cmd := *externalCmd

	if len(*fp) != 0 {
		cfg := config.New(*fp)
		cfg.Parse()

		ctx.Config = cfg
		ctx.GitHub = neighbor.GitHubDetails{
			AccessToken: cfg.Contents.AccessToken,
			SearchType:  cfg.Contents.SearchType,
			Query:       cfg.Contents.Query,
		}

		cmd = cfg.Contents.ExternalCmdStr
	}

	ctx.SetExternalCmd(cmd)

	glog.V(2).Infof("external command to be run on each project: %s\n", ctx.ExternalCmd)

	if err = ctx.Validate(); err != nil {
		glog.Exit(errors.Wrap(err, "error validating context"))
	}

	err = ctx.CreateExternalResultDir()
	if err != nil {
		glog.Exitf("error creating results directory: %+v", err)
	}

	svc := github.NewSearchService(github.Connect(ctx.Context, ctx.GitHub.AccessToken))
	res, resp := svc.Search(ctx, ctx.GitHub.SearchType, ctx.GitHub.Query, nil)
	glog.V(3).Infof("github search response: %+v", resp)
	glog.V(2).Infof("github search result: %+v", res)

	ch := github.CloneFromResult(ctx, svc.Client, res)
	external.Run(ctx, ch)
}
