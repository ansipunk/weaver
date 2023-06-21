package modrinth

import (
	"git.sr.ht/~ansipunk/weaver/pkg/utils"
	"io"
	"net/http"
)

const baseUrl string = "https://api.modrinth.com/v2"

func makeRequest(url string) ([]byte, error) {
	resp, getErr := http.Get(url)

	if getErr != nil {
		return []byte{}, getErr
	}

	body, readErr := io.ReadAll(resp.Body)

	if readErr != nil {
		return []byte{}, readErr
	}

	return body, nil
}

func deduplicateVersions(versions *[]Version) []Version {
	deduplicated := []Version{}
	modNames := []string{}

	for _, version := range *versions {
		if !utils.Contains(&modNames, version.ProjectId) {
			modNames = append(modNames, version.ProjectId)
			deduplicated = append(deduplicated, version)
		}
	}

	return deduplicated
}
