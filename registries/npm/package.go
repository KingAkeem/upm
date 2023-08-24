package npm

import (
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
	Author      Author `json:"author"`
	License     string `json:"license"`
	NodeVersion string `json:"_nodeVersion"`
	NpmVersion  string `json:"_npmVersion"`
}

type Package struct {
	ID                 string             `json:"_id"`
	PackageName        string             `json:"name"`
	Versions           map[string]Version `json:"versions"`
	PackageDescription string             `json:"description"`
	PackageAuthor      Author             `json:"author"`
	PackageLicense     string             `json:"license"`
}

func (p *Package) Name() string {
	return p.PackageName
}
func (p *Package) Author() string {
	return p.PackageAuthor.Name
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
