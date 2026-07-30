# UPM

UPM provides a small, common command-line interface for inspecting, fetching, and publishing packages across npm and PyPI.

## Supported registries

- npm
- PyPI

## Install

```bash
git clone git@github.com:KingAkeem/upm.git
cd upm
go build -o upm ./cmd/main
```

## Inspect package trust metadata

`inspect` reports registry-provided fields separately from UPM's inferred warnings. A warning is a review prompt, not verified evidence that a package is unsafe.

```bash
upm inspect -n requests -r pypi
upm inspect -n vue -r npm
upm inspect -n example-package -r all
```

Verified metadata includes the latest version, license, maintainers or authors, release and update times, deprecated or yanked status when the registry exposes it, and package and repository URLs.

Inferred warnings cover:

- missing license, maintainer, repository, or update metadata
- packages whose last reported update is more than two years old
- registry-reported deprecation or yanking

Always review the verified fields and linked upstream sources before making a dependency decision.

## Fetch package metadata

```bash
upm fetch -n vue -r npm
upm fetch -n requests -r pypi
```

Omit `-r` to query both supported registries.

## Publish safely

Publishing is a dry run by default. The preview identifies the registry, account label, package, and credential source without executing a build or upload command.

```bash
upm publish -n example-package -r npm
upm publish -n example-package -r pypi --dry-run
```

After reviewing the plan, add `--execute` to publish to exactly one registry:

```bash
upm publish -n example-package -r npm --execute
upm publish -n example-package -r pypi --execute
```

The package argument is used for the safety preview. Publishing commands run from the current package directory.

### Credentials

UPM does not accept password command-line flags and does not print credential values.

For npm, authenticate with `npm login` or configure an npm access token through npm's supported configuration or environment mechanism. For PyPI, configure Twine using an environment or CI secret, `.pypirc`, or a system keyring. Token-based PyPI authentication should use `__token__` as `TWINE_USERNAME`; provide `TWINE_PASSWORD` through a protected environment rather than shell history.

Use `--account` only to label the non-secret identity shown in a preview. It is not used as a credential.

## Development

```bash
go test ./...
go vet ./...
go build ./cmd/main
```

## Contributing

Contributions are welcome. Create a focused branch, add tests for behavioral changes, and open a pull request against this repository.
