package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	fdk "github.com/fnproject/fdk-go"
	zipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	zipkinHttpReporter "github.com/openzipkin/zipkin-go/reporter/http"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

type Person struct {
	Name string `json:"name"`
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	p := &Person{Name: "World"}
	json.NewDecoder(in).Decode(p)
	newCtx := fdk.GetContext(ctx)

	// Span main method
	// set up a span reporter
	reporter := zipkinHttpReporter.NewReporter(newCtx.TracingContextData().TraceCollectorURL())
	defer reporter.Close()

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(newCtx.ServiceName(), "")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))

	// Set Context
	sctx := setContext(newCtx)
	sopt := zipkin.Parent(sctx)

	span, ctx := tracer.StartSpanFromContext(context.Background(), "Main Method", sopt)

	time.Sleep(15 * time.Millisecond)
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	oneMethod(ctx, tracer)

	msg := struct {
		Msg string `json:"message"`
	}{
		Msg: fmt.Sprintf("Hello %s :: and Service :: %s", newCtx.AppID(), newCtx.AppName()),
	}
	log.Print("Inside Go Hello World function")
	json.NewEncoder(out).Encode(&msg)
	span.Finish()
}

func setContext(ctx fdk.Context) model.SpanContext {
	traceId, err := model.TraceIDFromHex(ctx.TracingContextData().TraceId())

	if err != nil {
		log.Println("TRACE ID NOT DEFINED.....")
		return model.SpanContext{}
	}

	id, err := strconv.ParseUint(ctx.TracingContextData().SpanId(), 16, 64)
	if err != nil {
		log.Println("SPAN ID NOT DEFINED.....")
		return model.SpanContext{}
	}
	i := model.ID(id)

	sctx := model.SpanContext{
		TraceID:  traceId,
		ID:       i,
		ParentID: nil,
		Sampled:  BoolAddr(ctx.TracingContextData().IsSampled()),
	}

	return sctx
}

func BoolAddr(b bool) *bool {
	boolVar := b
	return &boolVar
}

func oneMethod(ctx context.Context, tracer *zipkin.Tracer) {
	span, newCtx := tracer.StartSpanFromContext(ctx, "OneChildSpan")
	time.Sleep(80 * time.Millisecond)
	ctx = newCtx
	secMethod(ctx, tracer)
	span.Finish()
}

func secMethod(ctx context.Context, tracer *zipkin.Tracer) {
	span, newCtx := tracer.StartSpanFromContext(ctx, "SecChildSpan")
	time.Sleep(100 * time.Millisecond)
	ctx = newCtx
	span.Finish()
}
