package utils

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"syscall"
)

type HTTPHandler struct {
	handler Handler
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := WithContext(r.Context(), &Ctx{
		Config: BuildConfig(),
	})

	fnDeadline := Context(ctx).Header.Get("FN_DEADLINE")
	ctx, cancel := CtxWithDeadline(ctx, fnDeadline)
	defer cancel()

	SetHeaders(ctx, r.Header)
	SetRequestURL(ctx, r.URL.String())
	SetMethod(ctx, r.Method)
	h.handler.Serve(ctx, r.Body, w)

	// TODO can we get away with no buffer? set content length is 'nice' but they
	// can do it if they really need to... i lack ideas, now that we have a real
	// resp writer tho it's really worth considering.
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
	var oldmask int
	if uri.Scheme == "unix" {
		os.Remove(uri.Path)

		// this will give user perms to write to the sock file
		oldmask = syscall.Umask(0000)
	}

	listener, err := net.Listen(uri.Scheme, uri.Path)
	if err != nil {
		log.Fatalln("net.Listen error: ", err)
	}

	if uri.Scheme == "unix" {
		// set this back after we create the socket file. (of course there's a race, but who will ever hit it?)
		syscall.Umask(oldmask)
	}

	err = server.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln("serve error: ", err)
	}
}
