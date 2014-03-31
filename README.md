GoRequest
=========

GoRequest -- Simplified HTTP client ( inspired by famous SuperAgent lib in Node.js )

## Installation

```
$ go get github.com/parnurzeal/gorequest
```

## Documentation
See [Go Doc](http://godoc.org/github.com/parnurzeal/gorequest) for usage and details.

## Status

[![Drone Build Status](https://drone.io/github.com/jmcvetta/restclient/status.png)](https://drone.io/github.com/parnurzeal/gorequest/latest)
[![Travis Build Status](https://travis-ci.org/parnurzeal/gorequest.svg?branch=master)](https://travis-ci.org/parnurzeal/gorequest)

## Why should you use GoRequest?

GoRequest makes thing much more simple for you, making http client more awesome and fun like SuperAgent + golang style usage.

This is what you normally do for a simple GET without GoRequest:

```
resp, err := http.Get("http://example.com/")
```

With GoRequest:

```
resp, err := gorequest.Get("http://example.com/").End()
```

How about getting control over HTTP client headers, redirect policy, and etc. Things is getting more complicated in golang. You need to create a Client, setting header in different comamnd, ... to do just only one __GET__

```
client := &http.Client{
  CheckRedirect: redirectPolicyFunc,
}

req, err := http.NewRequest("GET", "http://example.com", nil)

req.Header.Add("If-None-Match", `W/"wyzzy"`)
resp, err := client.Do(req)
```

Why making things ugly while you can just do as follows:

```
### policy is not supported yet
resp, err := gorequest.Get("http://example.com").
  Set("If-None-Match", `W/"wyzzy"`).
  End()
```

For a __JSON POST__ with standard libraries, you might need to marshal map data structure to json format, setting header to 'application/json' (and other headers if you need to) and declare http.Client. So, you code become longer and hard to maintain:

```
m := map[string]interface{}{
  "name": "backy",
  "species": "dog",
}
mJson, _ := json.Marshal(m)
contentReader := bytes.NewReader(mJson)
req, _ := http.NewRequest("POST", "http://example.com", contentReader)
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Notes","GoRequest is coming!")
client := &http.Client{}
resp, _ := client.Do(req)
```

Compared to our GoRequest version, JSON is for sure a default. So, it turns out to be just one simple line!:

```
resp, err := gorequest.Post("http://example.com").
  Set("Notes","gorequst is coming!").
  Send(`{"name":"backy", "species":"dog"}`).
  End()
```

Note: This is a work in progress and not totally support all specifications.
Right now, you can only do simple get and post with easy to specify header like in examples which is enough in many cases.

## License

GoRequest is MIT License.


