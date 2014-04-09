// Package gorequest inspired by Nodejs SuperAgent provides easy-way to write http client
package gorequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Response *http.Response

type SuperAgent struct {
	Url        string
	Method     string
	Header     map[string]string
	TargetType string
	ForceType  string
	Data       map[string]interface{}
	FormData   url.Values
	QueryData  url.Values
	Client     *http.Client
}

func New() *SuperAgent {
	s := &SuperAgent{
		TargetType: "json",
		Data:       make(map[string]interface{}),
		Header:     make(map[string]string),
		FormData:   url.Values{},
		QueryData:  url.Values{},
		Client:     &http.Client{},
	}
	return s
}
func (s *SuperAgent) ClearSuperAgent() {
	s.Url = ""
	s.Method = ""
	s.Header = make(map[string]string)
	s.Data = make(map[string]interface{})
	s.FormData = url.Values{}
	s.QueryData = url.Values{}
	s.ForceType = ""
	s.TargetType = "json"
}

func (s *SuperAgent) Get(targetUrl string) *SuperAgent {
	s.ClearSuperAgent()
	s.Method = "GET"
	s.Url = targetUrl
	return s
}

func (s *SuperAgent) Post(targetUrl string) *SuperAgent {
	s.ClearSuperAgent()
	s.Method = "POST"
	s.Url = targetUrl
	return s
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

func (s *SuperAgent) Type(typeStr string) *SuperAgent {
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

// TODO: find a way to change it to gorequest's Request and Response itself
func (s *SuperAgent) RedirectPolicy(policy func(req *http.Request, via []*http.Request) error) *SuperAgent {
	s.Client.CheckRedirect = policy
	return s
}

func (s *SuperAgent) Send(content string) *SuperAgent {
	var val map[string]interface{}
	// check if it is json format
	if err := json.Unmarshal([]byte(content), &val); err == nil {
		for k, v := range val {
			s.Data[k] = v
		}
	} else {
		formVal, _ := url.ParseQuery(content)
		for k, _ := range formVal {
			// make it array if already have key
			if val, ok := s.Data[k]; ok {
				var strArray []string
				strArray = append(strArray, formVal.Get(k))
				// check if previous data is one string or array
				switch oldValue := val.(type) {
				case []string:
					strArray = append(strArray, oldValue...)
				case string:
					strArray = append(strArray, oldValue)
				}
				s.Data[k] = strArray
			} else {
				// make it just string if does not already have same key
				s.Data[k] = formVal.Get(k)
			}
		}
		s.TargetType = "form"
	}
	return s
}

func changeMapToURLValues(data map[string]interface{}) url.Values {
	var newUrlValues = url.Values{}
	for k, v := range data {
		switch val := v.(type) {
		case string:
			newUrlValues.Add(k, string(val))
		case []string:
			for _, element := range val {
				newUrlValues.Add(k, element)
			}
		}
	}
	return newUrlValues
}

func (s *SuperAgent) End(callback ...func(response Response, body string)) (Response, string, error) {
	var (
		req  *http.Request
		err  error
		resp Response
	)
	// check if there is forced type
	if s.ForceType == "json" {
		s.TargetType = "json"
	} else if s.ForceType == "form" {
		s.TargetType = "form"
	}
	if s.Method == "POST" {
		if s.TargetType == "json" {
			contentJson, _ := json.Marshal(s.Data)
			contentReader := bytes.NewReader(contentJson)
			req, err = http.NewRequest(s.Method, s.Url, contentReader)
			req.Header.Set("Content-Type", "application/json")
		} else if s.TargetType == "form" {
			formData := changeMapToURLValues(s.Data)
			req, err = http.NewRequest(s.Method, s.Url, strings.NewReader(formData.Encode()))
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
	fmt.Println(req.URL)
	resp, err = s.Client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	bodyCallback := body
	// deep copy response to give it to both return and callback func
	respCallback := *resp
	if len(callback) != 0 {
		callback[0](&respCallback, string(bodyCallback))
	}
	return resp, string(body), nil
}

func main() {
	New().Post("http://requestb.in/1f7ur5s1").
		Send(`nickname=a`).
		Set("Accept", "application/json").
		End(func(response Response, body string) {
		fmt.Println(response)
	})
}
