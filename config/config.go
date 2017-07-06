package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config is the root config object
type Config struct {
	file             string   
	WorkingDirectory string   `yaml:"workdir,omitempty"`
	Excludes         []string `yaml:"excludes"`
	Pipeline         []string `yaml:"pipeline"`
	Growl            bool     `yaml:"growl,omitempty"`
	Server           string   `yaml:"server,omitempty"`
}

var config Config

// Setup bootstraps the config object with the config file.
func Setup(cfgFileRel string) *Config {
	configFile, err := filepath.Abs(cfgFileRel)
	if err != nil {
		panic(err)
	}
	cfg, err := ParseConfig(configFile)
	if err != nil {
		panic(err)
	}
	return cfg
}

// Returns the loaded config file name.
func (c *Config) GetConfigFile() string {
	return c.file
}

// Reloads this configuration. It will panic if an error occurs.
func (c *Config) Reload() {
	var err error
	_, err = ParseConfig(c.file)
	if err != nil {
		panic(err)
	}
}

// ParseConfig reads and parses a config file.
func ParseConfig(cfgFile string) (*Config, error) {
	source, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(source, &config); err != nil {
		return nil, err
	}
	config.file = cfgFile
	return &config, nil
}

func SerializeConfig(cfg *Config) ([]byte, error) {
	return yaml.Marshal(cfg)
}
