package modrinth

import (
	"encoding/json"
	"errors"
	"net/url"
)

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

func GetSpecificVersion(versionId string) (Version, error) {
	var version Version
	url := baseUrl + "/version/" + url.QueryEscape(versionId)
	body, getErr := makeRequest(url)

	if getErr != nil {
		return Version{}, getErr
	}

	jsonErr := json.Unmarshal(body, &version)

	if jsonErr != nil {
		return Version{}, jsonErr
	}

	return version, nil
}

func GetAllVersionsToDownload(modNames *[]string, loader, gameVersion *string) ([]Version, error) {
	versionsToDownload := []Version{}

	for _, modName := range *modNames {
		version, versionErr := GetLatestVersion(modName, *loader, *gameVersion)

		if versionErr != nil {
			return versionsToDownload, versionErr
		}

		versionsToDownload = append(versionsToDownload, version)
	}

	for _, version := range versionsToDownload {
		dependencies, depErr := version.GetDependencies()

		if depErr != nil {
			return versionsToDownload, depErr
		}

		versionsToDownload = append(versionsToDownload, dependencies...)
	}

	return deduplicateVersions(&versionsToDownload), nil
}
