package npm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"upm/pkg/registry"
)

var Type string = "npm"

type NpmRegistry struct {
	URL string
}

func init() {
	registry.Register(Type, NewNpmRegistry)
}

func (n *NpmRegistry) Get(name string) (registry.Package, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name must be specified")
	}

	uri, err := url.JoinPath(n.URL, name)
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

func (n *NpmRegistry) Type() string {
	return Type
}

func NewNpmRegistry(settings map[string]interface{}) registry.Registry {
	return &NpmRegistry{URL: "https://registry.npmjs.com"}
}
