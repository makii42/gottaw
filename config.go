package main

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

// Config is the root config object
type Config struct {
	File             string
	WorkingDirectory string   `yaml:"workdir"`
	Excludes         []string `yaml:"excludes"`
	Pipeline         []string `yaml:"pipeline"`
	Growl            bool     `yaml:"growl"`
}

func setup(c *cli.Context) time.Duration {
	configFile, err := filepath.Abs(c.String("config"))
	if err != nil {
		panic(err)
	}
	err = parseConfig(configFile)
	if err != nil {
		panic(err)
	}
	delay, err := time.ParseDuration(c.String("delay"))
	if err != nil {
		panic(err)
	}
	return delay
}

func parseConfig(cfgFile string) error {
	source, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(source, &cfg); err != nil {
		return err
	}
	cfg.File = cfgFile
	return nil
}
