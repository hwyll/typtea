// internal/config/config.go
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Keybindings struct {
	Quit     []string `yaml:"quit"`
	Restart  []string `yaml:"restart"`
	NextLine []string `yaml:"next_line"`
}

type Config struct {
	Keybindings Keybindings `yaml:"keybindings"`
}

func DefaultConfig() *Config {
	return &Config{
		Keybindings: Keybindings{
			Quit:     []string{"ctrl+c", "esc"},
			Restart:  []string{"enter"},
			NextLine: []string{"enter"},
		},
	}
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "typtea", "config.yaml")
}

func Load() (*Config, error) {
	path := GetConfigPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig(), err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}

	return &cfg, nil
}
