package config

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

func Setup(c *cli.Context) (*Config, time.Duration) {
	configFile, err := filepath.Abs(c.GlobalString("config"))
	if err != nil {
		panic(err)
	}
	cfg, err := ParseConfig(configFile)
	if err != nil {
		panic(err)
	}
	delay, err := time.ParseDuration(c.GlobalString("delay"))
	if err != nil {
		panic(err)
	}
	return cfg, delay
}

func ParseConfig(cfgFile string) (*Config, error) {
	source, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(source, &cfg); err != nil {
		return nil, err
	}
	cfg.File = cfgFile
	return &cfg, nil
}
