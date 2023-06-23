package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"git.sr.ht/~ansipunk/weaver/pkg/cfg"
)

// contains checks if a given string is present in a string slice.
func contains(arr []string, target string) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}

// Remove is a CLI command that removes specified mods from the configuration.
func Remove(cCtx *cli.Context) error {
	// Read configuration
	modsToRemove := cCtx.Args().Slice()

	if len(modsToRemove) == 0 {
		fmt.Println("No mods to remove.")
		return nil
	}

	// Read the configuration file
	config, readErr := cfg.ReadConfig(ConfigFileName)
	if readErr != nil {
		return fmt.Errorf("Failed to read configuration: %w", readErr)
	}

	reducedMods := []string{}
	for _, installedMod := range config.Mods {
		if contains(modsToRemove, installedMod) {
			fmt.Printf("Removing mod: %s\n", installedMod)
		} else {
			reducedMods = append(reducedMods, installedMod)
		}
	}

	for _, modToRemove := range modsToRemove {
		if !contains(config.Mods, modToRemove) {
			fmt.Printf("Mod not installed: %s\n", modToRemove)
		}
	}

	config.Mods = reducedMods

	// Write the updated configuration file
	writeErr := config.Write(ConfigFileName)
	if writeErr != nil {
		return fmt.Errorf("Failed to write to configuration: %w", writeErr)
	}

	return nil
}
