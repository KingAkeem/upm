package registry

import (
	"fmt"
	"time"
)

type Registry interface {
	Type() string
	Fetch(name string) (Package, error)
	Inspect(name string, now time.Time) (Inspection, error)
	Publish(options PublishOptions) error
}

type PublishOptions struct {
	Package string
	Account string
}

type RegistryFactory func(settings map[string]interface{}) Registry

var registry = make(map[string]RegistryFactory)

func Register(registryType string, r RegistryFactory) {
	registry[registryType] = r
}

func Create(registryType string, settings map[string]interface{}) (Registry, error) {
	factory, exists := registry[registryType]
	if !exists {
		return nil, fmt.Errorf("registry type '%s' can not be found", registryType)
	}
	return factory(settings), nil
}
