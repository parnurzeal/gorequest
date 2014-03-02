package gorequest

import (
  "fmt"
  "net/http"
  "io/ioutil"
  _"encoding/json"
  "strings"
)

type Options struct{
  Url string
  Method string
  Body string
  Json string
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
  if options.Method == "POST" {
    if options.Json != ""{
      content := strings.NewReader(options.Json)
      req, err = http.NewRequest(options.Method, options.Url, content)
      req.Header.Set("Content-Type","application/json")
    }
  }
  resp, err = client.Do(req)
  if err!=nil{
    return err, nil, ""
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  return nil, resp, string(body)
}


func main(){
  /*err, response, body:= Get("http://localhost:1337")
  if err==nil && response.StatusCode == 200 {
    fmt.Println(body)
  }
  fmt.Println(err, response, body)*/
  options:= Options{ Url: "http://localhost:1337", Method:"POST", Body:"hello", Json:`{ "hello":"hello"}`}
  CustomRequest(options)


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
