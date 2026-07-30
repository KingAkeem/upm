package registry

import "fmt"

type PackageNotFoundError struct {
	Name         string
	RegistryType string
}

func (p PackageNotFoundError) Error() string {
	return fmt.Sprintf("package %s was not found in %s", p.Name, p.RegistryType)
}

func NewPackageNotFoundError(name, registryType string) error {
	return PackageNotFoundError{Name: name, RegistryType: registryType}
}

type Package interface {
	Type() string
	Name() string
	Author() string
	Description() string
	License() string
	Version() string
}
