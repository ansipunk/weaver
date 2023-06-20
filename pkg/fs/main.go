package fs

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

func isDir(path string) (bool, error) {
	file, openErr := os.Open(path)

	if openErr != nil {
		return false, openErr
	}

	defer file.Close()
	fileStat, statErr := file.Stat()

	if statErr != nil {
		return false, statErr
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

	return hex.EncodeToString(hash.Sum(nil)[:20]), nil
}

func EnsureDir(path string) error {
	pathIsDir, isDirErr := isDir(path)

	if isDirErr == nil && !pathIsDir {
		return errors.New(path + ": is not a directory")
	}

	if isDirErr != nil {
		mkdirErr := os.Mkdir(path, 0755)

		if mkdirErr != nil {
			return mkdirErr
		}
	}

	return nil
}

func fileExists(path string) error {
	_, isDirErr := isDir(path)
	return isDirErr
}

func ShouldDownload(path string, hash string) (bool, error) {
	fileExistsErr := fileExists(path)

	if fileExistsErr != nil {
		return true, nil
	}

	fileHash, getFileHashError := getFileHash(path)

	if getFileHashError != nil {
		return false, getFileHashError
	}

	return hash != fileHash, nil
}

func DeleteFile(path string) error {
	return os.Remove(path)
}

func SaveFile(contents io.ReadCloser, path string) error {
	defer contents.Close()

	file, createErr := os.Create(path)

	if createErr != nil {
		return createErr
	}

	defer file.Close()

	_, err := io.Copy(file, contents)
	return err
}
