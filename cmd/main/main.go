package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"upm/pkg/registry"
	"upm/registries/npm"
	"upm/registries/pypi"
)

const (
	fetchAction   = "fetch"
	inspectAction = "inspect"
	publishAction = "publish"
)

var availableRegistries = []string{npm.Type, pypi.Type}

type commandOptions struct {
	action      string
	packageName string
	registries  []string
	account     string
	dryRun      bool
	execute     bool
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	options, err := parseOptions(args, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return 2
	}

	exitCode := 0
	for _, registryType := range options.registries {
		value, err := registry.Create(registryType, nil)
		if err != nil {
			fmt.Fprintf(stderr, "Error: unable to create %s registry: %v\n", registryType, err)
			exitCode = 1
			continue
		}
		switch options.action {
		case fetchAction:
			project, fetchErr := value.Fetch(options.packageName)
			if fetchErr != nil {
				printRegistryError(stderr, registryType, options.packageName, fetchErr)
				exitCode = 1
				continue
			}
			printPackage(stdout, project)
		case inspectAction:
			inspection, inspectErr := value.Inspect(options.packageName, time.Now().UTC())
			if inspectErr != nil {
				printRegistryError(stderr, registryType, options.packageName, inspectErr)
				exitCode = 1
				continue
			}
			printInspection(stdout, inspection)
		case publishAction:
			plan := publishPlan(registryType, options.packageName, options.account)
			printPublishPlan(stdout, plan, !options.execute)
			if !options.execute {
				continue
			}
			if publishErr := value.Publish(registry.PublishOptions{
				Package: options.packageName,
				Account: plan.Account,
			}); publishErr != nil {
				fmt.Fprintf(stderr, "Error: unable to publish %s to %s: %v\n", options.packageName, registryType, publishErr)
				exitCode = 1
			}
		}
	}
	return exitCode
}

func parseOptions(args []string, stderr io.Writer) (commandOptions, error) {
	if len(args) == 0 {
		return commandOptions{}, errors.New("action is required: fetch, inspect, or publish")
	}
	options := commandOptions{action: strings.ToLower(strings.TrimSpace(args[0]))}
	switch options.action {
	case fetchAction, inspectAction, publishAction:
	default:
		return commandOptions{}, fmt.Errorf("unsupported action %q", options.action)
	}

	flags := flag.NewFlagSet(options.action, flag.ContinueOnError)
	flags.SetOutput(stderr)
	var registriesArg string
	flags.StringVar(&options.packageName, "n", "", "package name")
	flags.StringVar(&options.packageName, "name", "", "package name")
	flags.StringVar(&registriesArg, "r", "", "registry list (npm,pypi,all)")
	flags.StringVar(&registriesArg, "registries", "", "registry list (npm,pypi,all)")
	flags.StringVar(&options.account, "account", "", "non-secret account label for publish preview")
	flags.BoolVar(&options.dryRun, "dry-run", false, "preview publishing without running registry commands")
	flags.BoolVar(&options.execute, "execute", false, "perform publishing after previewing the target")
	if err := flags.Parse(args[1:]); err != nil {
		return commandOptions{}, err
	}
	if flags.NArg() > 0 {
		return commandOptions{}, fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}
	options.packageName = strings.TrimSpace(options.packageName)
	if options.packageName == "" {
		return commandOptions{}, errors.New("package name is required; use -n or --name")
	}
	if options.dryRun && options.execute {
		return commandOptions{}, errors.New("--dry-run and --execute cannot be combined")
	}
	if options.action == publishAction && strings.TrimSpace(registriesArg) == "" {
		return commandOptions{}, errors.New("publish requires an explicit registry with -r or --registries")
	}
	registries, err := parseRegistries(registriesArg)
	if err != nil {
		return commandOptions{}, err
	}
	if options.action == publishAction && options.execute && len(registries) != 1 {
		return commandOptions{}, errors.New("--execute publishes to exactly one registry at a time")
	}
	options.registries = registries
	return options, nil
}

func parseRegistries(value string) ([]string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" || value == "all" {
		return append([]string(nil), availableRegistries...), nil
	}
	requested := strings.Split(value, ",")
	result := make([]string, 0, len(requested))
	seen := make(map[string]bool)
	for _, item := range requested {
		item = strings.TrimSpace(item)
		valid := false
		for _, available := range availableRegistries {
			if item == available {
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("unsupported registry %q", item)
		}
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result, nil
}

func printPackage(out io.Writer, value registry.Package) {
	fmt.Fprintf(out, "Registry: %s\n", value.Type())
	fmt.Fprintf(out, "Name: %s\n", value.Name())
	fmt.Fprintf(out, "Author: %s\n", value.Author())
	fmt.Fprintf(out, "Description: %s\n", value.Description())
	fmt.Fprintf(out, "License: %s\n", value.License())
	fmt.Fprintf(out, "Version: %s\n\n", value.Version())
}

func printInspection(out io.Writer, value registry.Inspection) {
	fmt.Fprintln(out, "Verified registry metadata:")
	fmt.Fprintf(out, "  Registry: %s\n", value.Registry)
	fmt.Fprintf(out, "  Name: %s\n", value.Name)
	fmt.Fprintf(out, "  Latest version: %s\n", display(value.LatestVersion))
	fmt.Fprintf(out, "  License: %s\n", display(value.License))
	fmt.Fprintf(out, "  Maintainers/authors: %s\n", display(strings.Join(value.Maintainers, ", ")))
	fmt.Fprintf(out, "  Latest published: %s\n", formatTime(value.PublishedAt))
	fmt.Fprintf(out, "  Last updated: %s\n", formatTime(value.UpdatedAt))
	fmt.Fprintf(out, "  Deprecated: %s\n", formatBool(value.Deprecated))
	fmt.Fprintf(out, "  Yanked: %s\n", formatBool(value.Yanked))
	fmt.Fprintf(out, "  Package URL: %s\n", display(value.PackageURL))
	fmt.Fprintf(out, "  Repository URL: %s\n", display(value.RepositoryURL))
	fmt.Fprintln(out, "Inferred warnings:")
	if len(value.Warnings) == 0 {
		fmt.Fprintln(out, "  None.")
	} else {
		for _, warning := range value.Warnings {
			fmt.Fprintf(out, "  - [%s] %s\n", warning.Code, warning.Message)
		}
	}
	fmt.Fprintln(out)
}

func publishPlan(registryType, packageName, account string) registry.PublishPlan {
	credentialSource := ""
	account = strings.TrimSpace(account)
	switch registryType {
	case npm.Type:
		credentialSource = "npm configuration/environment (run npm login)"
		if account == "" {
			account = strings.TrimSpace(os.Getenv("NPM_USERNAME"))
		}
	case pypi.Type:
		credentialSource = "Twine environment, .pypirc, or keyring"
		if account == "" {
			account = strings.TrimSpace(os.Getenv("TWINE_USERNAME"))
		}
	}
	if account == "" {
		account = "(configured credential identity)"
	}
	return registry.PublishPlan{
		Registry:         registryType,
		Package:          packageName,
		Account:          account,
		CredentialSource: credentialSource,
	}
}

func printPublishPlan(out io.Writer, plan registry.PublishPlan, dryRun bool) {
	mode := "execute"
	if dryRun {
		mode = "dry-run (no commands executed)"
	}
	fmt.Fprintln(out, "Publish plan:")
	fmt.Fprintf(out, "  Mode: %s\n", mode)
	fmt.Fprintf(out, "  Registry: %s\n", plan.Registry)
	fmt.Fprintf(out, "  Account: %s\n", plan.Account)
	fmt.Fprintf(out, "  Package: %s\n", plan.Package)
	fmt.Fprintf(out, "  Credentials: %s\n\n", plan.CredentialSource)
}

func printRegistryError(out io.Writer, registryType, packageName string, err error) {
	var notFound registry.PackageNotFoundError
	if errors.As(err, &notFound) {
		fmt.Fprintf(out, "Error: %s was not found in %s\n", packageName, registryType)
		return
	}
	fmt.Fprintf(out, "Error: unable to query %s in %s: %v\n", packageName, registryType, err)
}

func display(value string) string {
	if strings.TrimSpace(value) == "" {
		return "Not reported"
	}
	return strings.TrimSpace(value)
}

func formatTime(value *time.Time) string {
	if value == nil || value.IsZero() {
		return "Not reported"
	}
	return value.UTC().Format(time.RFC3339)
}

func formatBool(value *bool) string {
	if value == nil {
		return "Not reported"
	}
	return fmt.Sprintf("%t", *value)
}
