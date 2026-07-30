package pypi

import (
	"fmt"
	"strings"
	"time"
)

type Info struct {
	Name            string            `json:"name"`
	Author          string            `json:"author"`
	AuthorEmail     string            `json:"author_email"`
	Maintainer      string            `json:"maintainer"`
	MaintainerEmail string            `json:"maintainer_email"`
	Description     string            `json:"description"`
	License         string            `json:"license"`
	Classifiers     []string          `json:"classifiers"`
	ProjectURL      string            `json:"project_url"`
	PackageURL      string            `json:"package_url"`
	HomePage        string            `json:"home_page"`
	ProjectURLs     map[string]string `json:"project_urls"`
	Version         string            `json:"version"`
}

type ReleaseFile struct {
	UploadTime   time.Time `json:"upload_time_iso_8601"`
	Yanked       bool      `json:"yanked"`
	YankedReason string    `json:"yanked_reason"`
}

type Package struct {
	Info     Info                     `json:"info"`
	Releases map[string][]ReleaseFile `json:"releases"`
}

func (p *Package) Name() string {
	return p.Info.Name
}

func (p *Package) Author() string {
	return strings.Join(p.MaintainerNames(), ", ")
}

func (p *Package) MaintainerNames() []string {
	values := make([]string, 0, 2)
	seen := make(map[string]bool)
	for _, value := range []string{p.Info.Maintainer, p.Info.Author} {
		value = strings.TrimSpace(value)
		if value != "" && !seen[value] {
			seen[value] = true
			values = append(values, value)
		}
	}
	return values
}

func (p *Package) Description() string {
	return p.Info.Description
}

func (p *Package) License() string {
	if value := strings.TrimSpace(p.Info.License); value != "" {
		return value
	}
	for _, classifier := range p.Info.Classifiers {
		if strings.HasPrefix(classifier, "License ::") {
			return classifier
		}
	}
	return ""
}

func (p *Package) ProjectURL() string {
	return p.Info.ProjectURL
}

func (p *Package) PackageURL() string {
	if value := strings.TrimSpace(p.Info.PackageURL); value != "" {
		return value
	}
	return fmt.Sprintf("https://pypi.org/project/%s/", p.Name())
}

func (p *Package) RepositoryURL() string {
	for _, preferred := range []string{"source", "repository", "code", "github"} {
		for label, value := range p.Info.ProjectURLs {
			if strings.EqualFold(strings.TrimSpace(label), preferred) && strings.TrimSpace(value) != "" {
				return strings.TrimSpace(value)
			}
		}
	}
	return strings.TrimSpace(p.Info.HomePage)
}

func (p *Package) Version() string {
	return p.Info.Version
}

func (p *Package) Type() string {
	return Type
}
