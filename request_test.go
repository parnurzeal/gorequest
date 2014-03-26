package gorequest

import (
	"fmt"
	_ "io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var robotsTxtHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Last-Modified", "sometime")
	fmt.Fprintf(w, "User-agent: go\nDisallow: /something/")
})

func TestGetFormat(t *testing.T) {
	//defer afterTest(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected method %q; got %q", "GET", r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
	}))
	defer ts.Close()

	Get(ts.URL).
		End()
}

func TestGetSetHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected method %q; got %q", "GET", r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		if r.Header.Get("API-Key") != "fookey" {
			t.Errorf("Expected 'API-Key' == %q; got %q", "fookey", r.Header.Get("API-Key"))
		}
	}))
	defer ts.Close()

	Get(ts.URL).
		Set("API-Key", "fookey").
		End()
}

func TestPostFormat(t *testing.T) {

}

func TestPostSetHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected method %q; got %q", "POST", r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		if r.Header.Get("API-Key") != "fookey" {
			t.Errorf("Expected 'API-Key' == %q; got %q", "fookey", r.Header.Get("API-Key"))
		}
	}))
	defer ts.Close()

	Post(ts.URL).
		Set("API-Key", "fookey").
		End()
}
