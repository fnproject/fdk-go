package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	var person struct {
		Name string `json:"name"`
	}
	json.NewDecoder(in).Decode(&person)

	if person.Name == "" {
		person.Name = "world"
	}

	out.Write([]byte(fmt.Sprintf("Hello %s!\n", person.Name)))
}
