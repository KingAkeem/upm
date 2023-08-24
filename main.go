package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"upm/pkg/registry"
	"upm/registries/npm"
	"upm/registries/pypi"

	"golang.org/x/exp/slices"
)

func toString(p registry.Package) string {
	metadata := map[string]string{
		"Registry":    p.Type(),
		"Name":        p.Name(),
		"Author":      p.Author(),
		"Description": p.Description(),
		"License":     p.License(),
		"Version":     p.Version(),
	}
	s, _ := json.MarshalIndent(metadata, "", "\t")
	return string(s)
}

func checkRegistries(registries []string) bool {
	for _, r := range registries {
		if slices.Contains(AvailableRegistries, r) {
			return true
		}
	}
	return false
}

var AvailableRegistries = []string{npm.Type, pypi.Type}

func main() {
	actionArg := flag.String("a", "GET", "Action to perform on package")
	registriesArg := flag.String("r", "all", "Which registries to use")
	packageName := flag.String("p", "", "Package name")
	flag.Parse()

	action := strings.ToLower(strings.TrimSpace(*actionArg))
	if action == "" {
		panic("action must be specified")
	}

	registries := strings.Split(strings.ToLower(strings.TrimSpace(*registriesArg)), ",")
	if len(registries) < 1 {
		panic(fmt.Errorf("registries argument malformed, found %s", registries))
	}

	if len(registries) == 1 {
		// a single registry could refer to "all" or be a single registry within the list
		if !checkRegistries(registries) && registries[0] != "all" {
			panic(fmt.Errorf("no valid registry passed, found %v", registries))
		}

	} else if !checkRegistries(registries) {
		panic(fmt.Errorf("no valid registry passed, found %v", registries))
	}

	var registryList []string
	if registries[0] == "all" {
		registryList = AvailableRegistries
	} else {
		registryList = registries
	}

	for _, registryType := range registryList {
		fmt.Printf("Creating registry for %s.\n", registryType)
		r, err := registry.Create(registryType, nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Success: Registry created for %s.\n", registryType)

		fmt.Printf("Performing action '%s' for package %s.\n", action, *packageName)
		switch action {
		case "get":
			project, err := r.Get(*packageName)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Success: %s has been found within %s.\n", *packageName, registryType)
			fmt.Printf("%s\n\n", toString(project))
		default:
			panic(fmt.Errorf("invalid action found, %v", err))
		}
	}
}
