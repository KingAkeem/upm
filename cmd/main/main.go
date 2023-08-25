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
var GET = "get"
var PUBLISH = "publish"
var AvailableActions = []string{GET, PUBLISH}

func main() {
	actionArg := flag.String("a", GET, "Action to perform on package")
	registriesArg := flag.String("r", "all", "Which registries to use")
	packageName := flag.String("n", "", "Package name")
	username := flag.String("u", "", "username for registry account")
	password := flag.String("p", "", "password for registry account")
	flag.Parse()

	action := strings.ToLower(strings.TrimSpace(*actionArg))
	if action == "" {
		fmt.Println("Error: Action is blank. Please provide a valid action.")
		os.Exit(0)
	}

	if strings.TrimSpace(*packageName) == "" && action != PUBLISH {
		fmt.Println("Error: Package name is blank. Please provide a package name.")
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
			fmt.Printf("Error: No valid registry passed, found %v\n", registries)
			os.Exit(0)
		}

	} else if !checkRegistries(registries) {
		fmt.Printf("Error: No valid registry passed, found %v\n", registries)
		os.Exit(0)
	}

	if !slices.Contains(AvailableActions, action) {
		fmt.Printf("Error: Invalid action found, found %s.\n", action)
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

		switch action {
		case PUBLISH:
			fmt.Printf("Performing '%s'.\n", action)
			err = r.Publish(*username, *password)
			if err != nil {
				fmt.Printf("Error: Unable to publish to %s, %s\n", registryType, err.Error())
				continue
			}
		case GET:
			fmt.Printf("Performing '%s' for package %s.\n", action, *packageName)
			project, err := r.Get(*packageName)
			switch {
			case errors.As(err, &registry.PackageNotFoundError{}):
				fmt.Printf("Error: %s not found within %s.\n\n", *packageName, registryType)
				continue
			case err == nil:
				fmt.Printf("Success: %s has been found within %s.\n", *packageName, registryType)
				fmt.Printf("%s\n\n", toString(project))
			default:
				fmt.Printf("Error: Unable to get %s from %s, %s.\n", *packageName, registryType, err.Error())
				continue
			}
		default:
			fmt.Printf("Error: Invalid action found, %s\n\n", action)
			continue
		}
	}
}
