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
	msg := struct {
		Msg string `json:"msg"`
	}{
		Msg: fmt.Sprintf("Hello %s!\n", person.Name),
	}

	fdk.WriteStatus(out, 200)
	json.NewEncoder(out).Encode(&msg)

}
