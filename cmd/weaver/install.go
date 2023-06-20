package main

import (
	"fmt"
	"git.sr.ht/~ansipunk/weaver/pkg/cfg"
	"git.sr.ht/~ansipunk/weaver/pkg/fs"
	"git.sr.ht/~ansipunk/weaver/pkg/modrinth"
	"github.com/urfave/cli/v2"
)

func Install(cCtx *cli.Context) error {
	config, configErr := cfg.ReadConfig("weaver.toml")

	if configErr != nil {
		return configErr
	}

	if ensureDirErr := fs.EnsureDir(modDirectory); ensureDirErr != nil {
		return ensureDirErr
	}

	fmt.Println(config.Loader)
	fmt.Println(config.GameVersion)

	for _, modName := range config.Mods {
		version, versionErr := modrinth.GetLatestVersion(modName, config.Loader, config.GameVersion)

		if versionErr != nil {
			return versionErr
		}

		primaryFile := version.GetPrimaryFile()

		shouldDownload, shouldErr := fs.ShouldDownload(
			modDirectory + modName + ".jar", primaryFile.Hashes.Sha1)

		if shouldErr != nil {
			return shouldErr
		}

		if shouldDownload {
			reader, readErr := primaryFile.Download()

			if readErr != nil {
				return readErr
			}

			defer reader.Close()

			if deleteErr := fs.DeleteFile(modDirectory + modName + ".jar"); deleteErr != nil {
				return deleteErr
			}

			if saveErr := fs.SaveFile(reader, modDirectory + modName + ".jar"); saveErr != nil {
				return saveErr
			}
		}
	}

	return nil
}