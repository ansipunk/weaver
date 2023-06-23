package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"git.sr.ht/~ansipunk/weaver/pkg/cfg"
)

func Add(cCtx *cli.Context) error {
	// Read configuration
	mods := cCtx.Args().Slice()

	config, readErr := cfg.ReadConfig("weaver.toml")
	if readErr != nil {
		return fmt.Errorf("Failed to read configuration: %w", readErr)
	}

	installErr := InstallMods(mods, config.Loader, config.GameVersion)
	if installErr != nil {
		return fmt.Errorf("Failed to install new mods: %w", installErr)
	}

	config.Mods = append(config.Mods, mods...)
	return config.Write("weaver.toml")
}
