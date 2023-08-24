package npm

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
)

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Version struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	License     string `json:"license"`
	NodeVersion string `json:"_nodeVersion"`
	NpmVersion  string `json:"_npmVersion"`
}

type Maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Package struct {
	ID                 string             `json:"_id"`
	PackageName        string             `json:"name"`
	Versions           map[string]Version `json:"versions"`
	PackageDescription string             `json:"description"`
	PackageAuthor      Author             `json:"author"`
	PackageLicense     string             `json:"license"`
	Maintainers        []Maintainer       `json:"maintainers"`
}

func (p *Package) Name() string {
	return p.PackageName
}
func (p *Package) Author() string {
	if strings.TrimSpace(p.PackageAuthor.Name) != "" {
		return p.PackageAuthor.Name
	}

	maintainers := []string{}
	for _, m := range p.Maintainers {
		maintainers = append(maintainers, fmt.Sprintf("%s - %s", m.Name, m.Email))
	}
	return strings.Join(maintainers, ", ")
}
func (p *Package) Description() string {
	return p.PackageDescription
}
func (p *Package) License() string {
	return p.PackageLicense
}
func (p *Package) Version() string {
	if len(p.Versions) == 0 {
		return ""
	}

	return maps.Values(p.Versions)[len(p.Versions)-1].Version
}
func (p *Package) Type() string {
	return Type
}
