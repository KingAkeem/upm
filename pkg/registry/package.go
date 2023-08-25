package registry

import "fmt"

type PackageNotFoundError struct {
	name         string
	registryType string
}

func (p PackageNotFoundError) Error() string {
	return fmt.Sprintf("Package %s was not found in %s", p.name, p.registryType)
}

type Package interface {
	Type() string
	Name() string
	Author() string
	Description() string
	License() string
	Version() string
}
