package main

import (
	// stdlib
	"context"
	"flag"

	// external

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
	ctx := neighbor.NewCtx(context.Background(), cfg)

	svc := github.NewSearchService(github.Connect(ctx.Context, cfg.Contents.AccessToken))
	res, resp := svc.Search(ctx, cfg.Contents.SearchType, cfg.Contents.Query, nil)
	ctx.Logger.Infof("github search response: %+v", resp)
	ctx.Logger.Infof("github search result: %+v", res)

	// populates the context's ProjectDirMap with cloned projects and where they were cloned
	github.CloneFromResult(ctx, svc.Client, res)

	ctx.TestCmd = cfg.Contents.TestCmd
	external.RunTests(ctx)
}
