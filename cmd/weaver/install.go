package main

import (
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

	versionsToDownload, verErr := modrinth.GetAllVersionsToDownload(&config.Mods, &config.Loader, &config.GameVersion)
	filenames := []string{}

	if verErr != nil {
		return verErr
	}

	for _, version := range versionsToDownload {
		primaryFile := version.GetPrimaryFile()
		filename := version.Slug + ".jar"

		shouldDownload, shouldErr := fs.ShouldDownload(
			modDirectory+version.Slug+".jar", primaryFile.Hashes.Sha1)

		if shouldErr != nil {
			return shouldErr
		}

		if shouldDownload {
			reader, readErr := primaryFile.Download()

			if readErr != nil {
				return readErr
			}

			defer reader.Close()

			if deleteErr := fs.DeleteFile(modDirectory + filename); deleteErr != nil {
				return deleteErr
			}

			if saveErr := fs.SaveFile(reader, modDirectory+filename); saveErr != nil {
				return saveErr
			}
		}

		filenames = append(filenames, filename)
	}

	return fs.RemoveOldFiles(filenames, modDirectory)
}
