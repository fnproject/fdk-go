package fdk

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// XXX(reed): test cloudevents in http-stream land

func echoHTTPHandler(_ context.Context, in io.Reader, out io.Writer) {
	io.Copy(out, in)
	WriteStatus(out, http.StatusTeapot+2)
	SetHeader(out, "yo", "dawg")
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

	go startHTTPServer(ctx, HandlerFunc(echoHTTPHandler), "unix:"+tmpSock)

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

	res, err := client.Do(req)
	if err != nil {
		t.Fatal("error doing uds request", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusTeapot+2 {
		t.Fatal("got wrong status code:", res.StatusCode)
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
