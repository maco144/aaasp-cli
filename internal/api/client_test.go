package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPIError_IncludesMessageWhenPresent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error":"not_runnable","message":"This deployment has no anthropic key"}`))
	}))
	defer srv.Close()

	err := New(srv.URL, "k").Post("/runs/stream", map[string]any{}, nil)

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("want *APIError, got %v", err)
	}
	if apiErr.Status != 422 {
		t.Errorf("status = %d, want 422", apiErr.Status)
	}
	if !strings.Contains(apiErr.Message, "not_runnable") || !strings.Contains(apiErr.Message, "no anthropic key") {
		t.Errorf("message = %q, want both the code and the human message", apiErr.Message)
	}
}

func TestAPIError_CodeOnly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not_found"}`))
	}))
	defer srv.Close()

	err := New(srv.URL, "k").Get("/deployments/x", nil)

	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Message != "not_found" {
		t.Fatalf("want message not_found, got %v", err)
	}
}
