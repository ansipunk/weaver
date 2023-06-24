package main

import (
	"fmt"
	"sync"

	"github.com/schollz/progressbar/v3"

	"git.sr.ht/~ansipunk/weaver/pkg/modrinth"
	"git.sr.ht/~ansipunk/weaver/pkg/types"
)

// processErrors is a helper function that processes errors from the error channel.
// If an error is received, it returns the error. Otherwise, it returns nil.
func processErrors(errCh <-chan error) error {
	select {
	case err := <-errCh:
		// Handle the received error
		return err
	default:
		// No error received, continue execution
		return nil
	}
}

// Counter is a thread-safe counter that keeps track of a value.
type Counter struct {
	value int
	mutex sync.Mutex
}

// Increment increments the counter value by 1.
func (c *Counter) Increment() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value++
}

// Value returns the current value of the counter.
func (c *Counter) Value() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.value
}

// InstallMods installs mods based on the provided mod list, loader, and game version.
func InstallMods(mods []string, loader string, gameVersion string) error {
	// Ensure mod directory exists
	if err := EnsureDir(ModDirectory); err != nil {
		return err
	}

	// Get versions to download
	versionsToDownload, err := modrinth.GetAllVersionsToDownload(mods, loader, gameVersion)
	if err != nil {
		return fmt.Errorf("Failed to get versions to download: %w", err)
	}

	if len(versionsToDownload) == 0 {
		return nil
	}

	versionSlugs := make([]string, len(versionsToDownload))
	filenames := make([]string, len(versionsToDownload))

	var (
		downloadedCounter Counter
		skippedCounter    Counter
	)

	errCh := make(chan error) // Channel to receive errors from goroutines

	defer close(errCh) // Close the error channel when InstallMods returns

	var wg sync.WaitGroup // WaitGroup for synchronization

	fmt.Println("Mods to install:")
	fmt.Println("================")

	for i, version := range versionsToDownload {
		primaryFile := modrinth.GetVersionPrimaryFile(&version)
		filename := version.Slug + ".jar"
		versionSlugs[i] = version.Slug
		filenames[i] = filename

		shouldDownload, err := ShouldDownload(ModDirectory+filename, primaryFile.Hashes.Sha1)
		if err != nil {
			return fmt.Errorf("Failed to check if download is needed: %w", err)
		}

		if shouldDownload {
			wg.Add(1) // Increment WaitGroup counter for each goroutine

			go func(version types.Version, filename string) {
				defer wg.Done() // Signal the WaitGroup that the goroutine is done

				reader, err := modrinth.DownloadFile(primaryFile)
				if err != nil {
					errCh <- fmt.Errorf("Failed to download file for version %s: %v", version.Slug, err)
					return
				}
				defer reader.Close()

				if err := DeleteFile(ModDirectory + filename); err != nil {
					errCh <- fmt.Errorf("Failed to delete existing file for version %s: %v", version.Slug, err)
					return
				}

				pb := progressbar.DefaultBytes(primaryFile.Size, version.Slug)
				if err := SaveFile(reader, ModDirectory+filename, pb); err != nil {
					errCh <- fmt.Errorf("Failed to save file for version %s: %v", version.Slug, err)
					return
				}

				downloadedCounter.Increment()
			}(version, filename)
		} else {
			fmt.Printf("Skipped: %s (already up to date)\n", version.Slug)
			skippedCounter.Increment()
		}
	}

	wg.Wait() // Wait for all goroutines to finish

	if err := processErrors(errCh); err != nil {
		// Handle the returned error
		return err
	}

	fmt.Println("================")
	fmt.Println("Downloaded:", downloadedCounter.Value())
	fmt.Println("Skipped:", skippedCounter.Value())

	if err := RemoveOldFiles(filenames, ModDirectory); err != nil {
		return fmt.Errorf("Failed to remove old files: %w", err)
	}

	return nil
}
