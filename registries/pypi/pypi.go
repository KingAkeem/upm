package pypi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"upm/registry"
)

var Type string = "pypi"

type PypiRegistry struct {
	URL string
}

func init() {
	registry.Register(Type, NewPypiRegistry)
}

func (p *PypiRegistry) Get(name string) (registry.Package, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name must be specified")
	}

	uri, err := url.JoinPath(p.URL, "pypi", name, "json")
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	item := new(Package)
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (p *PypiRegistry) Type() string {
	return Type
}

func NewPypiRegistry(settings map[string]interface{}) registry.Registry {
	return &PypiRegistry{URL: "https://pypi.org"}
}
