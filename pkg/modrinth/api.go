package modrinth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"git.sr.ht/~ansipunk/weaver/pkg/types"
)

// GetLatestVersion retrieves the latest version of a project for the specified loader and game version.
func GetLatestVersion(projectSlug, loader, gameVersion string) (types.Version, error) {
	loaders := "[" + `"` + loader + `"` + "]"
	gameVersions := "[" + `"` + gameVersion + `"` + "]"
	requestURL := baseURL + "/project/" + url.QueryEscape(projectSlug) + "/version" +
		"?loaders=" + url.QueryEscape(loaders) +
		"&game_versions=" + url.QueryEscape(gameVersions) +
		"&featured="

	body, err := makeRequest(requestURL + "true")
	if err != nil {
		return types.Version{}, err
	}

	var versions []types.Version
	err = json.Unmarshal(body, &versions)
	if err != nil {
		return types.Version{}, err
	}

	if len(versions) < 1 {
		body, err = makeRequest(requestURL + "false")
		if err != nil {
			return types.Version{}, err
		}

		err = json.Unmarshal(body, &versions)
		if err != nil {
			return types.Version{}, err
		}

		if len(versions) < 1 {
			return types.Version{}, errors.New("no versions available")
		}
	}

	versions[0].Slug = projectSlug
	return versions[0], nil
}

// GetSpecificVersion retrieves a specific version of a project.
func GetSpecificVersion(versionId string) (types.Version, error) {
	var version types.Version
	requestURL := baseURL + "/version/" + url.QueryEscape(versionId)
	body, err := makeRequest(requestURL)
	if err != nil {
		return types.Version{}, err
	}

	err = json.Unmarshal(body, &version)
	if err != nil {
		return types.Version{}, err
	}

	err = SetVersionProjectSlug(&version)
	if err != nil {
		return types.Version{}, err
	}

	return version, nil
}

// GetAllVersionsToDownload retrieves all versions to download for the given mod names, loader, and game version.
func GetAllVersionsToDownload(modNames []string, loader, gameVersion string) ([]types.Version, error) {
	versionsToDownload := []types.Version{}
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
		go func(version types.Version) {
			defer wg.Done()
			dependencies, err := GetVersionDependencies(&version)
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
func GetProject(projectId string) (types.Project, error) {
	var project types.Project
	requestURL := baseURL + "/project/" + url.QueryEscape(projectId)
	body, err := makeRequest(requestURL)
	if err != nil {
		return types.Project{}, err
	}

	err = json.Unmarshal(body, &project)
	if err != nil {
		return types.Project{}, err
	}

	return project, nil
}

// GetVersionDependencies retrieves the dependencies of a Modrinth version.
func GetVersionDependencies(version *types.Version) ([]types.Version, error) {
	dependencyVersions := []types.Version{}

	for _, dependency := range version.Dependencies {
		if dependency.VersionID != "" {
			specificVersion, err := GetSpecificVersion(dependency.VersionID)
			if err != nil {
				return dependencyVersions, err
			}
			dependencyVersions = append(dependencyVersions, specificVersion)
		}
	}

	version.Dependencies = nil
	return dependencyVersions, nil
}

// GetVersionPrimaryFile returns the primary file associated with a Modrinth version.
func GetVersionPrimaryFile(version *types.Version) *types.File {
	if len(version.Files) == 1 {
		return &version.Files[0]
	}

	primaryIndex := 0

	for i, file := range version.Files {
		if file.Primary {
			primaryIndex = i
			break
		}
	}

	return &version.Files[primaryIndex]
}

// DownloadFile downloads the file associated with a Modrinth version.
func DownloadFile(file *types.File) (io.ReadCloser, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(file.URL)
	if err != nil {
		return nil, fmt.Errorf("Failed to download file: %w", err)
	}

	return resp.Body, nil
}

// SetVersionProjectSlug sets the project slug for a Modrinth version.
func SetVersionProjectSlug(version *types.Version) error {
	project, err := GetProject(version.ProjectID)
	if err != nil {
		return err
	}

	version.Slug = project.Slug
	return nil
}
