package modrinth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
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

func GetLatestVersion(projectId string, loader string, gameVersion string) (Version, error) {
	loaders := "[\"" + loader + "\"]"
	gameVersions := "[\"" + gameVersion + "\"]"
	url := baseUrl + "/project/" + url.QueryEscape(projectId) + "/version" +
		"?loaders=" + url.QueryEscape(loaders) +
		"&game_versions=" + url.QueryEscape(gameVersions) +
		"&featured="

	body, getErr := makeRequest(url + "true")

	if getErr != nil {
		return Version{}, getErr
	}

	var versions []Version
	jsonErr := json.Unmarshal(body, &versions)

	if jsonErr != nil {
		return Version{}, jsonErr
	}

	if len(versions) < 1 {
		body, getErr = makeRequest(url + "false")

		if getErr != nil {
			return Version{}, getErr
		}

		jsonErr = json.Unmarshal(body, &versions)

		if jsonErr != nil {
			return Version{}, jsonErr
		}

		if len(versions) < 1 {
			err := "no versions available"
			return Version{}, errors.New(err)
		}
	}

	return versions[0], nil
}
