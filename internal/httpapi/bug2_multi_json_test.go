package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRejectMultipleJSONObjects verifies that the API rejects request bodies
// containing more than one JSON object. Without this check, only the first
// object is decoded and the trailing data is silently ignored, which can
// mask client-side serialization bugs or injection attempts.
func TestRejectMultipleJSONObjects(t *testing.T) {
	api := New()
	srv := httptest.NewServer(api.Handler())
	defer srv.Close()

	body := `{"expr":"* * * * *"}{"expr":"0 0 * * *"}`
	resp, err := http.Post(srv.URL+"/api/validate", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for multi-object body", resp.StatusCode)
	}
}

// TestRejectTrailingGarbage verifies rejection of valid JSON followed by
// garbage bytes that would be silently ignored without the extra-decode check.
func TestRejectTrailingGarbage(t *testing.T) {
	api := New()
	srv := httptest.NewServer(api.Handler())
	defer srv.Close()

	body := `{"expr":"* * * * *"}trailing`
	resp, err := http.Post(srv.URL+"/api/validate", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for trailing garbage", resp.StatusCode)
	}
}
