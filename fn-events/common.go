package fn_events

import (
    "context"
    "github.com/fnproject/fdk-go"
    "log"
)

type QueryParameters map[string][]string
type Headers map[string][]string

type FatalFuncWrapper struct {
    fatal func(args ...interface{})
}

func (wrapper *FatalFuncWrapper) Fatal(args ...interface{}) {
    wrapper.fatal(args...)
}

func (wrapper *FatalFuncWrapper) HandleError(err error) {
    log.SetFlags(0) // Hides redundant timestamp in log content
    wrapper.fatal(err)
}

func FetchContext(ctx context.Context) fdk.Context {
    fdkContext := fdk.GetContext(ctx)

    return fdkContext
}
