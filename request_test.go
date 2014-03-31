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

/* TODO:Testing post for application/x-www-form-urlencoded
post.query(json), post.query(string), post.send(json), post.send(string), post.query(both).send(both)
*/
func TestPostFormSendString(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		fmt.Println(r.URL.Query())
	}))
	defer ts.Close()
	Post(ts.URL).
		Send("query1=test").
		Send("query2=test").
		End()
}
func TestPostFormSendJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		fmt.Println(r.URL.Query())
	}))
	defer ts.Close()
	Post(ts.URL).
		Send(`{"query1":"test"}`).
		Send(`{"query2":"test"}`).
		End()
}
func TestPostFormSendJsonAndString(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		fmt.Println(r.URL.Query())
	}))
	defer ts.Close()
	Post(ts.URL).
		Send("query1=test").
		Send(`{"query2":"test"}`).
		End()
}

// TODO: check url query (all testcases)
func TestQueryFunc(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		fmt.Println(r.URL.Query())
	}))
	defer ts.Close()
	resp, _ := Post(ts.URL).
		Query("query1=test").
		Query("query2=test").
		End(func(r Response) {
		r.Status = "10"
	})
	fmt.Println(resp.Status)

}

func TestIntegration(t *testing.T) {

}
