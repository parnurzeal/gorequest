package gorequest

import (
	"fmt"
	_ "io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
