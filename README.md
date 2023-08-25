# UPM - Universal Package Manager for Multi-Language Repositories

Simplify your development workflow with UPM, a versatile package manager designed for multi-language repositories. Whether you're coding in Python, JavaScript, Java, Golang, or any other language, our tool streamlines the process of retrieving and publishing packages across diverse ecosystems. Say goodbye to language-specific complexities and embrace a unified solution for managing dependencies. Effortlessly fetch, version, and share packages, enhancing collaboration and boosting productivity. Join us in revolutionizing cross-language and cross-platform package management today!

## Features

- **Universal Solution:** Seamlessly manage packages across various programming languages within a unified tool.
- **Effortless Retrieval:** Fetch metadata for a single package effortlessly using UPM's command-line interface.
- **Streamlined Publishing:** Effortlessly publish your packages using simple commands.
- **Secure Credentials:** Securely manage authentication using encrypted local files.

## Supported Ecosystems

- npm
- pypi

## Getting Started

1. Install UPM by running: `go install github.com/KingAkeem/upm`.
2. Retrieve package metadata: `upm fetch -n {package_name}`.
3. Publish a package: `upm publish -u {username} -p {password}`.
4. To use credentials from a local file, run: `upm publish -c path/to/credentials.json`.

## Example

### Fetching
```bash
$ upm fetch -n vue
Creating registry for npm.
Success: Registry created for npm.
Performing 'fetch' for package 'vue'.
Success: 'vue' has been found within npm.
------------------------
Description: The progressive JavaScript framework for building modern web UI.
License: MIT
Version: 0.11.0-rc3
Registry: npm
Name: vue
Author: Evan You
------------------------
Creating registry for pypi.
Success: Registry created for pypi.
Performing 'fetch' for package 'vue'.
Success: 'vue' has been found within pypi.
------------------------
Author: titan - i@qtitan.com
Description: 
License: 
Version: 0.0.1
Registry: pypi
Name: vue
------------------------
```

### Publishing

### pypi
```bash
$ upm publish -u ${username} -p ${password} 
```

### npm
```bash
$ npm login
$ upm publish 
```
## Contributing
Contributions are welcome! To contribute:

1. Fork the KingAkeem/upm repository.
2. Create a branch: git checkout -b feature/your-feature-name.
3. Commit your changes: git commit -am 'Add new feature'.
4. Push the branch: git push origin feature/your-feature-name.
5. Open a pull request.
