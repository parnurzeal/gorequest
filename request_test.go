package gorequest

import (
	"fmt"
	"github.com/elazarl/goproxy"
	_ "io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

	New().Get(ts.URL).
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

	New().Get(ts.URL).
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

	New().Post(ts.URL).
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
		//fmt.Println(r.URL.Query())
	}))
	defer ts.Close()
	New().Post(ts.URL).
		Send("query1=test").
		Send("query2=test").
		End()
}
func TestPostFormSendJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		//fmt.Println(r.URL.Query())
	}))
	defer ts.Close()
	New().Post(ts.URL).
		Send(`{"query1":"test"}`).
		Send(`{"query2":"test"}`).
		End()
}
func TestPostFormSendJsonAndString(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		//fmt.Println(r.URL.Query())
	}))
	defer ts.Close()
	New().Post(ts.URL).
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
		//fmt.Println(r.URL.Query())
	}))
	defer ts.Close()
	New().Post(ts.URL).
		Query("query1=test").
		Query("query2=test").
		End(func(r Response, body string, errs []error) {
		r.Status = "10"
	})
	//fmt.Println(resp.Status)

}

// TODO: check redirect
func TestRedirectPolicyFunc(t *testing.T) {
	redirectSuccess := false
	redirectFuncGetCalled := false
	tsRedirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectSuccess = true
	}))
	defer tsRedirect.Close()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, tsRedirect.URL, http.StatusMovedPermanently)
	}))
	defer ts.Close()

	New().
		Get(ts.URL).
		RedirectPolicy(func(req Request, via []Request) error {
		redirectFuncGetCalled = true
		return nil
	}).End()
	if !redirectSuccess {
		t.Errorf("Expected reaching another redirect url not original one")
	}
	if !redirectFuncGetCalled {
		t.Errorf("Expected redirect policy func to get called")
	}
}

func TestProxyFunc(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "proxy passed")
	}))
	defer ts.Close()
	// start proxy
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			return r, nil
		})
	ts2 := httptest.NewServer(proxy)
	// sending request via Proxy
	resp, body, _ := New().Proxy(ts2.URL).Get(ts.URL).End()
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200 Status code")
	}
	if body != "proxy passed" {
		t.Errorf("Expected 'proxy passed' body string")
	}
}

// TODO: added check for the correct timeout error string
// Right now, I see 2 different errors from timeout. Need to check why there are two of them. (i/o timeout and operation timed out)
func TestTimeoutFunc(t *testing.T) {
	_, _, errs := New().Timeout(1000 * time.Millisecond).Get("http://www.google.com:81").End()
	if errs == nil {
		t.Errorf("Expected timeout error but get nothing")
	}
}

func TestIntegration(t *testing.T) {

}
