package utils

import (
	"net"
	"net/http"
	"os"
	"strings"
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
