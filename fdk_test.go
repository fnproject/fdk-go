package fdk

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"testing"
)

func echoHTTPHandler(ctx context.Context, in io.Reader, out io.Writer) {
	io.Copy(out, in)
	WriteStatus(out, http.StatusTeapot+2)
	SetHeader(out, "yo", "dawg")
}

func TestHandler(t *testing.T) {
	inString := "yodawg"
	var in bytes.Buffer
	io.WriteString(&in, inString)

	var out bytes.Buffer
	echoHTTPHandler(buildCtx(), &in, &out)

	if out.String() != inString {
		t.Fatalf("this was supposed to be easy. strings no matchy: %s got: %s", inString, out.String())
	}
}

func TestDefault(t *testing.T) {
	inString := "yodawg"
	var in bytes.Buffer
	io.WriteString(&in, inString)

	var out bytes.Buffer

	doDefault(HandlerFunc(echoHTTPHandler), buildCtx(), &in, &out)

	if out.String() != inString {
		t.Fatalf("strings no matchy: %s got: %s", inString, out.String())
	}
}

func JSONHandler(ctx context.Context, in io.Reader, out io.Writer) {
	var person struct {
		Name string `json:"name"`
	}
	json.NewDecoder(in).Decode(&person)

	if person.Name == "" {
		person.Name = "world"
	}

	out.Write([]byte(fmt.Sprintf("Hello %s!\n", person.Name)))
}

type testJSONOut struct {
	jsonio
	Protocol *callResponseHTTP `json:"protocol,omitempty"`
}

func TestJSON(t *testing.T) {
	req := jsonIn{
		jsonio{
			Body:        `{"name":"john"}`,
			ContentType: "application/json",
		},
		"someid",
		&callRequestHTTP{
			Type:       "json",
			RequestURL: "someURL",
			Headers:    http.Header{},
		},
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal("Unable to marshal request")
	}

	in := bytes.NewReader(b)
	var out bytes.Buffer

	doJSONOnce(HandlerFunc(JSONHandler), buildCtx(), in, &out, &bytes.Buffer{}, make(http.Header))
	output := &testJSONOut{}
	err = json.Unmarshal(out.Bytes(), output)
	if err != nil {
		t.Fatalf("Unable to unmarshal response, err: %v", err)
	}
	if !strings.Contains(output.Body, "Hello john!") {
		t.Fatalf("Output assertion mismatch. Expected: `Hello john!\n`. Actual: %v", output.Body)
	}
	if output.Protocol.StatusCode != 200 {
		t.Fatalf("Response code must equal to 200, but have: %v", output.Protocol.StatusCode)
	}
}

func TestFailedJSON(t *testing.T) {
	dummyBody := "should fail with this"
	in := strings.NewReader(dummyBody)
	var out bytes.Buffer
	doJSONOnce(HandlerFunc(JSONHandler), buildCtx(), in, &out, &bytes.Buffer{}, make(http.Header))
	output := &testJSONOut{}
	json.Unmarshal(out.Bytes(), output)
	if output.Protocol.StatusCode != 500 {
		t.Fatalf("Response code must equal to 500, but have: %v", output.Protocol.StatusCode)
	}
}

func TestHTTP(t *testing.T) {
	// simulate fn writing us http requests...

	bodyString := "yodawg"
	in := HTTPreq(t, bodyString)

	var out bytes.Buffer
	ctx := buildCtx()
	doHTTPOnce(HandlerFunc(echoHTTPHandler), ctx, in, &out, &bytes.Buffer{}, make(http.Header))

	res, err := http.ReadResponse(bufio.NewReader(&out), nil)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusTeapot+2 {
		t.Fatal("got wrong status code", res.StatusCode)
	}

	outBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if res.Header.Get("yo") != "dawg" {
		t.Fatal("expected yo dawg header, didn't get it")
	}

	if string(outBody) != bodyString {
		t.Fatal("strings no matchy for http", string(outBody), bodyString)
	}
}

func HTTPreq(t *testing.T, bod string) io.Reader {
	req, err := http.NewRequest("GET", "http://localhost:8080/r/myapp/yodawg", strings.NewReader(bod))
	if err != nil {
		t.Fatal(err)
	}

	byts, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(byts)
}
