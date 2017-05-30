package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config is the root cfg object
type Config struct {
	File             string             `yaml:",omitempty"`
	WorkingDirectory string             `yaml:"workdir,omitempty"`
	Excludes         []string           `yaml:"excludes"`
	Pipeline         []string           `yaml:"pipeline"`
	Growl            bool               `yaml:"growl,omitempty"`
	Server           string             `yaml:"server,omitempty"`
	Sidecars         map[string]Sidecar `yaml:"sidecars,omitempty"`
}

// Sidecar defines a background ("sidecar") service that is kept running
type Sidecar struct {
	Image       string            `yaml:"image"`
	Environment map[string]string `yaml:"env,omitempty"`
	Script      string            `yaml:"script,omitempty"`
	Volumes     map[string]string `yaml:"volumes,omitempty"`
}

var (
	cfg  *Config
	File string
)

// Load bootstraps the cfg object with the cfg file.
func Load() *Config {
	if cfg != nil {
		return cfg
	}
	configFile, err := filepath.Abs(File)
	if err != nil {
		panic(err)
	}
	cfg, err := ParseConfig(configFile)
	if err != nil {
		panic(err)
	}
	return cfg
}

// Returns the loaded cfg file name.
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

// ParseConfig reads and parses a cfg file.
func ParseConfig(cfgFile string) (*Config, error) {
	source, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(source, &cfg); err != nil {
		return nil, err
	}
	cfg.file = cfgFile
	return cfg, nil
}

func SerializeConfig(cfg *Config) ([]byte, error) {
	return yaml.Marshal(cfg)
}
