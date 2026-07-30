package npm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
	"upm/pkg/registry"
)

var Type string = "npm"

type NpmRegistry struct {
	URL    string
	Client *http.Client
}

func init() {
	registry.Register(Type, NewNpmRegistry)
}

func (n *NpmRegistry) Fetch(name string) (registry.Package, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name must be specified")
	}

	uri, err := url.JoinPath(n.URL, name)
	if err != nil {
		return nil, err
	}

	resp, err := n.Client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, registry.NewPackageNotFoundError(name, Type)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("npm registry returned %s", resp.Status)
	}

	item := new(Package)
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (n *NpmRegistry) Inspect(name string, now time.Time) (registry.Inspection, error) {
	value, err := n.Fetch(name)
	if err != nil {
		return registry.Inspection{}, err
	}
	item := value.(*Package)
	latest := item.Latest()
	deprecated := strings.TrimSpace(latest.Deprecated) != ""
	inspection := registry.Inspection{
		Registry:      Type,
		Name:          item.Name(),
		LatestVersion: item.Version(),
		License:       item.License(),
		Maintainers:   item.MaintainerNames(),
		Deprecated:    &deprecated,
		PackageURL:    "https://www.npmjs.com/package/" + item.Name(),
		RepositoryURL: item.RepositoryURL(),
	}
	if publishedAt, ok := item.Times[item.Version()]; ok {
		inspection.PublishedAt = &publishedAt
	}
	if updatedAt, ok := item.Times["modified"]; ok {
		inspection.UpdatedAt = &updatedAt
	}
	inspection.InferWarnings(now)
	return inspection, nil
}

func (n *NpmRegistry) Type() string {
	return Type
}

func (n *NpmRegistry) Publish(options registry.PublishOptions) error {
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
	value := &NpmRegistry{URL: "https://registry.npmjs.com", Client: http.DefaultClient}
	if configuredURL, ok := settings["url"].(string); ok && configuredURL != "" {
		value.URL = configuredURL
	}
	if client, ok := settings["client"].(*http.Client); ok && client != nil {
		value.Client = client
	}
	return value
}
