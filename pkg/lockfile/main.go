package lockfile

import (
	"git.sr.ht/~ansipunk/weaver/pkg/modrinth"
	"git.sr.ht/~ansipunk/weaver/pkg/types"
)

type VersionLock struct {
	VersionID    string   `toml:"version_id"`
	ProjectSlug  string   `toml:"project_slug"`
	Filename     string   `toml:"filename"`
	Dependencies []string `toml:"dependencies"`
	Hash         string   `toml:"hash"`
}

func flattenDependencies(version *types.Version) ([]string, error) {
	flat := []string{}

	for _, dependency := range version.Dependencies {
		projectID := dependency.ProjectID

		if projectID == "" {
			version, err := modrinth.GetSpecificVersion(dependency.VersionID)
			if err != nil {
				return flat, err
			}
			projectID = version.ProjectID
		}

		project, err := modrinth.GetProject(projectID)
		if err != nil {
			return flat, err
		}
		flat = append(flat, project.Slug)
	}

	return flat, nil
}

func LockVersion(version *types.Version, projectSlug string, filename string, hash string) (VersionLock, error) {
	dependencies, err := flattenDependencies(version)
	if err != nil {
		return VersionLock{}, err
	}

	return VersionLock{
		VersionID:    version.ID,
		ProjectSlug:  projectSlug,
		Filename:     filename,
		Dependencies: dependencies,
		Hash:         hash,
	}, nil
}
