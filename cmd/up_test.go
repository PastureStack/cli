package cmd

import (
	"strings"
	"testing"
)

func TestWithEnvironmentReplacesCredentialsExactlyOnce(t *testing.T) {
	environment := withEnvironment([]string{
		"PATH=/usr/bin",
		"PLATFORM_ACCESS_KEY=old-access",
		"platform_secret_key=old-secret",
		"KEEP=value",
	}, map[string]string{
		"PLATFORM_ACCESS_KEY": "new-access",
		"PLATFORM_SECRET_KEY": "new-secret",
	})

	want := map[string]string{
		"PLATFORM_ACCESS_KEY": "new-access",
		"PLATFORM_SECRET_KEY": "new-secret",
	}
	seen := map[string]int{}
	for _, entry := range environment {
		name, value, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		for key, expected := range want {
			if strings.EqualFold(name, key) {
				seen[key]++
				if value != expected {
					t.Fatalf("%s retained stale value %q", key, value)
				}
			}
		}
	}
	for key := range want {
		if seen[key] != 1 {
			t.Fatalf("%s appears %d times, want exactly once", key, seen[key])
		}
	}
}
