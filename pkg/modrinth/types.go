package modrinth

import (
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Dependency struct {
	VersionID      string `json:"version_id,omitempty"`
	ProjectID      string `json:"project_id,omitempty"`
	FileName       string `json:"file_name,omitempty"`
	DependencyType string `json:"dependency_type,omitempty"`
}

type Hashes struct {
	Sha512 string `json:"sha512,omitempty"`
	Sha1   string `json:"sha1,omitempty"`
}

type File struct {
	Hashes   Hashes `json:"hashes,omitempty"`
	URL      string `json:"url,omitempty"`
	Filename string `json:"filename,omitempty"`
	Primary  bool   `json:"primary,omitempty"`
	Size     int64  `json:"size,omitempty"`
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
	ID              string       `json:"id,omitempty"`
	ProjectID       string       `json:"project_id,omitempty"`
	AuthorID        string       `json:"author_id,omitempty"`
	DatePublished   time.Time    `json:"date_published,omitempty"`
	Downloads       uint         `json:"downloads,omitempty"`
	Files           []File       `json:"files,omitempty"`
	Slug            string       `json:",omitempty"`
}

type DonationURL struct {
	ID       string `json:"id,omitempty"`
	Platform string `json:"platform,omitempty"`
	URL      string `json:"url,omitempty"`
}

type License struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type GalleryImage struct {
	URL      string `json:"url,omitempty"`
	Featured bool   `json:"featured,omitempty"`
	Title    string `json:"title,omitempty"`
	Desc     string `json:"description,omitempty"`
	Created  string `json:"created,omitempty"`
	Ordering int    `json:"ordering,omitempty"`
}

type Project struct {
	Slug                 string         `json:"slug,omitempty"`
	Title                string         `json:"title,omitempty"`
	Description          string         `json:"description,omitempty"`
	Categories           []string       `json:"categories,omitempty"`
	ClientSide           bool           `json:"client_side,omitempty"`
	ServerSide           bool           `json:"server_side,omitempty"`
	Body                 string         `json:"body,omitempty"`
	AdditionalCategories []string       `json:"additional_categories,omitempty"`
	IssuesURL            string         `json:"issues_url,omitempty"`
	SourceURL            string         `json:"source_url,omitempty"`
	WikiURL              string         `json:"wiki_url,omitempty"`
	DiscordURL           string         `json:"discord_url,omitempty"`
	DonationURLs         []DonationURL  `json:"donation_urls,omitempty"`
	ProjectType          string         `json:"project_type,omitempty"`
	Downloads            int            `json:"downloads,omitempty"`
	IconURL              string         `json:"icon_url,omitempty"`
	Color                int            `json:"color,omitempty"`
	ID                   string         `json:"id,omitempty"`
	Team                 string         `json:"team,omitempty"`
	BodyURL              *string        `json:"body_url,omitempty"`
	ModeratorMsg         *string        `json:"moderator_message,omitempty"`
	Published            string         `json:"published,omitempty"`
	Updated              string         `json:"updated,omitempty"`
	Approved             string         `json:"approved,omitempty"`
	Followers            int            `json:"followers,omitempty"`
	Status               string         `json:"status,omitempty"`
	License              License        `json:"license,omitempty"`
	Versions             []string       `json:"versions,omitempty"`
	GameVersions         []string       `json:"game_versions,omitempty"`
	Loaders              []string       `json:"loaders,omitempty"`
	Gallery              []GalleryImage `json:"gallery,omitempty"`
}

func (version *Version) GetDependencies() ([]Version, error) {
	dependencyVersions := []Version{}

	for _, dependency := range version.Dependencies {
		if dependency.VersionID != "" {
			specificVersion, err := GetSpecificVersion(dependency.VersionID)
			if err != nil {
				return dependencyVersions, errors.Wrap(err, "failed to get specific version")
			}
			dependencyVersions = append(dependencyVersions, specificVersion)
		}
	}

	version.Dependencies = nil
	return dependencyVersions, nil
}

func (version *Version) GetPrimaryFile() *File {
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

func (file *File) Download() (io.ReadCloser, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(file.URL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download file")
	}

	return resp.Body, nil
}

func (version *Version) SetProjectSlug() error {
	project, err := GetProject(version.ProjectID)
	if err != nil {
		return errors.Wrap(err, "failed to get project")
	}

	version.Slug = project.Slug
	return nil
}
