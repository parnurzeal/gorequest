// Package gorequest inspired by Nodejs SuperAgent provides easy-way to write http client
package gorequest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Request *http.Request
type Response *http.Response

// A SuperAgent is a object storing all request data for client.
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
	Errors     []error
}

// Used to create a new SuperAgent object.
func New() *SuperAgent {
	s := &SuperAgent{
		TargetType: "json",
		Data:       make(map[string]interface{}),
		Header:     make(map[string]string),
		FormData:   url.Values{},
		QueryData:  url.Values{},
		Client:     &http.Client{},
		Errors:     make([]error, 0),
	}
	return s
}

// Clear SuperAgent data for another new request.
func (s *SuperAgent) ClearSuperAgent() {
	s.Url = ""
	s.Method = ""
	s.Header = make(map[string]string)
	s.Data = make(map[string]interface{})
	s.FormData = url.Values{}
	s.QueryData = url.Values{}
	s.ForceType = ""
	s.TargetType = "json"
	s.Errors = make([]error, 0)
}

func (s *SuperAgent) Get(targetUrl string) *SuperAgent {
	s.ClearSuperAgent()
	s.Method = "GET"
	s.Url = targetUrl
	s.Errors = make([]error, 0)
	return s
}

func (s *SuperAgent) Post(targetUrl string) *SuperAgent {
	s.ClearSuperAgent()
	s.Method = "POST"
	s.Url = targetUrl
	s.Errors = make([]error, 0)
	return s
}

// Set is used for setting header fields.
// Example. To set `Accept` as `application/json`
//
//    gorequest.New().
//      Post("/gamelist").
//      Set("Accept", "application/json").
//      End()
// TODO: make Set be able to get multiple fields
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

// Type is a convenience function to specify the data type to send.
// For example, to send data as `application/x-www-form-urlencoded` :
//
//    gorequest.New().
//      Post("/recipe").
//      Type("form").
//      Send(`{ name: "egg benedict", category: "brunch" }`).
//      End()
//
// This will POST the body "name=egg benedict&category=brunch" to url /recipe
//
// GoRequest supports
//
//    "text/html" uses "html"
//    "application/json" uses "json"
//    "application/xml" uses "xml"
//    "application/x-www-form-urlencoded" uses "urlencoded", "form" or "form-data"
//
func (s *SuperAgent) Type(typeStr string) *SuperAgent {
	if _, ok := Types[typeStr]; ok {
		s.ForceType = typeStr
	} else {
		s.Errors = append(s.Errors, errors.New("Type func: incorrect type \""+typeStr+"\""))
	}
	return s
}

// Query function accepts either json string or strings which will form a query-string in url of GET method or body of POST method.
// For example, making "/search?query=bicycle&size=50x50&weight=20kg" using GET method:
//
//      gorequest.New().
//        Get("/search").
//        Query(`{ query: 'bicycle' }`).
//        Query(`{ size: '50x50' }`).
//        Query(`{ weight: '20kg' }`).
//        End()
//
// Or you can put multiple json values:
//
//      gorequest.New().
//        Get("/search").
//        Query(`{ query: 'bicycle', size: '50x50', weight: '20kg' }`).
//        End()
//
// Strings are also acceptable:
//
//      gorequest.New().
//        Get("/search").
//        Query("query=bicycle&size=50x50").
//        Query("weight=20kg").
//        End()
//
// Or even Mixed! :)
//
//      gorequest.New().
//        Get("/search").
//        Query("query=bicycle").
//        Query(`{ size: '50x50', weight:'20kg' }`).
//        End()
//
// TODO: check error
func (s *SuperAgent) Query(content string) *SuperAgent {
	var val map[string]string
	if err := json.Unmarshal([]byte(content), &val); err == nil {
		for k, v := range val {
			s.QueryData.Add(k, v)
		}
	} else {
		if queryVal, err := url.ParseQuery(content); err == nil {
			for k, _ := range queryVal {
				s.QueryData.Add(k, queryVal.Get(k))
			}
		} else {
			s.Errors = append(s.Errors, err)
		}
		// TODO: need to check correct format of 'field=val&field=val&...'
	}
	return s
}

func (s *SuperAgent) RedirectPolicy(policy func(req Request, via []Request) error) *SuperAgent {
	s.Client.CheckRedirect = func(r *http.Request, v []*http.Request) error {
		vv := make([]Request, len(v))
		for i, r := range v {
			vv[i] = Request(r)
		}
		return policy(Request(r), vv)
	}
	return s
}

// Send function accepts either json string or query strings which is usually used to assign data to POST or PUT method.
// Without specifying any type, if you give Send with json data, you are doing requesting in json format:
//
//      gorequest.New().
//        Post("/search").
//        Send(`{ query: 'sushi' }`).
//        End()
//
// While if you use at least one of querystring, GoRequest understands and automatically set the Content-Type to `application/x-www-form-urlencoded`
//
//      gorequest.New().
//        Post("/search").
//        Send("query=tonkatsu").
//        End()
//
// So, if you want to strictly send json format, you need to use Type func to set it as `json` (Please see more details in Type function).
// You can also do multiple chain of Send:
//
//      gorequest.New().
//        Post("/search").
//        Send("query=bicycle&size=50x50").
//        Send(`{ wheel: '4'}`).
//        End()
//
// TODO: check error
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

func (s *SuperAgent) End(callback ...func(response Response, body string)) (Response, string, []error) {
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
		s.Errors = append(s.Errors, err)
		return nil, "", s.Errors
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
