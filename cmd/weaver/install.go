package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"git.sr.ht/~ansipunk/weaver/pkg/cfg"
)

func Install(cCtx *cli.Context) error {
	config, err := cfg.ReadConfig(ConfigFileName)
	if err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	if err := InstallMods(config.Mods, config.Loader, config.GameVersion); err != nil {
		return fmt.Errorf("failed to install mods: %w", err)
	}

	return nil
}
