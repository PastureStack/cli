package catalog

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestCatalogHTTPClientBlocksCrossOriginRedirect(t *testing.T) {
	var targetReached atomic.Bool
	target := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		targetReached.Store(true)
		_, _ = response.Write([]byte(`{}`))
	}))
	defer target.Close()

	source := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, target.URL, http.StatusFound)
	}))
	defer source.Close()

	client := &RancherBaseClientImpl{Opts: &ClientOpts{AccessKey: "access", SecretKey: "secret", Timeout: time.Second}}
	err := client.doGet(source.URL, nil, &map[string]interface{}{})
	if err == nil || !strings.Contains(err.Error(), "cross-origin redirect blocked") {
		t.Fatalf("unexpected redirect result: %v", err)
	}
	if targetReached.Load() {
		t.Fatal("cross-origin redirect reached the target")
	}
}

func TestCatalogSafeURLForLogRemovesSecrets(t *testing.T) {
	value := safeURLForLog("https://user:password@example.test/path?token=secret#fragment")
	for _, secret := range []string{"password", "token", "secret", "fragment"} {
		if strings.Contains(value, secret) {
			t.Fatalf("safe URL still contains %q: %s", secret, value)
		}
	}
}
