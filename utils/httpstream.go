package utils

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type HTTPHandler struct {
	handler Handler
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var buf bytes.Buffer
	hdr := make(http.Header)

	ctx := WithContext(r.Context(), &Ctx{
		Config: BuildConfig(),
	})

	buf.Reset()
	ResetHeaders(hdr)
	resp := Response{
		Writer: &buf,
		Status: 200,
		Header: hdr,
	}

	fnDeadline := Context(ctx).Header.Get("FN_DEADLINE")
	ctx, cancel := CtxWithDeadline(ctx, fnDeadline)
	defer cancel()

	SetHeaders(ctx, r.Header)
	SetRequestURL(ctx, r.URL.String())
	SetMethod(ctx, r.Method)
	h.handler.Serve(ctx, r.Body, &resp)

	hResp := GetHTTPStreamResp(&buf, &resp, r)
	hResp.Write(w)
}

func StartHTTPServer(handler Handler, path, format string) {

	if format != "httpstream" {
		panic("expecting httpstream, invalid format: " + format)
	}

	tokens := strings.Split(path, ":")
	if len(tokens) != 2 {
		panic("cannot process listener path: " + path)
	}

	server := http.Server{
		Handler: &HTTPHandler{
			handler: handler,
		},
	}

	// try to remove pre-existing UDS: ignore errors here
	if tokens[0] == "unix" {
		os.Remove(tokens[1])
	}

	listener, err := net.Listen(tokens[0], tokens[1])
	if err != nil {
		panic("net.Listen error: " + err.Error())
	}

	err = server.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		panic("server.Serve error: " + err.Error())
	}
}

func GetHTTPStreamResp(buf *bytes.Buffer, fnResp *Response, req *http.Request) http.Response {

	fnResp.Header.Set("Content-Length", strconv.Itoa(buf.Len()))

	hResp := http.Response{
		ProtoMajor:    1,
		ProtoMinor:    1,
		StatusCode:    fnResp.Status,
		Request:       req,
		Body:          ioutil.NopCloser(buf),
		ContentLength: int64(buf.Len()),
		Header:        fnResp.Header,
	}

	return hResp
}
