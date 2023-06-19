package cfg

import (
	"github.com/BurntSushi/toml"
	"io"
	"os"
)

type Config struct {
	Loader      string   `toml:"loader"`
	GameVersion string   `toml:"game_version"`
	Mods        []string `toml:"mods"`
}

func ParseConfig(config []byte) (Config, error) {
	var result Config
	_, err := toml.Decode(string(config), &result)

	if err != nil {
		return Config{}, err
	}

	return result, nil
}

func ReadConfig(path string) (Config, error) {
	file, openErr := os.Open(path)

	if openErr != nil {
		return Config{}, openErr
	}

	defer file.Close()
	content, readErr := io.ReadAll(file)

	if readErr != nil {
		return Config{}, readErr
	}

	parsed, parseErr := ParseConfig(content)

	if parseErr != nil {
		return Config{}, parseErr
	}

	return parsed, nil
}
