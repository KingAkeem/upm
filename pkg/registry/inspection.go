package registry

import (
	"strings"
	"time"
)

const staleAfter = 2 * 365 * 24 * time.Hour

type Warning struct {
	Code    string
	Message string
}

type Inspection struct {
	Registry      string
	Name          string
	LatestVersion string
	License       string
	Maintainers   []string
	PublishedAt   *time.Time
	UpdatedAt     *time.Time
	Deprecated    *bool
	Yanked        *bool
	PackageURL    string
	RepositoryURL string
	Warnings      []Warning
}

type PublishPlan struct {
	Registry         string
	Package          string
	Account          string
	CredentialSource string
}

func (inspection *Inspection) InferWarnings(now time.Time) {
	warnings := make([]Warning, 0, 6)
	if strings.TrimSpace(inspection.License) == "" {
		warnings = append(warnings, Warning{Code: "missing-license", Message: "registry metadata does not declare a license"})
	}
	if len(inspection.Maintainers) == 0 {
		warnings = append(warnings, Warning{Code: "missing-maintainer", Message: "registry metadata does not identify an author or maintainer"})
	}
	if strings.TrimSpace(inspection.RepositoryURL) == "" {
		warnings = append(warnings, Warning{Code: "missing-repository", Message: "registry metadata does not link to a source repository"})
	}
	if inspection.UpdatedAt == nil {
		warnings = append(warnings, Warning{Code: "missing-update-time", Message: "registry metadata does not provide a usable update timestamp"})
	} else if now.After(*inspection.UpdatedAt) && now.Sub(*inspection.UpdatedAt) > staleAfter {
		warnings = append(warnings, Warning{Code: "stale", Message: "package has not been updated in more than two years"})
	}
	if inspection.Deprecated != nil && *inspection.Deprecated {
		warnings = append(warnings, Warning{Code: "deprecated", Message: "the registry marks the latest package version as deprecated"})
	}
	if inspection.Yanked != nil && *inspection.Yanked {
		warnings = append(warnings, Warning{Code: "yanked", Message: "all files for the latest package version are yanked"})
	}
	inspection.Warnings = warnings
}
