package main

import (
	"errors"
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

	packageData := ""
	i := 0
	for metadataKey, metadataValue := range metadata {
		packageData += fmt.Sprintf("%s: %s", metadataKey, metadataValue)
		i++
		if i != len(metadata) {
			packageData += "\n"
		}
	}
	return packageData
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
var fetch = "fetch"
var publish = "publish"
var AvailableActions = []string{fetch, publish}

func main() {
	if len(os.Args) < 2 {
		os.Exit(0)
	}

	action := strings.ToLower(strings.TrimSpace(os.Args[1]))
	if action == "" {
		fmt.Println("Error: Action is blank. Please provide a valid action.")
		os.Exit(0)
	}

	var registriesArg string
	var packageName string
	var username string
	var password string
	for i, arg := range os.Args[1:] {
		switch strings.TrimSpace(arg) {
		case "-n":
			packageName = os.Args[i+2]
		case "-r":
			registriesArg = os.Args[i+2]
		case "-u":
			username = os.Args[i+2]
		case "-p":
			password = os.Args[i+2]

		}
	}

	// defaults to all available registries
	if strings.TrimSpace(registriesArg) == "" {
		registriesArg = "all"
	}

	// publish action doesn't require a package name
	if strings.TrimSpace(packageName) == "" && action != publish {
		fmt.Println("Error: Package name is blank. Please provide a package name.")
		os.Exit(0)
	}

	registries := strings.Split(strings.ToLower(strings.TrimSpace(registriesArg)), ",")
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
		case publish:
			fmt.Printf("Performing '%s'.\n", action)
			err = r.Publish(username, password)
			if err != nil {
				fmt.Printf("Error: Unable to publish to %s, %s\n", registryType, err.Error())
				continue
			}
		case fetch:
			fmt.Printf("Performing '%s' for package '%s'.\n", action, packageName)
			project, err := r.Fetch(packageName)
			switch {
			case errors.As(err, &registry.PackageNotFoundError{}):
				fmt.Printf("Error: %s not found within %s.\n\n", packageName, registryType)
				continue
			case err == nil:
				fmt.Printf("Success: '%s' has been found within %s.\n", packageName, registryType)
				fmt.Println("------------------------")
				fmt.Printf("%s\n", toString(project))
				fmt.Println("------------------------")
			default:
				fmt.Printf("Error: Unable to fetch %s from %s, %s.\n", packageName, registryType, err.Error())
				continue
			}
		default:
			fmt.Printf("Error: Invalid action found, %s\n\n", action)
			continue
		}
	}
}
