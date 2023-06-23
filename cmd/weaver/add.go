package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"git.sr.ht/~ansipunk/weaver/pkg/cfg"
)

func Add(cCtx *cli.Context) error {
	// Read configuration
	mods := cCtx.Args().Slice()

	config, readErr := cfg.ReadConfig("weaver.toml")
	if readErr != nil {
		return fmt.Errorf("failed to read configuration: %w", readErr)
	}
	config.Mods = append(config.Mods, mods...)

	file, openErr := os.Create("weaver.toml")
	if openErr != nil {
		return fmt.Errorf("failed to open configuration: %w", openErr)
	}
	defer file.Close()

	if writeErr := config.Dump(file); writeErr != nil {
		return fmt.Errorf("failed to write configuration: %w", openErr)
	}

	return nil
}
