package config

import (
	// stdlib
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/prometheus/common/log"
	// external
	// internal
)

// Contents contains the contents of the parsed config file.
type Contents struct {
	AccessToken string          `json:"access_token"`
	SSHKey      string          `json:"ssh_key"`
	SearchType  string          `json:"search_type"`
	Query       json.RawMessage `json:"query"` // RawMessage allows us to handle parsing this bit later

	TestCmdStr string `json:"external_test_command"`
	TestCmd    *exec.Cmd
}

// Config specifies information about the config file used for performing the experiment.
type Config struct {
	FilePath string
	Contents *Contents
}

// New is a constructor that returns a pointer to a Config object.
func New(fp string) *Config {
	return &Config{
		FilePath: fp,
	}
}

// Parse opens a file and calls a private `parse` method.
func (cfg *Config) Parse() {
	f, err := os.Open(cfg.FilePath)
	if err != nil {
		log.Errorf("error opening config file %+v", err)
		return
	}
	defer f.Close()

	c := &Contents{}
	if err := parse(f, c); err != nil {
		log.Errorf("error parsing config file %+v", err)
		return
	}

	cfg.Contents = c

	// set external test command
	if len(cfg.Contents.TestCmdStr) != 0 {
		s := strings.Split(cfg.Contents.TestCmdStr, " ")
		cfg.Contents.TestCmd = exec.Command(s[0], s[1:]...)
	}

	return
}

func parse(f io.Reader, d interface{}) error {
	b, err := ioutil.ReadAll(f)

	if err = json.Unmarshal(b, d); err != nil {
		return err
	}

	return nil
}
