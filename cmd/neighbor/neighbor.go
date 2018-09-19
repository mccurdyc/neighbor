package main

import (
	"context"
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/mccurdyc/neighbor/pkg/config"
	"github.com/mccurdyc/neighbor/pkg/github"
)

func main() {
	fp := flag.String("filepath", "config.json", "absolute filepath to config [default: \"$(PWD)/config.json\"].")

	cfg := config.New(*fp)
	cfg.Parse()

	ctx := context.Background()

	svc := github.NewSearchService(github.Connect(ctx, cfg.Contents.AccessToken))
	res, resp := svc.Search(ctx, cfg.Contents.SearchType, cfg.Contents.Query, nil)
	log.Infof("github search response: %+v", resp)

	nctx := github.CloneFromResult(ctx, svc.Client, res)

	// WIP; i dont like this because there is information leakage about the type of the response (map[string]string)
	// we will need to do something similar to this for cleanup of the temp directories
	for k, v := range nctx.Value(github.ClonedRepositoriesCtxKey{}).(map[string]string) {
		log.Debugf("K: %s, V: %s\n", k, v)
	}
}
