package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
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

	if strings.TrimSpace(*packageName) == "" {
		fmt.Println("Error: Package name is blank. Please provide a package name.")
		os.Exit(0)
	}

	action := strings.ToLower(strings.TrimSpace(*actionArg))
	if action == "" {
		fmt.Println("Error: Action is blank. Please provide a valid action.")
		os.Exit(0)
	}

	registries := strings.Split(strings.ToLower(strings.TrimSpace(*registriesArg)), ",")
	if len(registries) < 1 {
		fmt.Printf("Error: Malformed registries argument passed, found %v.\n", registries)
		os.Exit(0)
	}

	if len(registries) == 1 {
		// a single registry could refer to "all" or be a single registry within the list
		if !checkRegistries(registries) && registries[0] != "all" {
			fmt.Printf("No valid registry passed, found %v\n", registries)
			os.Exit(0)
		}

	} else if !checkRegistries(registries) {
		fmt.Printf("No valid registry passed, found %v\n", registries)
		os.Exit(0)
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
			fmt.Printf("Error: Unable to create %s registry, %s\n", registryType, err.Error())
			continue
		}
		fmt.Printf("Success: Registry created for %s.\n", registryType)

		fmt.Printf("Performing action '%s' for package %s.\n", action, *packageName)
		switch action {
		case "get":
			project, err := r.Get(*packageName)
			switch {
			case errors.As(err, &registry.PackageNotFoundError{}):
				fmt.Printf("Error: %s not found within %s.\n\n", *packageName, registryType)
				continue
			case err == nil:
				fmt.Printf("Success: %s has been found within %s.\n", *packageName, registryType)
				fmt.Printf("%s\n\n", toString(project))
			default:
				fmt.Printf("Error: Unable to get %s from %s, %s", *packageName, registryType, err.Error())
				continue
			}
		default:
			fmt.Printf("Error: Invalid action found, %s\n\n", action)
			continue
		}
	}
}
