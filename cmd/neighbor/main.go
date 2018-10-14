package main

import (
	// stdlib

	"context"
	"flag"
	"fmt"
	"os"

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

	// create a context object that will be used for the life of the program and passed around
	ctx := &neighbor.Ctx{
		Config:  cfg,
		Context: context.Background(),
		GitHub: neighbor.GitHubDetails{
			AccessToken: cfg.Contents.AccessToken,
		},
		Logger:        l,
		NeighborDir:   wd,
		ProjectDirMap: make(map[string]string),
		ExtResultDir:  fmt.Sprintf("%s/%s", wd, "ext-results"),
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

	// populates the context's ProjectDirMap with cloned projects and where they were cloned
	github.CloneFromResult(ctx, svc.Client, res)

	neighbor.SetTestCmd(ctx)

	external.RunTests(ctx)
}
