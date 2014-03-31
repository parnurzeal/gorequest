package gorequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Response *http.Response

type SuperAgent struct {
	Url       string
	Method    string
	Header    map[string]string
	Type      string
	ForceType string
	Data      map[string]interface{}
	FormData  url.Values
	QueryData url.Values
}

func New() *SuperAgent {
	s := SuperAgent{
		Type:   "json",
		Data:   make(map[string]interface{}),
		Header: make(map[string]string)}
	return &s
}

func Get(targetUrl string) *SuperAgent {
	newReq := &SuperAgent{
		Url:       targetUrl,
		Method:    "GET",
		Header:    make(map[string]string),
		Data:      make(map[string]interface{}),
		FormData:  url.Values{},
		QueryData: url.Values{}}
	return newReq
}

func Post(targetUrl string) *SuperAgent {
	newReq := &SuperAgent{
		Url:       targetUrl,
		Method:    "POST",
		Type:      "json",
		Header:    make(map[string]string),
		Data:      make(map[string]interface{}),
		FormData:  url.Values{},
		QueryData: url.Values{}}
	return newReq
}

func (s *SuperAgent) Set(param string, value string) *SuperAgent {
	s.Header[param] = value
	return s
}

var Types = map[string]string{
	"html":       "text/html",
	"json":       "application/json",
	"xml":        "application/xml",
	"urlencoded": "application/x-www-form-urlencoded",
	"form":       "application/x-www-form-urlencoded",
	"form-data":  "application/x-www-form-urlencoded",
}

func (s *SuperAgent) SetType(typeStr string) *SuperAgent {
	if _, ok := Types[typeStr]; ok {
		s.ForceType = typeStr
	}
	return s
}

// TODO: check error
func (s *SuperAgent) Query(content string) *SuperAgent {
	var val map[string]string
	if err := json.Unmarshal([]byte(content), &val); err == nil {
		for k, v := range val {
			s.QueryData.Add(k, v)
		}
	} else {
		queryVal, _ := url.ParseQuery(content)
		for k, _ := range queryVal {
			s.QueryData.Add(k, queryVal.Get(k))
		}
		// TODO: need to check correct format of 'field=val&field=val&...'
	}
	return s
}

func (s *SuperAgent) Send(content string) *SuperAgent {
	var val map[string]interface{}
	// check if it is json format
	if err := json.Unmarshal([]byte(content), &val); err == nil {
		if s.Type == "form" {
			for k, v := range val {
				// TODO: check if cannot convert to string, return error
				// Also, check that this is the right way to do. (Check superagent)
				s.FormData.Add(k, v.(string))
			}
			// in case previously sending json before knowing it's a form style, we need to include previous added data to formData as well
			for k, v := range s.Data {
				s.FormData.Add(k, v.(string))
			}
			// clear data
			s.Data = nil
		} else {
			s.Type = "json"
			for k, v := range val {
				s.Data[k] = v
			}
		}
	} else {
		// not json format (just normal string)
		s.Type = "form"
		formVal, _ := url.ParseQuery(content)
		for k, _ := range formVal {
			s.FormData.Add(k, formVal.Get(k))
		}
		// change all json data to form style
		for k, v := range s.Data {
			s.FormData.Add(k, v.(string))
		}
		// clear data
		s.Data = make(map[string]interface{})
	}

	return s
}

func (s *SuperAgent) End(callback ...func(response Response)) (Response, error) {
	var (
		req  *http.Request
		err  error
		resp Response
	)
	client := &http.Client{}
	if s.Method == "POST" {
		// if force type, change all data to that type
		if s.ForceType == "json" {
			s.Type = "json"
			// change all json data to form style
			for k, v := range s.FormData {
				s.Data[k] = v
			}
			// clear data
			s.FormData = make(url.Values)
		} else if s.ForceType == "form" {
			s.Type = "form"
			// in case previously sending json before knowing it's a form style, we need to include previous added data to formData as well
			for k, v := range s.Data {
				s.FormData.Add(k, v.(string))
			}
			// clear data
			s.Data = make(map[string]interface{})
		}
		if s.Type == "json" {
			contentJson, _ := json.Marshal(s.Data)
			contentReader := bytes.NewReader(contentJson)
			req, err = http.NewRequest(s.Method, s.Url, contentReader)
			req.Header.Set("Content-Type", "application/json")
		} else if s.Type == "form" {
			req, err = http.NewRequest(s.Method, s.Url, strings.NewReader(s.FormData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	} else if s.Method == "GET" {
		req, err = http.NewRequest(s.Method, s.Url, nil)
	}
	for k, v := range s.Header {
		req.Header.Set(k, v)
	}
	// Add all querystring from Query func
	q := req.URL.Query()
	for k, v := range s.QueryData {
		for _, vv := range v {
			q.Add(k, vv)
		}
	}
	req.URL.RawQuery = q.Encode()
	// Send request
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// deep copy response to give it to both return and callback func
	respCallback := *resp
	if len(callback) != 0 {
		callback[0](&respCallback)
	}
	return resp, nil
}

func main() {
	/*err, response, body:= Get("http://localhost:1337")
	  if err==nil && response.StatusCode == 200 {
	    fmt.Println(body)
	  }
	  fmt.Println(err, response, body)*/

	//s.post("/api/pet").send(`{"name":"tg"}`).end(
	Post("http://requestb.in/1f7ur5s1").
		Send(`nickname=a`).
		Set("Accept", "application/json").
		End(func(response Response) {
		fmt.Println(response)
	})
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
