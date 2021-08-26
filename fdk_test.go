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

package fdk

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// XXX(reed): test cloudevents in http-stream land

// echoHandler echos the body and all headers back
func echoHandler(ctx context.Context, in io.Reader, out io.Writer) {
	for k, vs := range GetContext(ctx).Header() {
		for _, v := range vs {
			AddHeader(out, k, v)
		}
	}

	// XXX(reed): could configure this to test too
	WriteStatus(out, http.StatusTeapot+2)
	io.Copy(out, in)
}

// echoContextHandler echos the tracing context back
func echoContextHandler(ctx context.Context, in io.Reader, out io.Writer) {
	nctx := GetContext(ctx)

	if resp, ok := out.(http.ResponseWriter); ok {
		resp.Header().Add("traceId", nctx.TracingContextData().TraceId())
		resp.Header().Add("spanId", nctx.TracingContextData().SpanId())
		resp.Header().Add("parentSpanId", nctx.TracingContextData().ParentSpanId())
		resp.Header().Add("sampled", strconv.FormatBool(nctx.TracingContextData().IsSampled()))
		resp.Header().Add("ociURL", nctx.TracingContextData().TraceCollectorURL())
		resp.Header().Add("flag", nctx.TracingContextData().Flags())
		resp.Header().Add("ServiceName", nctx.TracingContextData().ServiceName())
	}
	// XXX(reed): could configure this to test too

	WriteStatus(out, http.StatusTeapot+2)
	io.Copy(out, in)
}

// echoContextHandler echos the context app data back
func echoContextHandlerAppData(ctx context.Context, in io.Reader, out io.Writer) {
	nctx := GetContext(ctx)

	if resp, ok := out.(http.ResponseWriter); ok {
		resp.Header().Add("AppId", nctx.AppID())
		resp.Header().Add("AppName", nctx.AppName())
		resp.Header().Add("FnId", nctx.FnID())
		resp.Header().Add("FnName", nctx.FnName())
		resp.Header().Add("CallId", nctx.CallID())
	}
	// XXX(reed): could configure this to test too

	WriteStatus(out, http.StatusTeapot+2)
	io.Copy(out, in)
}

func TestHTTPStreamSock(t *testing.T) {
	// XXX(reed): move to fdk_linux_test.go with build tag
	// XXX(reed): extract the underlying server handler / write tests against it instead of starting uds for other tests

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tmpSock, err := ioutil.TempDir("/tmp", "fdk-go-test")
	if err != nil {
		t.Fatal("couldn't make tmpdir for testing")
	}
	defer os.RemoveAll(tmpSock)

	tmpSock = filepath.Join(tmpSock, "fn.sock")

	go startHTTPServer(ctx, HandlerFunc(echoHandler), "unix:"+tmpSock)

	// let the uds server start... could inotify but don't want the dependency for tests...
	time.Sleep(1 * time.Second)

	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
			IdleConnTimeout:     1 * time.Second,
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, "unix", tmpSock)
			},
		},
	}

	// TODO headers?
	bodyString := "yodawg"
	req, err := http.NewRequest("POST", "http://localhost/call", strings.NewReader(bodyString))
	if err != nil {
		t.Fatal("error making req", err)
	}
	req.Header.Set("yo", "dawg")

	res, err := client.Do(req)
	if err != nil {
		t.Fatal("error doing uds request", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatal("got wrong status code:", res.StatusCode)
	}

	if res.Header.Get("Fn-Fdk-Version") != versionHeader {
		t.Errorf("Expected \"%s\" but got \"%s\"", versionHeader, res.Header.Get("Fn-Fdk-Version"))
	}

	if res.Header.Get("Fn-Fdk-Runtime") != runtimeHeader {
		t.Errorf("Expected \"%s\" but got \"%s\"", runtimeHeader, res.Header.Get("Fn-Fdk-Runtime"))
	}

	outBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if res.Header.Get("yo") != "dawg" {
		t.Fatal("expected yo dawg header, didn't get it", res.Header)
	}

	if string(outBody) != bodyString {
		t.Fatal("body mismatch:", string(outBody), bodyString)
	}
}

func TestHandler(t *testing.T) {
	tests := []struct {
		name      string
		inBody    string
		inHeader  http.Header
		outBody   string
		outHeader http.Header
	}{
		{"invoke", "yodawg", http.Header{"Yo": {"dawg"}}, "yodawg", http.Header{"Yo": {"dawg"}}},
		{"httpgw", "yodawg", http.Header{"Fn-Intent": {"httprequest"}, "Fn-Http-H-Yo": {"dawg"}}, "yodawg", http.Header{"Fn-Http-H-Yo": {"dawg"}, "Fn-Http-Status": {"420"}}},
		{"httpgw-rm-nongw", "yodawg", http.Header{"Fn-Intent": {"httprequest"}, "Yo": {"dawg"}}, "yodawg", http.Header{"Fn-Http-Status": {"420"}}},
		// TODO(reed): test Fn-Http-Request-Url, Fn-Http-Method, Fn-Call-Id, Fn-Deadline...
	}

	// TODO make it so echoHandler takes expected headers to test
	handler := &httpHandler{HandlerFunc(echoHandler)}

	for _, test := range tests {
		req, err := http.NewRequest("POST", "http://localhost/call", strings.NewReader(test.inBody))
		if err != nil {
			t.Fatal("error making req", err)
		}
		req.Header = test.inHeader

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		resp := w.Result()

		if w.Body.String() != test.outBody {
			t.Error("body mismatch", test.name, w.Body.String(), test.outBody)
		}

		for k := range test.outHeader {
			if resp.Header.Get(k) != test.outHeader.Get(k) {
				t.Error("header mismatch", test.name, k, resp.Header.Get(k), test.outHeader.Get(k))
			}
		}
	}
}

func TestTracingContextInHandlerWithOCITracingDisabled(t *testing.T) {
	tests := []struct {
		name      string
		inBody    string
		inHeader  http.Header
		outBody   string
		outHeader http.Header
	}{
		{"testinvoke", "yodawg", http.Header{"Fn-Intent": []string{"httprequest"}, "Fn-Call-Id": []string{"dawg"}}, "yodawg",
			http.Header{"Fn-Http-Status": {"420"}, "Fn-Http-H-Ociurl": {""}, "Fn-Http-H-Parentspanid": {""},
				"Fn-Http-H-Sampled": {"false"}, "Fn-Http-H-Spanid": {""}, "Fn-Http-H-Traceid": {""},
				"Fn-Http-H-Flag": {""}}},
	}

	os.Setenv("FN_APP_ID", "GOLANGAPP1234323dsw")
	os.Setenv("FN_FN_ID", "12342r321422b41")
	os.Setenv("FN_APP_NAME", "GOLANGAPP")
	os.Setenv("FN_FN_NAME", "FNTESTAPP")
	os.Setenv("OCI_TRACING_ENABLED", "0")

	handler := &httpHandler{HandlerFunc(echoContextHandler)}

	for _, test := range tests {
		req, err := http.NewRequest("POST", "http://localhost/invoke", strings.NewReader(test.inBody))
		if err != nil {
			t.Fatal("error making req", err)
		}
		req.Header = test.inHeader
		req.Header.Set("Fn-Call-Id", "fncall")

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		resp := w.Result()

		if w.Body.String() != test.outBody {
			t.Error("body mismatch", test.name, w.Body.String(), test.outBody)
		}

		for k := range test.outHeader {
			if resp.Header.Get(k) != test.outHeader.Get(k) {
				t.Error("header mismatch", test.name, k, resp.Header.Get(k), test.outHeader.Get(k))
			}
		}
	}
}

func TestTracingContextInHandler(t *testing.T) {
	tests := []struct {
		name      string
		inBody    string
		inHeader  http.Header
		outBody   string
		outHeader http.Header
	}{
		{"testinvoke", "yodawg", http.Header{"Fn-Intent": []string{"httprequest"}, "Fn-Call-Id": []string{"dawg"}}, "yodawg",
			http.Header{"Fn-Http-Status": {"420"}, "Fn-Http-H-Ociurl": {"localhost"}, "Fn-Http-H-Parentspanid": {""},
				"Fn-Http-H-Sampled": {"true"}, "Fn-Http-H-Spanid": {"12345675685432"}, "Fn-Http-H-Traceid": {"12345675685432"},
				"Fn-Http-H-Flag": {"1"}, "Fn-Http-H-Servicename": {"golangapp::fntestapp"}}},
	}

	os.Setenv("FN_APP_ID", "GOLANGAPP1234323dsw")
	os.Setenv("FN_FN_ID", "12342r321422b41")
	os.Setenv("FN_APP_NAME", "GOLANGAPP")
	os.Setenv("FN_FN_NAME", "FNTESTAPP")
	os.Setenv("OCI_TRACING_ENABLED", "1")
	os.Setenv("OCI_TRACE_COLLECTOR_URL", "localhost")

	handler := &httpHandler{HandlerFunc(echoContextHandler)}

	for _, test := range tests {
		req, err := http.NewRequest("POST", "http://localhost/invoke", strings.NewReader(test.inBody))
		if err != nil {
			t.Fatal("error making req", err)
		}
		req.Header = test.inHeader
		req.Header.Set("x-b3-traceid", "12345675685432")
		req.Header.Set("x-b3-spanid", "12345675685432")
		req.Header.Set("x-b3-parentspanid", "")
		req.Header.Set("x-b3-flags", "1")
		req.Header.Set("x-b3-sampled", "1")
		req.Header.Set("Fn-Call-Id", "fncall")

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		resp := w.Result()

		if w.Body.String() != test.outBody {
			t.Error("body mismatch", test.name, w.Body.String(), test.outBody)
		}

		for k := range test.outHeader {
			if resp.Header.Get(k) != test.outHeader.Get(k) {
				t.Error("header mismatch", test.name, k, resp.Header.Get(k), test.outHeader.Get(k))
			}
		}
	}
}

func TestContextInHandlerWithAppFunctionData(t *testing.T) {
	tests := []struct {
		name      string
		inBody    string
		inHeader  http.Header
		outBody   string
		outHeader http.Header
	}{
		{"testinvoke", "yodawg", http.Header{"Fn-Intent": []string{"httprequest"}, "Fn-Call-Id": []string{"dawg"}}, "yodawg",
			http.Header{"Fn-Http-Status": {"420"}, "Fn-Http-H-Appid": {"GOLANGAPP1234323dsw"}, "Fn-Http-H-Appname": {"GOLANGAPP"},
				"Fn-Http-H-Callid": {"fncall"}, "Fn-Http-H-Fnid": {"12342r321422b41"}, "Fn-Http-H-Fnname": {"FNTESTAPP"}}},
	}

	os.Setenv("FN_APP_ID", "GOLANGAPP1234323dsw")
	os.Setenv("FN_FN_ID", "12342r321422b41")
	os.Setenv("FN_APP_NAME", "GOLANGAPP")
	os.Setenv("FN_FN_NAME", "FNTESTAPP")
	os.Setenv("OCI_TRACING_ENABLED", "0")

	handler := &httpHandler{HandlerFunc(echoContextHandlerAppData)}

	for _, test := range tests {
		req, err := http.NewRequest("POST", "http://localhost/invoke", strings.NewReader(test.inBody))
		if err != nil {
			t.Fatal("error making req", err)
		}
		req.Header = test.inHeader
		req.Header.Set("Fn-Call-Id", "fncall")

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		resp := w.Result()

		if w.Body.String() != test.outBody {
			t.Error("body mismatch", test.name, w.Body.String(), test.outBody)
		}

		for k := range test.outHeader {
			if resp.Header.Get(k) != test.outHeader.Get(k) {
				t.Error("header mismatch", test.name, k, resp.Header.Get(k), test.outHeader.Get(k))
			}
		}
	}
}

// NOTE: the below may serve as a reminder that memory allocs suck and you can do better

const mappers = 10

func memory1nMap(m map[int]int) map[int]int {
	rm := make(map[int]int, len(m))
	for _, i := range m {
		rm[mappers-i] = i
	}
	return rm
}

func memory2nMap(m map[int]int) map[int]int {
	for _, i := range m {
		_ = i
		continue
	}

	for _, i := range m {
		m[mappers-i] = i
	}
	return m
}

func BenchmarkMapCrap1(b *testing.B) {
	m := make(map[int]int)
	for i := 0; i < mappers; i++ {
		m[i] = i
	}

	for i := 0; i < b.N; i++ {
		memory1nMap(m)
	}
}

func BenchmarkMapCrap2(b *testing.B) {
	m := make(map[int]int)
	for i := 0; i < mappers; i++ {
		m[i] = i
	}

	for i := 0; i < b.N; i++ {
		memory2nMap(m)
	}
}
