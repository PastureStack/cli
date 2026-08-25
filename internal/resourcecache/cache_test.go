package resourcecache

import (
	"testing"
	"time"
)

func TestSetGetAndExpire(t *testing.T) {
	cache := New(5 * time.Millisecond)
	cache.Set("resource:1", "value")
	if value, ok := cache.Get("resource:1"); !ok || value != "value" {
		t.Fatalf("unexpected cache read: value=%v ok=%v", value, ok)
	}
	time.Sleep(10 * time.Millisecond)
	if value, ok := cache.Get("resource:1"); ok || value != nil {
		t.Fatalf("expired value remained: value=%v ok=%v", value, ok)
	}
}
