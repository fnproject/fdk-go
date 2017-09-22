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
  fdk.Do(myHandler)
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) error {
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
  fdk.Do(myHandler)
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) error {
  fnctx := fdk.Context(ctx)

  contentType := fntctx.Headers["Content-Type"]
  if contentType != "application/json" {
    return fdk.Error(400, "invalid content type")
  }

  var person struct {
    Name `json:"name"`
  }
  json.NewDecoder(in).Decode(&person)

  // you can write your own headers, if you'd like to
  fdk.WriteStatus(out, 201)
  fdk.WriteHeader(out, "Content-Type", "application/json")

  all := struct {
    Name   string `json:"name"`
    Header map[string][]string `json:"header"`
    Config map[string]string `json:"config"`
  }{
    Name: person.Name,
    Header: fnctx.Header,
    Config: fnctx.Config,
  }

  return json.NewEncoder(out).Encode(&all)
}
```
