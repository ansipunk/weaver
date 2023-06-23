package modrinth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"
)

// GetLatestVersion retrieves the latest version of a project for the specified loader and game version.
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

// GetSpecificVersion retrieves a specific version of a project.
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

// GetAllVersionsToDownload retrieves all versions to download for the given mod names, loader, and game version.
func GetAllVersionsToDownload(modNames []string, loader, gameVersion string) ([]Version, error) {
	versionsToDownload := []Version{}
	var wg sync.WaitGroup
	var firstError error
	errorMutex := sync.Mutex{}

	modLock := sync.Mutex{}
	dependencyLock := sync.Mutex{}

	for _, modName := range modNames {
		wg.Add(1)
		go func(modName string) {
			defer wg.Done()
			version, err := GetLatestVersion(modName, loader, gameVersion)
			if err != nil {
				errorMutex.Lock()
				if firstError == nil {
					firstError = fmt.Errorf("failed to get latest version for mod '%s': %v", modName, err)
				}
				errorMutex.Unlock()
				return
			}

			modLock.Lock()
			versionsToDownload = append(versionsToDownload, version)
			modLock.Unlock()
		}(modName)
	}

	wg.Wait()

	if firstError != nil {
		return nil, firstError
	}

	for _, version := range versionsToDownload {
		wg.Add(1)
		go func(version Version) {
			defer wg.Done()
			dependencies, err := version.GetDependencies()
			if err != nil {
				errorMutex.Lock()
				if firstError == nil {
					firstError = fmt.Errorf("failed to get dependencies for version '%s': %v", version.Slug, err)
				}
				errorMutex.Unlock()
				return
			}

			dependencyLock.Lock()
			versionsToDownload = append(versionsToDownload, dependencies...)
			dependencyLock.Unlock()
		}(version)
	}

	wg.Wait()

	if firstError != nil {
		return nil, firstError
	}

	return deduplicateVersions(versionsToDownload), nil
}

// GetProject retrieves project information for the specified project ID.
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
