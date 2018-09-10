package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mccurdyc/neighbor/pkg/config"
	"github.com/mccurdyc/neighbor/pkg/github"
)

func main() {
	fp := flag.String("filepath", "config.yml", "absolute filepath to config [default: \"$(PWD)/config.yml\"].")

	cfg := config.New(*fp)
	err := cfg.Parse()
	if err != nil {
		fmt.Printf("error parsing config file at (%s): %+v\n", *fp, err)
		os.Exit(1)
	}

	ctx := context.Background()
	svc := github.NewSearchService(github.Connect(ctx, cfg.Contents.AccessToken))
	res, resp, err := svc.Search(ctx, cfg.Contents.SearchType, cfg.Contents.Query, nil)

	fmt.Println(res)
	fmt.Println(resp)
	fmt.Println(err)
}
