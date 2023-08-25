package npm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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

	if resp.StatusCode == http.StatusNotFound {
		return nil, registry.PackageNotFoundError{}
	}

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

func (n *NpmRegistry) Publish(username, password string) error {
	publishCmd := exec.Command("npm", "publish")
	publishCmd.Stdout = os.Stdout
	publishCmd.Stderr = os.Stderr
	err := publishCmd.Start()
	if err != nil {
		return err
	}
	err = publishCmd.Wait()
	if err != nil {
		return err
	}
	fmt.Printf("Success: NPM package successfully uploaded to %s.\n", Type)
	return nil
}

func NewNpmRegistry(settings map[string]interface{}) registry.Registry {
	return &NpmRegistry{URL: "https://registry.npmjs.com"}
}
