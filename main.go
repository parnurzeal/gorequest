package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  _"encoding/json"
  "strings"
)

type Options struct{
  url string
  method string
  body string
  json string
}

func Get(url string) (error, *http.Response, string){
  client := &http.Client{}
  req, err := http.NewRequest("GET", url, nil)
  resp, err:= client.Do(req)
  if err!=nil{
    return err, nil,""
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  return nil, resp, string(body)
}

func CustomRequest(options Options) (error, *http.Response, string){
  fmt.Println(options)
  var (
    req *http.Request
    err error
    resp *http.Response
  )
  client := &http.Client{}
  if options.json != ""{
    content := strings.NewReader(options.json)
    req, err = http.NewRequest(options.method, options.url, content)
    req.Header.Add("Content-Type","application/json")
    fmt.Println("inside", req)
  }
  fmt.Println("outside", req)
  resp, err = client.Do(req)
  if err!=nil{
    return err, nil, ""
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  return nil, resp, string(body)
}

func main(){
  err, response, body:= Get("http://localhost:1337")
  if err==nil && response.StatusCode == 200 {
    fmt.Println(body)
  }
  fmt.Println(err, response, body)
  options:= Options{ url: "http://localhost:1337", method:"GET", body:"hello", json:`{ "hello":"hello"}`}
  CustomRequest(options)

  const jsonStream =`{"sn":"sn1"}`
  reader:=strings.NewReader(jsonStream)
  http.Post("http://localhost:1337", "application/json", reader)

}
