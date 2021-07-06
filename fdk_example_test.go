/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fdk_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	fdk "github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func Example() { println("use main()") }

// TODO make http.Handler example

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	fnctx, ok := fdk.GetContext(ctx).(fdk.HTTPContext)
	if !ok {
		// optionally, this may be a good idea
		fdk.WriteStatus(out, 400)
		fdk.SetHeader(out, "Content-Type", "application/json")
		io.WriteString(out, `{"error":"function not invoked via http trigger"}`)
		return
	}

	contentType := fnctx.Header().Get("Content-Type")
	if contentType != "application/json" {
		// can assert content type for your api this way
		fdk.WriteStatus(out, 400)
		fdk.SetHeader(out, "Content-Type", "application/json")
		io.WriteString(out, `{"error":"invalid content type"}`)
		return
	}

	if fnctx.RequestMethod() != "PUT" {
		// can assert certain request methods for certain endpoints
		fdk.WriteStatus(out, 404)
		fdk.SetHeader(out, "Content-Type", "application/json")
		io.WriteString(out, `{"error":"route not found, method not supported"}`)
		return
	}

	var person struct {
		Name string `json:"name"`
	}
	json.NewDecoder(in).Decode(&person)

	// this is where you might insert person into a database or do something else

	all := struct {
		Name   string            `json:"name"`
		URL    string            `json:"url"`
		Header http.Header       `json:"header"`
		Config map[string]string `json:"config"`
	}{
		Name:   person.Name,
		URL:    fnctx.RequestURL(),
		Header: fnctx.Header(),
		Config: fnctx.Config(),
	}

	// you can write your own headers & status, if you'd like to
	fdk.SetHeader(out, "Content-Type", "application/json")
	fdk.WriteStatus(out, 201)
	json.NewEncoder(out).Encode(&all)
}
