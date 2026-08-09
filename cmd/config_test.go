package cmd

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestGetSecretConfigKeepsExistingValueWithoutDisplayingIt(t *testing.T) {
	const secret = "existing-secret-value"
	var output bytes.Buffer

	value, err := getSecretConfig(bufio.NewReader(strings.NewReader("\n")), &output, "Secret Key", secret)
	if err != nil {
		t.Fatal(err)
	}
	if value != secret {
		t.Fatalf("expected the existing secret to be retained, got %q", value)
	}
	if strings.Contains(output.String(), secret) {
		t.Fatal("secret value was written to the prompt")
	}
	if got, want := output.String(), "Secret Key [configured]: "; got != want {
		t.Fatalf("unexpected prompt: got %q, want %q", got, want)
	}
}

func TestGetSecretConfigAcceptsReplacementWithoutDisplayingIt(t *testing.T) {
	const replacement = "replacement-secret-value"
	var output bytes.Buffer

	value, err := getSecretConfig(bufio.NewReader(strings.NewReader(replacement+"\n")), &output, "Secret Key", "")
	if err != nil {
		t.Fatal(err)
	}
	if value != replacement {
		t.Fatalf("unexpected replacement: got %q, want %q", value, replacement)
	}
	if strings.Contains(output.String(), replacement) {
		t.Fatal("replacement secret was written to the prompt")
	}
	if got, want := output.String(), "Secret Key: "; got != want {
		t.Fatalf("unexpected prompt: got %q, want %q", got, want)
	}
}
