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
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	// internal
	"github.com/mccurdyc/neighbor/pkg/config"
	"github.com/mccurdyc/neighbor/pkg/external"
	"github.com/mccurdyc/neighbor/pkg/github"
	"github.com/mccurdyc/neighbor/pkg/neighbor"
)

func main() {
	fp := flag.String("file", "", "absolute filepath to config [default: \"$(PWD)/config.json\"].")
	tkn := flag.String("access_token", "", "your personal GitHub access token.")
	searchType := flag.String("search_type", "", "the type of GitHub search to perform.")
	query := flag.String("query", "", "the GitHub search query to execute.")
	externalCmd := flag.String("external_command", "", "the command to execute on each project returned from the GitHub search query.")

	flag.Parse()

	l := log.New()

	wd, err := os.Getwd()
	if err != nil {
		l.Errorf("error getting current directory: %+v", err)
		os.Exit(1)
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
		Logger:       l,
		NeighborDir:  wd,
		ExtResultDir: fmt.Sprintf("%s/%s", wd, "_external-projects-wd"), // go tools handle directories prepended with '_' differently; often they ignore those directories
	}

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

	if err = ctx.SetExternalCmd(cmd); err != nil {
		err = errors.Wrap(err, "error parsing external command from config")
		ctx.Logger.Error(err)
		os.Exit(1)
	}
	ctx.Logger.Infof("external command to be run on each project: %s\n", ctx.ExternalCmd)

	if err = ctx.Validate(); err != nil {
		err = errors.Wrap(err, "error validating context")
		ctx.Logger.Error(err)
		os.Exit(1)
	}

	go func() {
		ch := make(chan os.Signal, 1)

		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(ch)

		select {
		case <-c.Done():
		case <-ch:
			cancel()
		}
	}()

	ll := os.Getenv("LOG_LEVEL")
	if len(ll) == 0 {
		ctx.Logger.SetLevel(log.InfoLevel)
	} else {
		ll, err := log.ParseLevel(ll)
		if err != nil {
			ctx.Logger.SetLevel(log.InfoLevel)
		}
		ctx.Logger.SetLevel(ll)
	}

	err = ctx.CreateExternalResultDir()
	if err != nil {
		l.Errorf("error creating results directory: %+v", err)
		os.Exit(1)
	}

	svc := github.NewSearchService(github.Connect(ctx.Context, ctx.GitHub.AccessToken))
	res, resp := svc.Search(ctx, ctx.GitHub.SearchType, ctx.GitHub.Query, nil)
	ctx.Logger.Debugf("github search response: %+v", resp)
	ctx.Logger.Debugf("github search result: %+v", res)

	ch := github.CloneFromResult(ctx, svc.Client, res)
	external.Run(ctx, ch)
}
