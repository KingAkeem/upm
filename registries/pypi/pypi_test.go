package pypi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestInspectMapsVerifiedMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/pypi/fixture-package/json" {
			t.Fatalf("unexpected request path: %s", request.URL.Path)
		}
		response.Header().Set("Content-Type", "application/json")
		fmt.Fprint(response, `{
			"info":{
				"name":"fixture-package",
				"version":"1.4.0",
				"license":"Apache-2.0",
				"author":"Primary Author",
				"maintainer":"Package Maintainer",
				"package_url":"https://pypi.org/project/fixture-package/",
				"project_urls":{"Source":"https://github.com/example/fixture-package"}
			},
			"releases":{
				"1.3.0":[{"upload_time_iso_8601":"2025-01-02T03:04:05Z","yanked":false}],
				"1.4.0":[{"upload_time_iso_8601":"2026-02-03T04:05:06Z","yanked":true,"yanked_reason":"broken wheel"}]
			}
		}`)
	}))
	defer server.Close()

	value := NewPypiRegistry(map[string]interface{}{"url": server.URL, "client": server.Client()})
	inspection, err := value.Inspect("fixture-package", time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if inspection.LatestVersion != "1.4.0" || inspection.License != "Apache-2.0" {
		t.Fatalf("unexpected metadata: %#v", inspection)
	}
	if inspection.RepositoryURL != "https://github.com/example/fixture-package" {
		t.Fatalf("unexpected repository URL: %q", inspection.RepositoryURL)
	}
	if inspection.Yanked == nil || !*inspection.Yanked {
		t.Fatal("expected yanked metadata")
	}
	if len(inspection.Warnings) != 1 || inspection.Warnings[0].Code != "yanked" {
		t.Fatalf("expected one inferred yanked warning, got %v", inspection.Warnings)
	}
}

func TestInspectDoesNotMarkMixedFilesYanked(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprint(response, `{
			"info":{"name":"fixture-package","version":"1.0.0"},
			"releases":{"1.0.0":[
				{"upload_time_iso_8601":"2026-01-01T00:00:00Z","yanked":true},
				{"upload_time_iso_8601":"2026-01-01T00:00:01Z","yanked":false}
			]}
		}`)
	}))
	defer server.Close()

	value := NewPypiRegistry(map[string]interface{}{"url": server.URL, "client": server.Client()})
	inspection, err := value.Inspect("fixture-package", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if inspection.Yanked == nil || *inspection.Yanked {
		t.Fatal("release should only be yanked when every current file is yanked")
	}
}
