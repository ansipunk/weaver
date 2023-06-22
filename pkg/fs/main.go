package fs

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

func isDir(path string) (bool, error) {
	fileStat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileStat.IsDir(), nil
}

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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

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

func DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

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

func RemoveOldFiles(requiredFiles []string, directory string) error {
	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() || !strings.HasSuffix(dirEntry.Name(), ".jar") || contains(requiredFiles, dirEntry.Name()) {
			continue
		}

		if err := DeleteFile(filepath.Join(directory, dirEntry.Name())); err != nil {
			return err
		}
	}

	return nil
}

func contains(strings []string, s string) bool {
	for _, str := range strings {
		if str == s {
			return true
		}
	}
	return false
}
