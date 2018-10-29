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
	log "github.com/sirupsen/logrus"

	// internal
	"github.com/mccurdyc/neighbor/pkg/config"
	"github.com/mccurdyc/neighbor/pkg/external"
	"github.com/mccurdyc/neighbor/pkg/github"
	"github.com/mccurdyc/neighbor/pkg/neighbor"
)

func main() {
	fp := flag.String("filepath", "config.json", "absolute filepath to config [default: \"$(PWD)/config.json\"].")

	cfg := config.New(*fp)
	cfg.Parse()

	l := log.New()

	wd, err := os.Getwd()
	if err != nil {
		l.Errorf("error getting current directory: %+v", err)
		os.Exit(1)
	}

	c, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// create a context object that will be used for the life of the program and passed around
	ctx := &neighbor.Ctx{
		Config:  cfg,
		Context: c,
		GitHub: neighbor.GitHubDetails{
			AccessToken: cfg.Contents.AccessToken,
		},
		Logger:       l,
		NeighborDir:  wd,
		ExtResultDir: fmt.Sprintf("%s/%s", wd, "_ext-results"), // go tools handle directories prepended with '_' differently; often they ignore those directories
	}

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

	svc := github.NewSearchService(github.Connect(ctx.Context, cfg.Contents.AccessToken))
	res, resp := svc.Search(ctx, cfg.Contents.SearchType, cfg.Contents.Query, nil)
	ctx.Logger.Debugf("github search response: %+v", resp)
	ctx.Logger.Debugf("github search result: %+v", res)

	neighbor.SetExternalCmd(ctx)
	ctx.Logger.Infof("external command to be run on each project: %s\n", ctx.ExternalCmd)

	ch := github.CloneFromResult(ctx, svc.Client, res)
	external.Run(ctx, ch)
}
