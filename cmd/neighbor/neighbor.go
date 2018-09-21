package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/mccurdyc/neighbor/pkg/config"
	"github.com/mccurdyc/neighbor/pkg/github"
	"github.com/mccurdyc/neighbor/pkg/neighbor"
)

func main() {
	fp := flag.String("filepath", "config.json", "absolute filepath to config [default: \"$(PWD)/config.json\"].")

	cfg := config.New(*fp)
	cfg.Parse()
	ctx := neighbor.NewCtx(cfg)

	svc := github.NewSearchService(github.Connect(context.Background(), cfg.Contents.AccessToken))
	res, resp := svc.Search(context.Background(), cfg.Contents.SearchType, cfg.Contents.Query, nil)
	log.Infof("github search response: %+v", resp)

	github.CloneFromResult(ctx, svc.Client, res)

	// WIP; i dont like this because there is information leakage about the type of the response (map[string]string)
	// we will need to do something similar to this for cleanup of the temp directories
	for name, dir := range ctx.ProjectDirMap {
		log.Debugf("K: %s, V: %s\n", name, dir)

		err := os.Chdir(dir)
		if err != nil {
			log.Error(err)
			continue
		}

		cmd := exec.Command("make", "test")
		var out bytes.Buffer
		cmd.Stdout = &out
		fmt.Println(out.String())

		err = cmd.Run()
		if err != nil {
			log.Error(err)
			continue
		}
	}
}
