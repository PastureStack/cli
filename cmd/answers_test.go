package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseStructuredAnswers(t *testing.T) {
	file := filepath.Join(t.TempDir(), "answers.yaml")
	if err := os.WriteFile(file, []byte("name: pasture\nenabled: true\nreplicas: 3\noptional: null\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	values, err := parseStructuredAnswers(file, ".yaml")
	if err != nil {
		t.Fatal(err)
	}
	if values["name"] != "pasture" || values["enabled"] != "true" || values["replicas"] != "3" || values["optional"] != "" {
		t.Fatalf("unexpected values: %#v", values)
	}
}

func TestParseEnvironmentAnswers(t *testing.T) {
	file := filepath.Join(t.TempDir(), "answers.env")
	if err := os.WriteFile(file, []byte("# comment\nexport NAME=pasture\nEMPTY=\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	values, err := parseEnvironmentAnswers(file)
	if err != nil {
		t.Fatal(err)
	}
	if values["NAME"] != "pasture" || values["EMPTY"] != "" || len(values) != 2 {
		t.Fatalf("unexpected values: %#v", values)
	}
}
