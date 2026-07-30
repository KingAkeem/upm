package npm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestInspectMapsVerifiedMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/fixture-package" {
			t.Fatalf("unexpected request path: %s", request.URL.Path)
		}
		response.Header().Set("Content-Type", "application/json")
		fmt.Fprint(response, `{
			"name":"fixture-package",
			"license":"MIT",
			"author":{"name":"Primary Author"},
			"maintainers":[{"name":"Package Maintainer"}],
			"repository":{"type":"git","url":"git+https://github.com/example/fixture-package.git"},
			"dist-tags":{"latest":"2.0.0"},
			"versions":{"2.0.0":{"name":"fixture-package","version":"2.0.0","deprecated":"use the replacement"}},
			"time":{"2.0.0":"2026-01-02T03:04:05Z","modified":"2026-02-03T04:05:06Z"}
		}`)
	}))
	defer server.Close()

	value := NewNpmRegistry(map[string]interface{}{"url": server.URL, "client": server.Client()})
	inspection, err := value.Inspect("fixture-package", time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if inspection.LatestVersion != "2.0.0" || inspection.License != "MIT" {
		t.Fatalf("unexpected metadata: %#v", inspection)
	}
	if inspection.RepositoryURL != "https://github.com/example/fixture-package" {
		t.Fatalf("unexpected repository URL: %q", inspection.RepositoryURL)
	}
	if inspection.Deprecated == nil || !*inspection.Deprecated {
		t.Fatal("expected deprecated metadata")
	}
	if len(inspection.Maintainers) != 2 {
		t.Fatalf("expected author and maintainer, got %v", inspection.Maintainers)
	}
	if len(inspection.Warnings) != 1 || inspection.Warnings[0].Code != "deprecated" {
		t.Fatalf("expected one inferred deprecation warning, got %v", inspection.Warnings)
	}
}

func TestFetchReturnsContextualNotFound(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	value := NewNpmRegistry(map[string]interface{}{"url": server.URL, "client": server.Client()})
	_, err := value.Fetch("missing")
	if err == nil || err.Error() != "package missing was not found in npm" {
		t.Fatalf("unexpected error: %v", err)
	}
}
