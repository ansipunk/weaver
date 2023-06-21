package main

import (
	"fmt"
	"git.sr.ht/~ansipunk/weaver/pkg/cfg"
	"git.sr.ht/~ansipunk/weaver/pkg/fs"
	"git.sr.ht/~ansipunk/weaver/pkg/modrinth"
	"github.com/schollz/progressbar/v3"
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
	versionSlugs := []string{}

	for _, version := range versionsToDownload {
		versionSlugs = append(versionSlugs, version.Slug)
	}

	filenames := []string{}

	if verErr != nil {
		return verErr
	}

	var downloaded uint16
	var skipped uint16

	fmt.Println("Mods to install:", versionSlugs)
	fmt.Println("================")

	for _, version := range versionsToDownload {
		primaryFile := version.GetPrimaryFile()
		filename := version.Slug + ".jar"

		shouldDownload, shouldErr := fs.ShouldDownload(
			modDirectory+filename, primaryFile.Hashes.Sha1)

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

			pb := progressbar.DefaultBytes(primaryFile.Size, version.Slug)

			if saveErr := fs.SaveFile(reader, modDirectory+filename, pb); saveErr != nil {
				return saveErr
			}

			downloaded++
		} else {
			fmt.Println(version.Slug, "is already up to date")
			skipped++
		}

		filenames = append(filenames, filename)
	}

	fmt.Println("================")
	fmt.Println("Downloaded:", downloaded)
	fmt.Println("Skipped:", skipped)

	return fs.RemoveOldFiles(filenames, modDirectory)
}
