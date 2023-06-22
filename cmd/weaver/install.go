package main

import (
	"fmt"
	"sync"

	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"

	"git.sr.ht/~ansipunk/weaver/pkg/cfg"
	"git.sr.ht/~ansipunk/weaver/pkg/fs"
	"git.sr.ht/~ansipunk/weaver/pkg/modrinth"
)

func Install(cCtx *cli.Context) error {
	// Read configuration
	config, err := cfg.ReadConfig("weaver.toml")
	if err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	// Ensure mod directory exists
	if err := fs.EnsureDir(modDirectory); err != nil {
		return fmt.Errorf("failed to ensure mod directory exists: %w", err)
	}

	// Get versions to download
	versionsToDownload, err := modrinth.GetAllVersionsToDownload(&config.Mods, &config.Loader, &config.GameVersion)
	if err != nil {
		return fmt.Errorf("failed to get versions to download: %w", err)
	}

	versionSlugs := make([]string, len(versionsToDownload))
	filenames := make([]string, len(versionsToDownload))

	var (
		downloaded uint16
		skipped    uint16
		wg         sync.WaitGroup // WaitGroup for synchronization
	)

	fmt.Println("Mods to install:")
	fmt.Println("================")

	for i, version := range versionsToDownload {
		primaryFile := version.GetPrimaryFile()
		filename := version.Slug + ".jar"
		versionSlugs[i] = version.Slug
		filenames[i] = filename

		shouldDownload, err := fs.ShouldDownload(modDirectory+filename, primaryFile.Hashes.Sha1)
		if err != nil {
			return fmt.Errorf("failed to check if download is needed: %w", err)
		}

		if shouldDownload {
			wg.Add(1) // Increment WaitGroup counter for each goroutine

			go func(version modrinth.Version) {
				defer wg.Done() // Signal the WaitGroup that the goroutine is done

				reader, err := primaryFile.Download()
				if err != nil {
					fmt.Printf("failed to download file for version %s: %v\n", version.Slug, err)
					return
				}
				defer reader.Close()

				if err := fs.DeleteFile(modDirectory + filename); err != nil {
					fmt.Printf("failed to delete existing file for version %s: %v\n", version.Slug, err)
					return
				}

				pb := progressbar.DefaultBytes(primaryFile.Size, version.Slug)
				if err := fs.SaveFile(reader, modDirectory+filename, pb); err != nil {
					fmt.Printf("failed to save file for version %s: %v\n", version.Slug, err)
					return
				}

				downloaded++
			}(version)
		} else {
			fmt.Printf("Skipped: %s (already up to date)\n", version.Slug)
			skipped++
		}
	}

	wg.Wait() // Wait for all goroutines to finish

	fmt.Println("================")
	fmt.Println("Downloaded:", downloaded)
	fmt.Println("Skipped:", skipped)

	if err := fs.RemoveOldFiles(filenames, modDirectory); err != nil {
		return fmt.Errorf("failed to remove old files: %w", err)
	}

	return nil
}
