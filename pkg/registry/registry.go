package registry

import "fmt"

type Registry interface {
	Type() string
	Get(name string) (Package, error)
	Publish(username, password string) error
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
