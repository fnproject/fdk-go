package fdk

import (
	"context"
	"io"
	"os"
	"strings"
)

type Handler interface {
	Serve(ctx context.Context, in io.Reader, out io.Writer)
}

type HandlerFunc func(ctx context.Context, in io.Reader, out io.Writer)

func Do(handler Handler) {
	format, _ := os.LookupEnv("FN_FORMAT")
	switch format {
	case "http":
		doHTTP(handler)
	case "default":
		doDefault(handler)
	default:
		panic("unknown format: " + format)
	}
}

type Context struct {
	Headers map[string][]string
	Config  map[string]string
}

type key struct{}

var ctxKey = new(key)

var (
	base = map[string]struct{}{
		`FN_APP_NAME`: struct{}{},
		`FN_PATH`:     struct{}{},
		`FN_METHOD`:   struct{}{},
		`FN_FORMAT`:   struct{}{},
		`FN_MEMORY`:   struct{}{},
		`FN_TYPE`:     struct{}{},
	}

	pres = [...]string{
		`FN_PARAM`,
		`FN_HEADER`,
	}

	exact = map[string]struct{}{
		`FN_CALL_ID`:     struct{}{},
		`FN_REQUEST_URL`: struct{}{},
	}
)

func header(key string) bool {
	for _, pre := range pres {
		if strings.HasPrefix(key, pre) {
			return true
		}
	}
	_, ok := exact[key]
	return ok
}

func buildConfig() map[string]string {
	cfg := make(map[string]string, len(base))

	for _, e := range os.Environ() {
		vs := strings.SplitN(e, "=", 1)
		if _, ok := base[vs[0]]; !ok {
			continue
		}
		if len(vs) < 2 {
			vs = append(vs, "")
		}
		cfg[vs[0]] = vs[1]
	}
	return cfg
}

func buildHeadersFromEnv() map[string][]string {
	env := os.Environ()
	hdr := make(map[string][]string, len(env)-len(base))

	for _, e := range env {
		vs := strings.SplitN(e, "=", 1)
		if !header(vs[0]) {
			continue
		}
		if len(vs) < 2 {
			vs = append(vs, "")
		}
		k := vs[0]
		// rebuild these as 'http' headers
		vs = strings.Split(vs[1], ", ")
		hdr[k] = vs
	}
	return hdr
}

// doDefault only runs once, since it is a 'cold' function
func doDefault(handler Handler) {
}
