package testutil

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func AssertStatus(t *testing.T, w *httptest.ResponseRecorder, want int) {
	t.Helper()
	if w.Code != want {
		t.Errorf("status=%d, want=%d — body: %s", w.Code, want, w.Body.String())
	}
}

func AssertJSONField(t *testing.T, w *httptest.ResponseRecorder, key string, want any) {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %s", w.Body.String())
	}
	got, ok := body[key]
	if !ok {
		t.Errorf("field %q missing in response: %s", key, w.Body.String())
		return
	}
	if got != want {
		t.Errorf("field %q = %v, want %v", key, got, want)
	}
}

func ParseJSON(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON object: %s", w.Body.String())
	}
	return body
}

func ParseJSONArray(t *testing.T, w *httptest.ResponseRecorder) []any {
	t.Helper()
	var body []any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON array: %s", w.Body.String())
	}
	return body
}
