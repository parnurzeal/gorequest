Note: This is a work in progress and not totally support all specifications.
Right now, you can only do simple get and post with easy to specify header like in examples

GoRequest
=========

GoRequest -- Simplified HTTP client ( inspired by famous SuperAgent lib in Node.js )

## Installation

```
$ go get github.com/parnurzeal/gorequest
```

## Why should you use GoRequest?

GoRequest makes thing much more simple for you, making http client more awesome and fun like SuperAgent + golang style usage.

This is what you normally do for a simple GET without GoRequest:

```
resp, err := http.Get("http://example.com/")
```

With GoRequest:

```
resp, body, err := gorequest.Get("http://example.com/").End()
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
resp, body, err := gorequest.Get("http://example.com").
  Set("If-None-Match", `W/"wyzzy"`).
  End()
```

For __a JSON POST__, you might need to marshal map data to json format, setting header to 'application/json' and declare http.Client with standard library. So, you code might look like:

```
m := map[string]interface{}{
  "value1": "one",
  "value2": "two",
}
mJson, _ := json.Marshal(m)
contentReader := bytes.NewReader(mJson)
req, _ := http.NewRequest("POST", "http://requestb.in/1kovd3s1", contentReader)
req.Header.Set("Content-Type", "application/json")
client := &http.Client{}
resp, _ := client.Do(req)
```



