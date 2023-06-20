package modrinth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseUrl string = "https://api.modrinth.com/v2"

type Dependency struct {
	VersionId      string `json:"version_id,omitempty"`
	ProjectId      string `json:"project_id,omitempty"`
	FileName       string `json:"file_name,omitempty"`
	DependencyType string `json:"dependency_type,omitempty"`
}

type Hashes struct {
	Sha512 string `json:"sha512,omitempty"`
	Sha1   string `json:"sha1,omitempty"`
}

type File struct {
	Hashes   Hashes `json:"hashes,omitempty"`
	Url      string `json:"url,omitempty"`
	Filename string `json:"filename,omitempty"`
	Primary  bool   `json:"primary,omitempty"`
	Size     uint   `json:"size,omitempty"`
	FileType string `json:"file_type,omitempty"`
}

type Version struct {
	Name            string       `json:"name,omitempty"`
	VersionNumber   string       `json:"version_number,omitempty"`
	Changelog       string       `json:"changelog,omitempty"`
	Dependencies    []Dependency `json:"dependencies,omitempty"`
	GameVersions    []string     `json:"game_versions,omitempty"`
	VersionType     string       `json:"version_type,omitempty"`
	Loaders         []string     `json:"loaders,omitempty"`
	Featured        bool         `json:"featured,omitempty"`
	Status          string       `json:"status,omitempty"`
	RequestedStatus string       `json:"requested_status,omitempty"`
	Id              string       `json:"id,omitempty"`
	ProjectId       string       `json:"project_id,omitempty"`
	AuthorId        string       `json:"author_id,omitempty"`
	DatePublished   time.Time    `json:"date_published,omitempty"`
	Downloads       uint         `json:"downloads,omitempty"`
	Files           []File       `json:"files,omitempty"`
}

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

func (v Version) GetPrimaryFile() *File {
	if len(v.Files) == 1 {
		return &v.Files[0]
	}

	primaryIndex := 0

	for i, f := range v.Files {
		if f.Primary {
			primaryIndex = i
			break
		}
	}

	return &v.Files[primaryIndex]
}

func (f File) Download() (io.ReadCloser, error) {
	resp, getErr := http.Get(f.Url)

	if getErr != nil {
		return nil, getErr
	}

	return resp.Body, nil
}
