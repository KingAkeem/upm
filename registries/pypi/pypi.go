package pypi

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

var Type string = "pypi"

type PypiRegistry struct {
	URL string
}

func init() {
	registry.Register(Type, NewPypiRegistry)
}

func (p *PypiRegistry) Publish(username, password string) error {
	buildCmd := exec.Command("python", "setup.py", "sdist", "bdist_wheel")
	result, err := buildCmd.Output()
	if err != nil {
		return err
	}
	fmt.Print(string(result))
	fmt.Println("Success: Python package successfully built.")

	uploadCmd := exec.Command("twine", "upload", "-u", username, "-p", password, "dist/*")

	uploadCmd.Stdout = os.Stdout
	uploadCmd.Stderr = os.Stderr

	err = uploadCmd.Start()
	if err != nil {
		return err
	}

	err = uploadCmd.Wait()
	if err != nil {
		return err
	}

	fmt.Print(string(result))
	fmt.Printf("Success: Python package successfuly uploaded to %s.\n", Type)
	return nil
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

func (p *PypiRegistry) Type() string {
	return Type
}

func NewPypiRegistry(settings map[string]interface{}) registry.Registry {
	return &PypiRegistry{URL: "https://pypi.org"}
}
