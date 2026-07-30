package pypi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"upm/pkg/registry"
)

var Type string = "pypi"

type PypiRegistry struct {
	URL    string
	Client *http.Client
}

func init() {
	registry.Register(Type, NewPypiRegistry)
}

func (p *PypiRegistry) Publish(options registry.PublishOptions) error {
	python := strings.TrimSpace(os.Getenv("UPM_PYTHON"))
	if python == "" {
		python = "python3"
	}
	buildCmd := exec.Command(python, "setup.py", "sdist", "bdist_wheel")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	err := buildCmd.Start()
	if err != nil {
		return err
	}
	err = buildCmd.Wait()
	if err != nil {
		return err
	}
	fmt.Println("Success: Python package successfully built.")

	distributions, err := filepath.Glob("dist/*")
	if err != nil {
		return err
	}
	if len(distributions) == 0 {
		return fmt.Errorf("no distributions found under dist/")
	}
	uploadCmd := exec.Command("twine", append([]string{"upload"}, distributions...)...)
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
	fmt.Printf("Success: Python package successfuly uploaded to %s.\n", Type)
	return nil
}

func (p *PypiRegistry) Fetch(name string) (registry.Package, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name must be specified")
	}

	uri, err := url.JoinPath(p.URL, "pypi", name, "json")
	if err != nil {
		return nil, err
	}

	resp, err := p.Client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, registry.NewPackageNotFoundError(name, Type)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("pypi registry returned %s", resp.Status)
	}

	item := new(Package)
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (p *PypiRegistry) Inspect(name string, now time.Time) (registry.Inspection, error) {
	value, err := p.Fetch(name)
	if err != nil {
		return registry.Inspection{}, err
	}
	item := value.(*Package)
	latestFiles := item.Releases[item.Version()]
	yanked := len(latestFiles) > 0
	var publishedAt *time.Time
	for _, file := range latestFiles {
		if !file.Yanked {
			yanked = false
		}
		if !file.UploadTime.IsZero() && (publishedAt == nil || file.UploadTime.Before(*publishedAt)) {
			value := file.UploadTime
			publishedAt = &value
		}
	}
	var updatedAt *time.Time
	for _, files := range item.Releases {
		for _, file := range files {
			if !file.UploadTime.IsZero() && (updatedAt == nil || file.UploadTime.After(*updatedAt)) {
				value := file.UploadTime
				updatedAt = &value
			}
		}
	}
	inspection := registry.Inspection{
		Registry:      Type,
		Name:          item.Name(),
		LatestVersion: item.Version(),
		License:       item.License(),
		Maintainers:   item.MaintainerNames(),
		PublishedAt:   publishedAt,
		UpdatedAt:     updatedAt,
		Yanked:        &yanked,
		PackageURL:    item.PackageURL(),
		RepositoryURL: item.RepositoryURL(),
	}
	inspection.InferWarnings(now)
	return inspection, nil
}

func (p *PypiRegistry) Type() string {
	return Type
}

func NewPypiRegistry(settings map[string]interface{}) registry.Registry {
	value := &PypiRegistry{URL: "https://pypi.org", Client: http.DefaultClient}
	if configuredURL, ok := settings["url"].(string); ok && configuredURL != "" {
		value.URL = configuredURL
	}
	if client, ok := settings["client"].(*http.Client); ok && client != nil {
		value.Client = client
	}
	return value
}
