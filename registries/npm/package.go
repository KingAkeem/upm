package npm

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type TextValue string

func (value *TextValue) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		*value = TextValue(text)
		return nil
	}
	var object struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &object); err != nil {
		return err
	}
	*value = TextValue(object.Type)
	return nil
}

type Person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (person *Person) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		parts := strings.SplitN(text, "<", 2)
		person.Name = strings.TrimSpace(parts[0])
		return nil
	}
	type alias Person
	return json.Unmarshal(data, (*alias)(person))
}

func (person Person) DisplayName() string {
	if strings.TrimSpace(person.Name) != "" {
		return strings.TrimSpace(person.Name)
	}
	return strings.TrimSpace(person.Email)
}

type Repository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

func (repository *Repository) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		repository.URL = text
		return nil
	}
	type alias Repository
	return json.Unmarshal(data, (*alias)(repository))
}

func (repository Repository) WebURL() string {
	value := strings.TrimSpace(repository.URL)
	value = strings.TrimPrefix(value, "git+")
	value = strings.TrimSuffix(value, ".git")
	if strings.HasPrefix(value, "git://github.com/") {
		return "https://github.com/" + strings.TrimPrefix(value, "git://github.com/")
	}
	if strings.HasPrefix(value, "git@github.com:") {
		return "https://github.com/" + strings.TrimPrefix(value, "git@github.com:")
	}
	return value
}

type Version struct {
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	License     TextValue  `json:"license"`
	Deprecated  string     `json:"deprecated"`
	Repository  Repository `json:"repository"`
	NodeVersion string     `json:"_nodeVersion"`
	NpmVersion  string     `json:"_npmVersion"`
}

type Package struct {
	ID                 string             `json:"_id"`
	PackageName        string             `json:"name"`
	Versions           map[string]Version `json:"versions"`
	PackageDescription string             `json:"description"`
	PackageAuthor      Person             `json:"author"`
	PackageLicense     TextValue          `json:"license"`
	Maintainers        []Person           `json:"maintainers"`
	Repository         Repository         `json:"repository"`
	DistTags           struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
	Times map[string]time.Time `json:"time"`
}

func (p *Package) Name() string {
	return p.PackageName
}

func (p *Package) Author() string {
	if name := p.PackageAuthor.DisplayName(); name != "" {
		return name
	}
	maintainers := p.MaintainerNames()
	return strings.Join(maintainers, ", ")
}

func (p *Package) MaintainerNames() []string {
	seen := make(map[string]bool)
	values := make([]string, 0, len(p.Maintainers)+1)
	for _, person := range append([]Person{p.PackageAuthor}, p.Maintainers...) {
		name := person.DisplayName()
		if name != "" && !seen[name] {
			seen[name] = true
			values = append(values, name)
		}
	}
	return values
}

func (p *Package) Description() string {
	return p.PackageDescription
}

func (p *Package) License() string {
	if value := strings.TrimSpace(string(p.PackageLicense)); value != "" {
		return value
	}
	return strings.TrimSpace(string(p.Latest().License))
}

func (p *Package) Version() string {
	if p.DistTags.Latest != "" {
		return p.DistTags.Latest
	}
	for version := range p.Versions {
		return version
	}
	return ""
}

func (p *Package) Latest() Version {
	return p.Versions[p.Version()]
}

func (p *Package) RepositoryURL() string {
	if value := p.Latest().Repository.WebURL(); value != "" {
		return value
	}
	return p.Repository.WebURL()
}

func (p *Package) Type() string {
	return Type
}

func (p *Package) String() string {
	return fmt.Sprintf("%s@%s", p.PackageName, p.Version())
}
