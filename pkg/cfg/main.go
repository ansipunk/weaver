package cfg

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Loader      string   `toml:"loader"`
	GameVersion string   `toml:"game_version"`
	Mods        []string `toml:"mods"`
}

func ParseConfig(config []byte) (Config, error) {
	var result Config
	if _, err := toml.Decode(string(config), &result); err != nil {
		return Config{}, err
	}
	return result, nil
}

func ReadConfig(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	return ParseConfig(content)
}
