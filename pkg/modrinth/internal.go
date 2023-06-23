package modrinth

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.modrinth.com/v2"

// makeRequest makes an HTTP GET request to the specified URL and returns the response body.
func makeRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode != 200 {
		if statusCode == 404 {
			return nil, errors.New("Mod or version was not found")
		} else if statusCode == 500 {
			return nil, errors.New("Modrinth server error")
		} else {
			return nil, errors.New("Request failed with status " + fmt.Sprint(statusCode))
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// deduplicateVersions removes duplicate versions from the provided slice.
func deduplicateVersions(versions []Version) []Version {
	deduplicated := []Version{}
	modNames := make(map[string]bool)

	for _, version := range versions {
		if !modNames[version.ProjectID] {
			modNames[version.ProjectID] = true
			deduplicated = append(deduplicated, version)
		}
	}

	return deduplicated
}
