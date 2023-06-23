package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"git.sr.ht/~ansipunk/weaver/pkg/cfg"
)

// Add is a CLI command that adds mods to the configuration and installs them.
func Add(cCtx *cli.Context) error {
	// Read configuration
	mods := cCtx.Args().Slice()

	if len(mods) == 0 {
		fmt.Println("No mods to add or install.")
	}

	// Read the configuration file
	config, readErr := cfg.ReadConfig(ConfigFileName)
	if readErr != nil {
		return fmt.Errorf("Failed to read configuration: %w", readErr)
	}

	// Install mods
	installErr := InstallMods(mods, config.Loader, config.GameVersion)
	if installErr != nil {
		return fmt.Errorf("Failed to install new mods: %w", installErr)
	}

	// Append the new mods to the configuration
	config.Mods = append(config.Mods, mods...)

	// Write the updated configuration file
	writeErr := config.Write(ConfigFileName)
	if writeErr != nil {
		return fmt.Errorf("Failed to write to configuration: %w", writeErr)
	}

	return nil
}
