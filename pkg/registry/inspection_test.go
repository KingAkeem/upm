package registry

import (
	"testing"
	"time"
)

func TestInferWarningsLabelsMetadataGapsAndStaleness(t *testing.T) {
	updatedAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	deprecated := true
	inspection := Inspection{
		UpdatedAt:  &updatedAt,
		Deprecated: &deprecated,
	}
	inspection.InferWarnings(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))

	wanted := map[string]bool{
		"missing-license":    false,
		"missing-maintainer": false,
		"missing-repository": false,
		"stale":              false,
		"deprecated":         false,
	}
	for _, warning := range inspection.Warnings {
		if _, ok := wanted[warning.Code]; ok {
			wanted[warning.Code] = true
		}
	}
	for code, found := range wanted {
		if !found {
			t.Errorf("expected warning %q, got %#v", code, inspection.Warnings)
		}
	}
}
