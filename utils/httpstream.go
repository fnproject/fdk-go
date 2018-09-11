package utils

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
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
	phonySock := "/tmp/phonyfn.sock"
	if uri.Scheme == "unix" {
		os.Remove(phonySock)
	}

	listener, err := net.Listen(uri.Scheme, phonySock)
	if err != nil {
		log.Fatalln("net.Listen error: ", err)
	}

	if uri.Scheme == "unix" {
		// somehow this is the best way to get a permissioned sock file, don't ask questions, life is sad and meaningless
		f, err := os.Create(phonySock)
		if err != nil {
			log.Fatalln("error creating sock file", err)
		}

		err = f.Chmod(0666)
		if err != nil {
			f.Close()
			log.Fatalln("error giving sock file a perm", err)
		}
		f.Close()

		err = os.Symlink(uri.Path, phonySock)
		if err != nil {
			log.Fatalln("error giving sock file a perm", err)
		}
	}

	err = server.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln("serve error: ", err)
	}
}
