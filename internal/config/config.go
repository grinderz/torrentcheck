package config

import (
	"fmt"
	"github.com/grinderz/grgo/osutils"
	"github.com/grinderz/grgo/yamlutils"
)

type Config struct {
	SessionDir      string   `yaml:"sessionDir"`
	RootTorrentsDir string   `yaml:"rootTorrentsDir"`
	TorrentsSubdirs []string `yaml:"torrentsSubdirs,flow"`
}

func Load(configPath string) (*Config, error) {
	isExists, err := osutils.IsExists(configPath)
	if err != nil {
		return nil, fmt.Errorf("config load: %v", err)
	}
	if !isExists {
		return nil, fmt.Errorf("config load: error file not exists [%s]", configPath)
	}

	c := &Config{}
	if err = yamlutils.Load(configPath, c); err != nil {
		return nil, fmt.Errorf("config load: %v", err)
	}

	if err = c.validate(); err != nil {
		return nil, fmt.Errorf("config load: %v", err)
	}

	return c, nil
}

func (c *Config) Save(configPath string) error {
	return yamlutils.Save(configPath, c)
}

func (c *Config) validate() error {
	isExists, err := osutils.IsExists(c.SessionDir)
	if err != nil {
		return fmt.Errorf("validate session dir [%s]: %v", c.SessionDir, err)
	}
	if !isExists {
		return fmt.Errorf("validate session dir [%s]: error not exists", c.SessionDir)
	}
	isExists, err = osutils.IsExists(c.RootTorrentsDir)
	if err != nil {
		return fmt.Errorf("validate root torrents dir [%s]: %v", c.RootTorrentsDir, err)
	}
	if !isExists {
		return fmt.Errorf("validate root torrents dir [%s]: error not exists", c.RootTorrentsDir)
	}
	if len(c.TorrentsSubdirs) == 0 {
		return fmt.Errorf("validate torrents subdirs: error empty")
	}
	return nil
}
