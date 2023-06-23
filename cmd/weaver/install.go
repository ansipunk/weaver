package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"git.sr.ht/~ansipunk/weaver/pkg/cfg"
)

// Install is a CLI command that installs mods based on the configuration.
func Install(cCtx *cli.Context) error {
	// Read the configuration file
	config, err := cfg.ReadConfig(ConfigFileName)
	if err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	// Install mods using the configuration
	if err := InstallMods(config.Mods, config.Loader, config.GameVersion); err != nil {
		return fmt.Errorf("failed to install mods: %w", err)
	}

	return nil
}
