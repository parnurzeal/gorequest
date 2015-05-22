package gorequest

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/elazarl/goproxy"
)

// testing for Get method
func TestGet(t *testing.T) {
	const case1_empty = "/"
	const case2_set_header = "/set_header"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check method is GET before going to check other features
		if r.Method != GET {
			t.Errorf("Expected method %q; got %q", GET, r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		switch r.URL.Path {
		default:
			t.Errorf("No testing for this case yet : %q", r.URL.Path)
		case case1_empty:
			t.Logf("case %v ", case1_empty)
		case case2_set_header:
			t.Logf("case %v ", case2_set_header)
			if r.Header.Get("API-Key") != "fookey" {
				t.Errorf("Expected 'API-Key' == %q; got %q", "fookey", r.Header.Get("API-Key"))
			}
		}
	}))

	defer ts.Close()

	New().Get(ts.URL + case1_empty).
		End()

	New().Get(ts.URL+case2_set_header).
		Set("API-Key", "fookey").
		End()
}

// testing for POST method
func TestPost(t *testing.T) {
	const case1_empty = "/"
	const case2_set_header = "/set_header"
	const case3_send_json = "/send_json"
	const case4_send_string = "/send_string"
	const case5_integration_send_json_string = "/integration_send_json_string"
	const case6_set_query = "/set_query"
	const case7_integration_send_json_struct = "/integration_send_json_struct"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check method is PATCH before going to check other features
		if r.Method != POST {
			t.Errorf("Expected method %q; got %q", POST, r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		switch r.URL.Path {
		default:
			t.Errorf("No testing for this case yet : %q", r.URL.Path)
		case case1_empty:
			t.Logf("case %v ", case1_empty)
		case case2_set_header:
			t.Logf("case %v ", case2_set_header)
			if r.Header.Get("API-Key") != "fookey" {
				t.Errorf("Expected 'API-Key' == %q; got %q", "fookey", r.Header.Get("API-Key"))
			}
		case case3_send_json:
			t.Logf("case %v ", case3_send_json)
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			if string(body) != `{"query1":"test","query2":"test"}` {
				t.Error(`Expected Body with {"query1":"test","query2":"test"}`, "| but got", string(body))
			}
		case case4_send_string:
			t.Logf("case %v ", case4_send_string)
			if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
				t.Error("Expected Header Content-Type -> application/x-www-form-urlencoded", "| but got", r.Header.Get("Content-Type"))
			}
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			if string(body) != "query1=test&query2=test" {
				t.Error("Expected Body with \"query1=test&query2=test\"", "| but got", string(body))
			}
		case case5_integration_send_json_string:
			t.Logf("case %v ", case5_integration_send_json_string)
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			if string(body) != "query1=test&query2=test" {
				t.Error("Expected Body with \"query1=test&query2=test\"", "| but got", string(body))
			}
		case case6_set_query:
			t.Logf("case %v ", case6_set_query)
			v := r.URL.Query()
			if v["query1"][0] != "test" {
				t.Error("Expected query1:test", "| but got", v["query1"][0])
			}
			if v["query2"][0] != "test" {
				t.Error("Expected query2:test", "| but got", v["query2"][0])
			}
		case case7_integration_send_json_struct:
			t.Logf("case %v ", case7_integration_send_json_struct)
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			comparedBody := []byte(`{"Lower":{"Color":"green","Size":1.7},"Upper":{"Color":"red","Size":0},"a":"a","name":"Cindy"}`)
			if !bytes.Equal(body, comparedBody) {
				t.Errorf(`Expected correct json but got ` + string(body))
			}
		}
	}))

	defer ts.Close()

	New().Post(ts.URL + case1_empty).
		End()

	New().Post(ts.URL+case2_set_header).
		Set("API-Key", "fookey").
		End()

	New().Post(ts.URL + case3_send_json).
		Send(`{"query1":"test"}`).
		Send(`{"query2":"test"}`).
		End()

	New().Post(ts.URL + case4_send_string).
		Send("query1=test").
		Send("query2=test").
		End()

	New().Post(ts.URL + case5_integration_send_json_string).
		Send("query1=test").
		Send(`{"query2":"test"}`).
		End()

	/* TODO: More testing post for application/x-www-form-urlencoded
	   post.query(json), post.query(string), post.send(json), post.send(string), post.query(both).send(both)
	*/
	New().Post(ts.URL + case6_set_query).
		Query("query1=test").
		Query("query2=test").
		End()
	// TODO:
	// 1. test normal struct
	// 2. test 2nd layer nested struct
	// 3. test struct pointer
	// 4. test lowercase won't be export to json
	// 5. test field tag change to json field name
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
	New().Post(ts.URL + case7_integration_send_json_struct).
		Send(`{"a":"a"}`).
		Send(myStyle).
		End()
}

// testing for Patch method
func TestPatch(t *testing.T) {
	const case1_empty = "/"
	const case2_set_header = "/set_header"
	const case3_send_json = "/send_json"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check method is PATCH before going to check other features
		if r.Method != PATCH {
			t.Errorf("Expected method %q; got %q", PATCH, r.Method)
		}
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		switch r.URL.Path {
		default:
			t.Errorf("No testing for this case yet : %q", r.URL.Path)
		case case1_empty:
			t.Logf("case %v ", case1_empty)
		case case2_set_header:
			t.Logf("case %v ", case2_set_header)
			if r.Header.Get("API-Key") != "fookey" {
				t.Errorf("Expected 'API-Key' == %q; got %q", "fookey", r.Header.Get("API-Key"))
			}
		case case3_send_json:
			t.Logf("case %v ", case3_send_json)
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			if string(body) != `{"query1":"test","query2":"test"}` {
				t.Error(`Expected Body with {"query1":"test","query2":"test"}`, "| but got", string(body))
			}
		}
	}))

	defer ts.Close()

	New().Patch(ts.URL + case1_empty).
		End()

	New().Patch(ts.URL+case2_set_header).
		Set("API-Key", "fookey").
		End()

	New().Patch(ts.URL + case3_send_json).
		Send(`{"query1":"test"}`).
		Send(`{"query2":"test"}`).
		End()
}

func checkQuery(t *testing.T, q map[string][]string, key string, want string) {
	v, ok := q[key]
	if !ok {
		t.Error(key, "Not Found")
	} else if len(v) < 1 {
		t.Error("No values for", key)
	} else if v[0] != want {
		t.Errorf("Expected %v:%v | but got %v", key, want, v[0])
	}
	return
}

// TODO: more check on url query (all testcases)
func TestQueryFunc(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header == nil {
			t.Errorf("Expected non-nil request Header")
		}
		v := r.URL.Query()
		checkQuery(t, v, "query1", "test1")
		checkQuery(t, v, "query2", "test2")
	}))
	defer ts.Close()

	New().Post(ts.URL).
		Query("query1=test1").
		Query("query2=test2").
		End()

	qq := struct {
		Query1 string
		Query2 string
	}{
		Query1: "test1",
		Query2: "test2",
	}
	New().Post(ts.URL).
		Query(qq).
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

func TestCookies(t *testing.T) {
	request := New().Timeout(60 * time.Second)
	_, _, errs := request.Get("https://github.com").End()
	if errs != nil {
		t.Errorf("Cookies test request did not complete")
		return
	}
	domain, _ := url.Parse("https://github.com")
	if len(request.Client.Jar.Cookies(domain)) == 0 {
		t.Errorf("Expected cookies | but get nothing")
	}
}

func TestGetSetCookies(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != GET {
			t.Errorf("Expected method %q; got %q", GET, r.Method)
		}
		c, err := r.Cookie("API-Cookie-Name")
		if err != nil {
			t.Error(err)
		}
		if c == nil {
			t.Errorf("Expected non-nil request Cookie 'API-Cookie-Name'")
		} else if c.Value != "api-cookie-value" {
			t.Errorf("Expected 'API-Cookie-Name' == %q; got %q", "api-cookie-value", c.Value)
		}
	}))
	defer ts.Close()

	New().Get(ts.URL).
		AddCookie(&http.Cookie{Name: "API-Cookie-Name", Value: "api-cookie-value"}).
		End()
}

func TestErrorTypeWrongKey(t *testing.T) {
	//defer afterTest(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, checkTypeWrongKey")
	}))
	defer ts.Close()

	_, _, err := New().
		Get(ts.URL).
		Type("wrongtype").
		End()
	if len(err) != 0 {
		if err[0].Error() != "Type func: incorrect type \"wrongtype\"" {
			t.Errorf("Wrong error message: " + err[0].Error())
		}
	} else {
		t.Errorf("Should have error")
	}
}

func TestBasicAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			t.Error("bad syntax")
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		if pair[0] != "myusername" || pair[1] != "mypassword" {
			t.Error("Wrong username/password")
		}
	}))
	defer ts.Close()
	New().Post(ts.URL).
		SetBasicAuth("myusername", "mypassword").
		End()
}
