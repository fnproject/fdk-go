package utils

import (
	"bytes"
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// in case we go over the timeout, need to use a pool since prev buffer may not be freed
var bufPool = &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}

type HTTPHandler struct {
	handler Handler
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	resp := Response{
		Writer: buf,
		Status: 200,
		Header: make(http.Header), // XXX(reed): pool these too
	}

	ctx := WithContext(r.Context(), &Ctx{
		Config: BuildConfig(),
	})

	ctx, cancel := decapHeaders(ctx, r)
	defer cancel()

	h.handler.Serve(ctx, r.Body, &resp)

	encapHeaders(w, resp)

	// XXX(reed): 504 if ctx is past due / handle errors with 5xx? just 200 for now
	// copy response from user back up now with headers in place...
	io.Copy(w, buf)

	// XXX(reed): handle streaming, we have to intercept headers but not necessarily body (ie no buffer)
}

func encapHeaders(fn http.ResponseWriter, user Response) {
	fnh := fn.Header()
	fnh.Set("Fn-Http-Status", strconv.Itoa(user.Status))

	for k, vs := range user.Header {
		switch k {
		case "Content-Type":
			// don't modify this one...
		default:
			// prepend this guy
			k = "Fn-Http-H-" + k
		}

		for _, v := range vs {
			fnh.Add(k, v)
		}
	}
}

// TODO can make this the primary means of context construction
func decapHeaders(ctx context.Context, r *http.Request) (_ context.Context, cancel func()) {
	// XXX(reed): could make a new header bucket to dump things into instead of futzing

	rctx := Context(ctx)
	var deadline string

	for k, vs := range r.Header {
		switch k {
		// XXX(reed): we should strip out request url and method too but for invoke they don't exist...
		case "Fn-Deadline":
			r.Header.Del(k)
			deadline = vs[0]
		case "Fn-Call-Id":
			r.Header.Del(k)
			rctx.callId = vs[0]
		case "Content-Type":
			// just leave this one instead of deleting
		default:
			continue
		}

		if !strings.HasPrefix(k, "Fn-Http-H-") {
			// XXX(reed): we need 2 header buckets on ctx, one for these and one for the 'original req' headers
			// for now just nuke so the headers are clean...
			r.Header.Del(k)
			continue
		}

		r.Header.Del(k)
		for _, v := range vs {
			r.Header.Add(strings.TrimPrefix(k, "Fn-Http-H-"), v)
		}
	}

	return CtxWithDeadline(ctx, deadline)
}

func StartHTTPServer(handler Handler, path, format string) {

	uri, err := url.Parse(path)
	if err != nil {
		log.Fatalln("url parse error: ", path, err)
	}

	server := http.Server{
		Handler: &HTTPHandler{
			handler: handler,
		},
	}

	// try to remove pre-existing UDS: ignore errors here
	phonySock := filepath.Join(filepath.Dir(uri.Path), "phony"+filepath.Base(uri.Path))
	if uri.Scheme == "unix" {
		os.Remove(phonySock)
		os.Remove(uri.Path)
	}

	listener, err := net.Listen(uri.Scheme, phonySock)
	if err != nil {
		log.Fatalln("net.Listen error: ", err)
	}

	if uri.Scheme == "unix" {
		sockPerm(phonySock, uri.Path)
	}

	err = server.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln("serve error: ", err)
	}
}

func sockPerm(phonySock, realSock string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// somehow this is the best way to get a permissioned sock file, don't ask questions, life is sad and meaningless
	err := os.Chmod(phonySock, 0666)
	if err != nil {
		log.Fatalln("error giving sock file a perm", err)
	}

	err = os.Link(phonySock, realSock)
	if err != nil {
		log.Fatalln("error linking fake sock to real sock", err)
	}
}
