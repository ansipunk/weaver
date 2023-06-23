package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"git.sr.ht/~ansipunk/weaver/pkg/cfg"
)

func Add(cCtx *cli.Context) error {
	// Read configuration
	mods := cCtx.Args().Slice()

	if len(mods) == 0 {
		fmt.Println("No mods to add or install.")
	}

	config, readErr := cfg.ReadConfig(ConfigFileName)
	if readErr != nil {
		return fmt.Errorf("Failed to read configuration: %w", readErr)
	}

	installErr := InstallMods(mods, config.Loader, config.GameVersion)
	if installErr != nil {
		return fmt.Errorf("Failed to install new mods: %w", installErr)
	}

	config.Mods = append(config.Mods, mods...)

	writeErr := config.Write(ConfigFileName)
	if writeErr != nil {
		return fmt.Errorf("Failed to write to configuration: %w", writeErr)
	}

	return nil
}
