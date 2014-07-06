package gorequest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/elazarl/goproxy"
)

func TestGetFormat(t *testing.T) {
	//defer afterTest(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != GET {
			t.Errorf("Expected method %q; got %q", GET, r.Method)
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
		if r.Method != GET {
			t.Errorf("Expected method %q; got %q", GET, r.Method)
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
	//defer afterTest(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != POST {
			t.Errorf("Expected method %q; got %q", POST, r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
	}))
	defer ts.Close()

	New().Post(ts.URL).
		End()
}

func TestPostSetHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != POST {
			t.Errorf("Expected method %q; got %q", POST, r.Method)
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

/* TODO: More testing post for application/x-www-form-urlencoded
post.query(json), post.query(string), post.send(json), post.send(string), post.query(both).send(both)
*/
func TestPostFormSendString(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Error("Expected Header Content-Type -> application/x-www-form-urlencoded", "| but got", r.Header.Get("Content-Type"))
		}
		body, _ := ioutil.ReadAll(r.Body)
		if string(body) != "query1=test&query2=test" {
			t.Error("Expected Body with \"query1=test&query2=test\"", "| but got", string(body))
		}
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
		body, _ := ioutil.ReadAll(r.Body)
		if string(body) != `{"query1":"test","query2":"test"}` {
			t.Error(`Expected Body with {"query1":"test","query2":"test"}`, "| but got", string(body))
		}
	}))
	defer ts.Close()
	New().Post(ts.URL).
		Send(`{"query1":"test"}`).
		Send(`{"query2":"test"}`).
		End()
}

func TestPostFormSendJsonSend_ShouldGive_FormBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		body, _ := ioutil.ReadAll(r.Body)
		if string(body) != "query1=test&query2=test" {
			t.Error("Expected Body with \"query1=test&query2=test\"", "| but got", string(body))
		}
	}))
	defer ts.Close()
	New().Post(ts.URL).
		Send("query1=test").
		Send(`{"query2":"test"}`).
		End()
}

// TODO: more check on url query (all testcases)
func TestQueryFunc(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		v := r.URL.Query()
		if v["query1"][0] != "test" {
			t.Error("Expected query1:test", "| but got", v["query1"][0])
		}
		if v["query2"][0] != "test" {
			t.Error("Expected query2:test", "| but got", v["query2"][0])
		}
	}))
	defer ts.Close()
	New().Post(ts.URL).
		Query("query1=test").
		Query("query2=test").
		End()
}

// TODO: more tests on redirect
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

// 1. test normal struct
// 2. test 2nd layer nested struct
// 3. test struct pointer
// 4. test lowercase won't be export to json
// 5. test field tag change to json field name
func TestSendStructFunc(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		comparedBody := []byte(`{"Lower":{"Color":"green","Size":1.7},"Upper":{"Color":"red","Size":0},"a":"a","name":"Cindy"}`)
		if !bytes.Equal(body, comparedBody) {
			t.Errorf(`Expected correct json but got ` + string(body))
		}
	}))
	defer ts.Close()
	type Upper struct {
		Color string
		Size  int
		note  string
	}
	type Lower struct {
		Color string
		Size  float64
		note  string
	}

	type Style struct {
		Upper Upper
		Lower Lower
		Name  string `json:"name"`
	}
	myStyle := Style{Upper: Upper{Color: "red"}, Name: "Cindy", Lower: Lower{Color: "green", Size: 1.7}}
	New().Post(ts.URL).
		Send(`{"a":"a"}`).
		Send(myStyle).
		End()
}

func TestTimeoutFunc(t *testing.T) {
	// 1st case, dial timeout
	startTime := time.Now()
	_, _, errs := New().Timeout(1000 * time.Millisecond).Get("http://www.google.com:81").End()
	elapsedTime := time.Since(startTime)
	if errs == nil {
		t.Errorf("Expected dial timeout error but get nothing")
	}
	if elapsedTime < 1000*time.Millisecond || elapsedTime > 1500*time.Millisecond {
		t.Errorf("Expected timeout in between 1000 -> 1500 ms | but got ", elapsedTime)
	}
	// 2st case, read/write timeout (Can dial to url but want to timeout because too long operation on the server)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(200)
	}))
	request := New().Timeout(1000 * time.Millisecond)
	startTime = time.Now()
	_, _, errs = request.Get(ts.URL).End()
	elapsedTime = time.Since(startTime)
	if errs == nil {
		t.Errorf("Expected dial+read/write timeout | but get nothing")
	}
	if elapsedTime < 1000*time.Millisecond || elapsedTime > 1500*time.Millisecond {
		t.Errorf("Expected timeout in between 1000 -> 1500 ms | but got ", elapsedTime)
	}
	// 3rd case, testing reuse of same request
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(200)
	}))
	startTime = time.Now()
	_, _, errs = request.Get(ts.URL).End()
	elapsedTime = time.Since(startTime)
	if errs == nil {
		t.Errorf("Expected dial+read/write timeout | but get nothing")
	}
	if elapsedTime < 1000*time.Millisecond || elapsedTime > 1500*time.Millisecond {
		t.Errorf("Expected timeout in between 1000 -> 1500 ms | but got ", elapsedTime)
	}

}

// TODO: complete integration test
func TestIntegration(t *testing.T) {

}
