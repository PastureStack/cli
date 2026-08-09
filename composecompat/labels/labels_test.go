package labels

import "testing"

func TestComposeHashLabelContract(t *testing.T) {
	if got, want := HASH.Str(), "com.docker.compose.config-hash"; got != want {
		t.Fatalf("hash label = %q, want %q", got, want)
	}
}
