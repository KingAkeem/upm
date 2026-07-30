package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPublishDefaultsToSafePreview(t *testing.T) {
	t.Setenv("NPM_USERNAME", "fixture-account")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := run([]string{"publish", "-n", "example-package", "-r", "npm"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected success, got %d: %s", code, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"dry-run (no commands executed)",
		"Registry: npm",
		"Account: fixture-account",
		"Package: example-package",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected %q in output:\n%s", expected, output)
		}
	}
}

func TestPlaintextPasswordFlagIsRejected(t *testing.T) {
	const credential = "do-not-echo-this-value"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := run([]string{"publish", "-n", "example-package", "-r", "pypi", "-p", credential}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected usage failure, got %d", code)
	}
	if strings.Contains(stdout.String()+stderr.String(), credential) {
		t.Fatal("secret was written to command output")
	}
}

func TestExecuteRequiresOneExplicitRegistry(t *testing.T) {
	_, err := parseOptions([]string{"publish", "-n", "example-package", "-r", "npm,pypi", "--execute"}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "exactly one registry") {
		t.Fatalf("expected explicit single-registry error, got %v", err)
	}
}
