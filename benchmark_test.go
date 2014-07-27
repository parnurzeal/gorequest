package gorequest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkPostSendJson(b *testing.B) {
	request := New()
	for n := 0; n < b.N; n++ {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))
		request.Post(ts.URL).
			Send(`{"query1":"test"}`).
			Send(`{"query2":"test"}`).
			End()
		ts.Close()
	}

}
