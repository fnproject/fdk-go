# Go Fn Development Kit (FDK)

fdk-go provides convenience functions for writing Go fn code

For getting started with fn, please refer to https://github.com/fnproject/fn/blob/master/README.md

# Example function using fdk-go

```go
package main

import (
  "bytes"
  "fmt"
  "io"

  fdk "github.com/fnproject/fdk-go"
)

func main() {
  fdk.Do(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
  fnctx := fdk.Context(ctx)

  var b bytes.Buffer
  io.Copy(&b, in)
  fmt.Fprintf(out, fmt.Sprintf("Hello %s\n", b.String()))

  for k, vs := range fnctx.Header {
    fmt.Fprintf(out, fmt.Sprintf("ENV: %s %#v\n", k, vs))
  }
  return nil
}
```

# Advanced example

```go
package main

import (
  fdk "github.com/fnproject/fdk-go"
)

func main() {
  fdk.Do(new(mydoer))
}

// use a struct to maintain state between invocations
type mydoer struct {
  count int
}

func (m *mydoer) Serve(ctx context.Context, in io.Reader, out io.Writer) {
  m.count++

  fnctx := fdk.Context(ctx)

  contentType := fntctx.Header.Get("Content-Type")
  if contentType != "application/json" {
    fdk.WriteStatus(out, 400)
    fdk.SetHeader(out, "Content-Type", "application/json")
    io.Copy(out, `{"error":"invalid content type"}`)
    return
  }

  if fnctx.Config["FN_METHOD"] != "PUT" {
    fdk.WriteStatus(out, 404)
    fdk.SetHeader(out, "Content-Type", "application/json")
    io.Copy(out, `{"error":"route not found"}`)
    return
  }

  var person struct {
    Name `json:"name"`
  }
  json.NewDecoder(in).Decode(&person)

  // you can write your own headers & status, if you'd like to
  fdk.WriteStatus(out, 201)
  fdk.SetHeader(out, "Content-Type", "application/json")

  all := struct {
    Name   string              `json:"name"`
    Header map[string][]string `json:"header"`
    Config map[string]string   `json:"config"`
  }{
    Name: person.Name,
    Header: fnctx.Header,
    Config: fnctx.Config,
  }

  json.NewEncoder(out).Encode(&all)
}
```
