package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fnproject/fdk-go"
	"os"
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
	body := fmt.Sprintf("Hello %s!\n", person.Name)
	err := json.NewEncoder(out).Encode(body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
