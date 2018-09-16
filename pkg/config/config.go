package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

// Contents contains the contents of the parsed config file.
type Contents struct {
	AccessToken string          `json:"access_token"`
	SearchType  string          `json:"search_type"`
	Query       json.RawMessage `json:"query"`
	// Code        struct {
	// 	Query github.CodeQuery `json:"query"`
	// } `json:"code"`
	// Repository struct {
	// 	Query github.RepositoryQuery `json:"query"`
	// } `json:"repository"`
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
func (cfg *Config) Parse() error {
	f, err := os.Open(cfg.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var c interface{}
	if err := parse(f, c); err != nil {
		return err
	}

	cfg.Contents = c.(*Contents)

	return nil
}

func parse(f io.Reader, d interface{}) error {
	b, err := ioutil.ReadAll(f)

	if err = json.Unmarshal(b, d); err != nil {
		return err
	}

	return nil
}
