package modrinth

import (
	"encoding/json"
	"errors"
	"net/url"
)

func GetLatestVersion(projectSlug, loader, gameVersion string) (Version, error) {
	loaders := "[" + `"` + loader + `"` + "]"
	gameVersions := "[" + `"` + gameVersion + `"` + "]"
	requestURL := baseURL + "/project/" + url.QueryEscape(projectSlug) + "/version" +
		"?loaders=" + url.QueryEscape(loaders) +
		"&game_versions=" + url.QueryEscape(gameVersions) +
		"&featured="

	body, err := makeRequest(requestURL + "true")
	if err != nil {
		return Version{}, err
	}

	var versions []Version
	err = json.Unmarshal(body, &versions)
	if err != nil {
		return Version{}, err
	}

	if len(versions) < 1 {
		body, err = makeRequest(requestURL + "false")
		if err != nil {
			return Version{}, err
		}

		err = json.Unmarshal(body, &versions)
		if err != nil {
			return Version{}, err
		}

		if len(versions) < 1 {
			return Version{}, errors.New("no versions available")
		}
	}

	versions[0].Slug = projectSlug
	return versions[0], nil
}

func GetSpecificVersion(versionId string) (Version, error) {
	var version Version
	requestURL := baseURL + "/version/" + url.QueryEscape(versionId)
	body, err := makeRequest(requestURL)
	if err != nil {
		return Version{}, err
	}

	err = json.Unmarshal(body, &version)
	if err != nil {
		return Version{}, err
	}

	err = version.SetProjectSlug()
	if err != nil {
		return Version{}, err
	}

	return version, nil
}

func GetAllVersionsToDownload(modNames *[]string, loader, gameVersion *string) ([]Version, error) {
	versionsToDownload := []Version{}

	for _, modName := range *modNames {
		version, err := GetLatestVersion(modName, *loader, *gameVersion)
		if err != nil {
			return nil, err
		}

		versionsToDownload = append(versionsToDownload, version)
	}

	for _, version := range versionsToDownload {
		dependencies, err := version.GetDependencies()
		if err != nil {
			return nil, err
		}

		versionsToDownload = append(versionsToDownload, dependencies...)
	}

	return deduplicateVersions(versionsToDownload), nil
}

func GetProject(projectId string) (Project, error) {
	var project Project
	requestURL := baseURL + "/project/" + url.QueryEscape(projectId)
	body, err := makeRequest(requestURL)
	if err != nil {
		return Project{}, err
	}

	err = json.Unmarshal(body, &project)
	if err != nil {
		return Project{}, err
	}

	return project, nil
}
