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
