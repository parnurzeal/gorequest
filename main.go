package gorequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Options struct {
	Url    string
	Method string
	Body   string
	Json   string
	Header string
}

type SuperAgent struct {
	Url    string
	Method string
	Header map[string]string
	Type   string
	Data   map[string]interface{}
}

func New() *SuperAgent {
	s := SuperAgent{
		Type:   "json",
		Data:   make(map[string]interface{}),
		Header: make(map[string]string)}
	return &s
}

func Post(url string) *SuperAgent {
	newReq := &SuperAgent{
		Url:    url,
		Method: "POST",
		Type:   "json",
		Header: make(map[string]string),
		Data:   make(map[string]interface{})}
	return newReq
}

func (s *SuperAgent) Set(param string, value string) *SuperAgent {
	s.Header[param] = value
	return s
}

func (s *SuperAgent) Send(content string) *SuperAgent {
	if s.Type == "json" {
		var val map[string]interface{}
		if err := json.Unmarshal([]byte(content), &val); err != nil {
			fmt.Println("ERROR to json.Unmarshal in send", err)
		}
		for k, v := range val {
			s.Data[k] = v
		}
	}
	return s
}

func (s *SuperAgent) End() (error, *http.Response, string) {
	var (
		req  *http.Request
		err  error
		resp *http.Response
	)
	client := &http.Client{}
	if s.Method == "POST" {
		if s.Type == "json" {
			contentJson, _ := json.Marshal(s.Data)
			contentReader := bytes.NewReader(contentJson)
			req, err = http.NewRequest(s.Method, s.Url, contentReader)
			req.Header.Set("Content-Type", "application/json")
		}
	}
	for k, v := range s.Header {
		req.Header.Set(k, v)
	}
	resp, err = client.Do(req)
	if err != nil {
		return err, nil, ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return nil, resp, string(body)
}

func Get(url string) (error, *http.Response, string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return err, nil, ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return nil, resp, string(body)
}

func CustomRequest(options Options) (error, *http.Response, string) {
	fmt.Println(options)
	var (
		req  *http.Request
		err  error
		resp *http.Response
	)
	client := &http.Client{}
	if options.Method == "POST" {
		if options.Json != "" {
			content := strings.NewReader(options.Json)
			req, err = http.NewRequest(options.Method, options.Url, content)
			req.Header.Set("Content-Type", "application/json")
		}
	}
	resp, err = client.Do(req)
	if err != nil {
		return err, nil, ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return nil, resp, string(body)
}

func main() {
	/*err, response, body:= Get("http://localhost:1337")
	  if err==nil && response.StatusCode == 200 {
	    fmt.Println(body)
	  }
	  fmt.Println(err, response, body)*/
	options := Options{Url: "http://localhost:1337", Method: "POST", Body: "hello", Json: `{ "hello":"hello"}`}
	CustomRequest(options)

	//s.post("/api/pet").send(`{"name":"tg"}`).end()
	Post("http://requestb.in/1f7ur5s1").
		Send(`{"nickname":"a"}`).
		Set("Accept", "application/json").
		End()
	/*client:= &http.Client{}
	  req,_ := http.NewRequest("GET", "http://localhost:1337", nil)
	  req.Header.Add("Content-Type","application/json")
	  fmt.Println("main",req)
	  res, _ :=  client.Do(req)
	  fmt.Println(res.Body)
	  /*const jsonStream =`{"sn":"sn1"}`
	  reader:=strings.NewReader(jsonStream)
	  resp,_ := http.Post("http://localhost:1337", "application/json", reader)
	  defer resp.Body.Close()
	  body,_ :=ioutil.ReadAll(resp.Body)
	  fmt.Println(resp)
	  fmt.Println(string(body))*/
}
