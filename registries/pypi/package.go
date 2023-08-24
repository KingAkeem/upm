package pypi

import "fmt"

type Info struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Email       string `json:"author_email"`
	Description string `json:"description"`
	License     string `json:"license"`
	ProjectURL  string `json:"project_url"`
	Version     string `json:"version"`
}

type Package struct {
	Info Info `json:"info"`
}

func (p *Package) Name() string {
	return p.Info.Name
}

func (p *Package) Author() string {
	return fmt.Sprintf("%s - %s", p.Info.Author, p.Info.Email)
}
func (p *Package) Description() string {
	return p.Info.Description
}
func (p *Package) License() string {
	return p.Info.License
}
func (p *Package) ProjectURL() string {
	return p.Info.ProjectURL
}
func (p *Package) Version() string {
	return p.Info.Version
}
func (p *Package) Type() string {
	return Type
}
