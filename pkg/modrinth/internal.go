package modrinth

import (
	"io"
	"net/http"
)

const baseURL = "https://api.modrinth.com/v2"

func makeRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

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
