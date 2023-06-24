package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// isDir checks if the given path is a directory.
func isDir(path string) (bool, error) {
	fileStat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileStat.IsDir(), nil
}

// getFileHash calculates the SHA1 hash of a file.
func getFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// EnsureDir ensures that the directory at the specified path exists.
// If the directory doesn't exist, it creates it with the default permissions.
func EnsureDir(path string) error {
	pathIsDir, err := isDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	} else if !pathIsDir {
		return errors.New(path + ": is not a directory")
	}

	return nil
}

// fileExists checks if the file exists at the specified path.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ShouldDownload checks if the file at the specified path should be downloaded based on its hash.
// It compares the hash of the existing file (if it exists) with the provided hash.
func ShouldDownload(path string, hash string) (bool, error) {
	if !fileExists(path) {
		return true, nil
	}

	fileHash, err := getFileHash(path)
	if err != nil {
		return false, err
	}

	return hash != fileHash, nil
}

// DeleteFile deletes the file at the specified path.
func DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// SaveFile saves the contents from the provided io.ReadCloser to the specified file path.
// If a progressBar is provided, it writes the contents with progress tracking.
func SaveFile(contents io.ReadCloser, path string, progressBar *progressbar.ProgressBar) error {
	defer contents.Close()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var writer io.Writer = file
	if progressBar != nil {
		writer = io.MultiWriter(file, progressBar)
	}

	_, err = io.Copy(writer, contents)
	return err
}

// convertToMap converts a string slice to a map with string keys and boolean values.
func convertToMap(strings []string) map[string]bool {
	result := make(map[string]bool, len(strings))
	for _, str := range strings {
		result[str] = true
	}
	return result
}

// RemoveOldFiles removes old files in the directory that are not in the list of required files.
func RemoveOldFiles(requiredFiles []string, directory string) error {
	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	requiredFilesMap := convertToMap(requiredFiles)

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() || !strings.HasSuffix(dirEntry.Name(), ".jar") || requiredFilesMap[dirEntry.Name()] {
			continue
		}

		if err := DeleteFile(filepath.Join(directory, dirEntry.Name())); err != nil {
			return err
		}
	}

	return nil
}
