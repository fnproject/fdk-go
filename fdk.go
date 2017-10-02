package fdk

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Handler interface {
	Serve(ctx context.Context, in io.Reader, out io.Writer)
}

type HandlerFunc func(ctx context.Context, in io.Reader, out io.Writer)

func (f HandlerFunc) Serve(ctx context.Context, in io.Reader, out io.Writer) {
	f(ctx, in, out)
}

func Do(handler Handler) {
	format, _ := os.LookupEnv("FN_FORMAT")
	switch format {
	case "http":
		doHTTP(handler)
	case "default":
		doDefault(handler)
	default:
		panic("unknown format (fdk-go): " + format)
	}
}

// doDefault only runs once, since it is a 'cold' function
func doDefault(handler Handler) {
	ctx := buildCtx()
	setHeaders(ctx, buildHeadersFromEnv())

	// TODO we need to set deadline on ctx here (need FN_DEADLINE header)
	handler.Serve(ctx, os.Stdin, os.Stdout)
}

func doHTTP(handler Handler) {
	ctx := buildCtx()

	var buf bytes.Buffer
	// maps don't get down-sized, so we can reuse this as it's likely that the
	// user sends in the same amount of headers over and over (but still clear
	// b/w runs) -- buf uses same principle
	hdr := make(map[string][]string)

	for {
		// TODO we need to set deadline on ctx here (need FN_DEADLINE header)
		// for now, just get a new ctx each go round
		ctx, _ := context.WithCancel(ctx)

		buf.Reset()
		resetHeaders(hdr)
		resp := response{
			Writer: &buf,
			status: 200,
			header: hdr,
		}

		req, err := http.ReadRequest(bufio.NewReader(os.Stdin))
		if err != nil {
			// TODO it would be nice if we could let the user format this response to their preferred style..
			resp.status = http.StatusInternalServerError
			io.WriteString(resp, err.Error())
		} else {
			setHeaders(ctx, req.Header)
			handler.Serve(ctx, req.Body, &resp)
		}

		hResp := http.Response{
			ProtoMajor:    1,
			ProtoMinor:    1,
			StatusCode:    resp.status,
			Request:       req,
			Body:          ioutil.NopCloser(&buf),
			ContentLength: int64(buf.Len()),
			Header:        resp.header,
		}
		hResp.Write(os.Stdout)
	}
}

func resetHeaders(m map[string][]string) {
	for k := range m { // compiler optimizes this to 1 instruction now
		delete(m, k)
	}
}

// response is a general purpose response struct any format can use to record
// user's code responses before formatting them appropriately.
type response struct {
	status int
	header map[string][]string

	io.Writer
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

func setHeaders(ctx context.Context, hdr map[string][]string) {
	fctx := ctx.Value(ctxKey).(*Context)
	fctx.Headers = hdr
}

func buildCtx() context.Context {
	ctx := &Context{
		Config: buildConfig(),
		// allow caller to build headers separately (to avoid map alloc)
	}

	return context.WithValue(context.Background(), ctxKey, ctx)
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

func header(key string) bool {
	for _, pre := range pres {
		if strings.HasPrefix(key, pre) {
			return true
		}
	}
	_, ok := exact[key]
	return ok
}
