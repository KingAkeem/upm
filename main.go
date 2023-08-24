package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"upm/registries/npm"
	"upm/registries/pypi"
	"upm/registry"

	"golang.org/x/exp/slices"
)

func toString(p registry.Package) string {
	metadata := map[string]string{
		"Name":        p.Name(),
		"Author":      p.Author(),
		"Description": p.Description(),
		"License":     p.License(),
		"Version":     p.Version(),
	}
	s, _ := json.MarshalIndent(metadata, "", "\t")
	return string(s)
}

func checkRegistries(userRegistries, availableRegistries []string) bool {
	for _, r := range userRegistries {
		if slices.Contains(availableRegistries, r) {
			return true
		}
	}

	return false
}

func handleAction(r registry.Registry, action string, packageName string) {
	switch action {
	case "get":
		project, err := r.Get(packageName)
		if err != nil {
			panic(err)
		}
		fmt.Println(toString(project))

	}
}

func main() {
	var action string
	flag.StringVar(&action, "a", "GET", "Action to be performed by upm")

	var registriesArg string
	flag.StringVar(&registriesArg, "r", "all", "Regstries to use for action")

	var moduleName string
	flag.StringVar(&moduleName, "n", "", "Name of the module to perform action upon")

	flag.Parse()

	action = strings.ToLower(strings.TrimSpace(action))
	if action == "" {
		panic("action must be specified")
	}

	if registriesArg == "" {
		panic("registries must be specified")
	}

	registries := strings.Split(registriesArg, ",")
	if len(registries) < 1 {
		panic(fmt.Errorf("registries argument malformed, found %s", registriesArg))
	}

	availableRegistries := []string{npm.Type, pypi.Type}
	if len(registries) == 1 {
		if !checkRegistries(registries, availableRegistries) && registries[0] != "all" {
			panic(fmt.Errorf("no valid registry passed, found %v", registries))
		}

	} else if !checkRegistries(registries, availableRegistries) {
		panic(fmt.Errorf("no valid registry passed, found %v", registries))
	}

	if registries[0] == "all" {
		// search all available registries
		for _, registryType := range availableRegistries {
			r, err := registry.Create(registryType, nil)
			if err != nil {
				panic(err)
			}
			handleAction(r, action, moduleName)
		}
	} else {
		// search user specified registries
		for _, registryType := range registries {
			r, err := registry.Create(registryType, nil)
			if err != nil {
				panic(err)
			}
			handleAction(r, action, moduleName)
		}
	}
}
