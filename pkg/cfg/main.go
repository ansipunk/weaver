package cfg

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Loader      string   `toml:"loader"`
	GameVersion string   `toml:"game_version"`
	Mods        []string `toml:"mods"`
}

func (c *Config) Dump(file *os.File) error {
	encoder := toml.NewEncoder(file)
	return encoder.Encode(c)
}

func (c *Config) Write(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to open file %v: %w", path, err)
	}
	defer file.Close()

	if err := c.Dump(file); err != nil {
		return fmt.Errorf("failed to dump config to file %v: %w", path, err)
	}

	return nil
}

func ParseConfig(config []byte) (Config, error) {
	var result Config
	if _, err := toml.Decode(string(config), &result); err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}
	return result, nil
}

func ReadConfig(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read file %v: %w", path, err)
	}
	return ParseConfig(content)
}
