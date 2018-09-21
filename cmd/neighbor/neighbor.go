package main

import (
	// stdlib
	"context"
	"flag"

	// external

	// internal
	"github.com/mccurdyc/neighbor/pkg/config"
	"github.com/mccurdyc/neighbor/pkg/github"
	"github.com/mccurdyc/neighbor/pkg/neighbor"
)

func main() {
	fp := flag.String("filepath", "config.json", "absolute filepath to config [default: \"$(PWD)/config.json\"].")

	cfg := config.New(*fp)
	cfg.Parse()
	ctx := neighbor.NewCtx(context.Background(), cfg)

	svc := github.NewSearchService(github.Connect(ctx, cfg.Contents.AccessToken))
	res, resp := svc.Search(ctx, cfg.Contents.SearchType, cfg.Contents.Query, nil)
	ctx.Logger.Infof("github search response: %+v", resp)

	// populates the context's ProjectDirMap
	github.CloneFromResult(ctx, svc.Client, res)
}
