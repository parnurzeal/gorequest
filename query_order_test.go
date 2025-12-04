package gorequest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueryParameterOrder(t *testing.T) {
	var receivedRawQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedRawQuery = r.URL.RawQuery
		w.WriteHeader(200)
	}))
	defer ts.Close()

	// Test that parameters are sent in the order they were added
	New().Get(ts.URL).
		Param("b", "1").
		Param("a", "2").
		End()

	// The issue is that currently this would be "a=2&b=1" instead of "b=1&a=2"
	expected := "b=1&a=2"
	if receivedRawQuery != expected {
		t.Errorf("Expected query string %q, got %q", expected, receivedRawQuery)
	}
}

func TestQueryParameterOrderWithURLParams(t *testing.T) {
	var receivedRawQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedRawQuery = r.URL.RawQuery
		w.WriteHeader(200)
	}))
	defer ts.Close()

	// Test that URL query parameters come first, then added parameters
	New().Get(ts.URL+"?z=0").
		Param("b", "1").
		Param("a", "2").
		End()

	// URL params first, then added params in order
	expected := "z=0&b=1&a=2"
	if receivedRawQuery != expected {
		t.Errorf("Expected query string %q, got %q", expected, receivedRawQuery)
	}
}

func TestQueryParameterOrderWithQuery(t *testing.T) {
	var receivedRawQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedRawQuery = r.URL.RawQuery
		w.WriteHeader(200)
	}))
	defer ts.Close()

	// Test with Query method using string
	New().Get(ts.URL).
		Query("b=1").
		Query("a=2").
		End()

	expected := "b=1&a=2"
	if receivedRawQuery != expected {
		t.Errorf("Expected query string %q, got %q", expected, receivedRawQuery)
	}
}

func TestQueryParameterOrderWithSpecialChars(t *testing.T) {
	var receivedRawQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedRawQuery = r.URL.RawQuery
		w.WriteHeader(200)
	}))
	defer ts.Close()

	// Test that special characters are properly encoded
	New().Get(ts.URL).
		Param("key with space", "value with space").
		Param("key&special", "value=special").
		End()

	// Special characters should be encoded
	expected := "key+with+space=value+with+space&key%26special=value%3Dspecial"
	if receivedRawQuery != expected {
		t.Errorf("Expected query string %q, got %q", expected, receivedRawQuery)
	}
}
