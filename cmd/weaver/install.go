package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"git.sr.ht/~ansipunk/weaver/pkg/cfg"
)

func Install(cCtx *cli.Context) error {
	config, err := cfg.ReadConfig("weaver.toml")
	if err != nil {
		return fmt.Errorf("Failed to read configuration: %w", err)
	}

	return InstallMods(config.Mods, config.Loader, config.GameVersion)
}
