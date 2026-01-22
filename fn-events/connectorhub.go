package fn_events

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/fnproject/fdk-go"
    "io"
    "log"
    "net/http"
    "reflect"
    "time"
)

type ConnectorHubHandleWrapper[T any] struct {
    handler  ConnectorHubFnHandler[T]
    bodyType reflect.Type
    FatalFuncWrapper
}

type ConnectorHubFnHandler[T any] interface {
    Serve(ctx context.Context, batch ConnectorHubBatch[T])
}

// ConnectorHubHandler creates and returns an fdk.Handler that wraps a customer-supplied handler
// function and provides type-safe automatic JSON deserialization for incoming ConnectorHub events.
//
// Parameters:
//   - fnHandler:  An implementation of ConnectorHubFnHandler[T], where T is the struct type that incoming
//                 events should deserialize into.
//   - sourceType: The reflect.Type of T, used to support generic decoding of input payloads.
//
// Returns:
//   - fdk.Handler: A handler ready for registration with the FN framework, with error handling
//                  defaulting to log.Fatal.
//
// Example usage:
//
//    type CustomPayload struct { ... }
//    type MyHandler struct {}
//    func (h *MyHandler) Serve(ctx context.Context, batch ConnectorHubBatch[CustomPayload]) { ... }
//    handler := ConnectorHubHandler(&MyHandler{}, reflect.TypeOf(CustomPayload{}))
//
func ConnectorHubHandler[T any](fnHandler ConnectorHubFnHandler[T], sourceType reflect.Type) fdk.Handler {
    return &ConnectorHubHandleWrapper[T]{
        handler:  fnHandler,
        bodyType: sourceType,
        FatalFuncWrapper: FatalFuncWrapper{
            fatal: func(args ...interface{}) {
                log.Fatal(args...)
            },
        },
    }
}

// ConnectorHubLoggingHandler is a convenience function for registering a ConnectorHub handler
// that processes LoggingData events. It instantiates a ConnectorHubHandler with LoggingData for strong typing.
//
// Parameters:
//   - fnHandler:  The customer's ConnectorHubFnHandler[LoggingData] implementation.
//
// Returns:
//   - fdk.Handler: A ready-to-register handler for logging events.
//
// Example usage:
//
//    func (h *MyLoggingHandler) Serve(ctx context.Context, batch ConnectorHubBatch[LoggingData]) { ... }
//    handler := ConnectorHubLoggingHandler(&MyLoggingHandler{})
//
func ConnectorHubLoggingHandler(fnHandler ConnectorHubFnHandler[LoggingData]) fdk.Handler {
    return ConnectorHubHandler(fnHandler, reflect.TypeOf(LoggingData{}))
}

// ConnectorHubMonitoringHandler is a convenience function for registering a ConnectorHub handler
// that processes MetricData events. It instantiates a ConnectorHubHandler with MetricData for strong typing.
//
// Parameters:
//   - fnHandler:  The customer's ConnectorHubFnHandler[MetricData] implementation.
//
// Returns:
//   - fdk.Handler: A ready-to-register handler for monitoring/metric events.
//
// Example usage:
//
//    func (h *MyMetricHandler) Serve(ctx context.Context, batch ConnectorHubBatch[MetricData]) { ... }
//    handler := ConnectorHubMonitoringHandler(&MyMetricHandler{})
//
func ConnectorHubMonitoringHandler(fnHandler ConnectorHubFnHandler[MetricData]) fdk.Handler {
    return ConnectorHubHandler(fnHandler, reflect.TypeOf(MetricData{}))
}

// ConnectorHubQueueHandler is a convenience function for registering a ConnectorHub handler
// that processes Queue events. It instantiates a ConnectorHubHandler with a custom struct for strong typing.
//
// Parameters:
//   - fnHandler:  The customer's ConnectorHubFnHandler[MyStruct] implementation.
//   - valueType: The reflect.Type of T, required so that queue payloads are correctly decoded and unmarshaled
//     at runtime. For example, use reflect.TypeOf(MyStruct{}) or reflect.TypeOf("") for string payloads.
//
// Returns:
//   - fdk.Handler: A ready-to-register handler for queue events.
//
// Example usage:
//
//   type MyPayload struct { ... }
//   type MyQueueHandler struct {}
//   func (h MyQueueHandler) Serve(ctx context.Context, batch ConnectorHubBatch[MyPayload]) { ... }
//
//   fdk.Handle(
//       ConnectorHubHandlerQueue[MyPayload](&MyQueueHandler{}, reflect.TypeOf(MyPayload{})),
//   )
//
// For string payloads, use:
//   ConnectorHubHandlerQueue[string](&MyQueueHandler{}, reflect.TypeOf(""))
//
func ConnectorHubQueueHandler[T any](fnHandler ConnectorHubFnHandler[T], contentType reflect.Type) fdk.Handler {
    return ConnectorHubHandler(fnHandler, contentType)
}

// Serve implements fdk.Handler for ConnectorHubHandleWrapper.
// It reads JSON input, deserializes to a slice of T, and dispatches to the handler.
func (wrapper *ConnectorHubHandleWrapper[T]) Serve(ctx context.Context, in io.Reader, out io.Writer) {
    batch, err := wrapper.readFnInput(in)
    if err != nil {
        wrapper.HandleError(err)
        return
    }
    var batchTyped []T
    if batch != nil {
        batchTyped = batch.([]T)
    }
    fdkContext := FetchContext(ctx)
    connectorHubBatch := ConnectorHubBatch[T]{
        Batch:   batchTyped,
        Headers: fdkContext.Header(),
    }
    wrapper.handler.Serve(ctx, connectorHubBatch)
}

func (wrapper *ConnectorHubHandleWrapper[T]) readFnInput(in io.Reader) (any, error) {
    data, err := io.ReadAll(in)
    if err != nil {
        return nil, fmt.Errorf("failed to read input: %w", err)
    }
    if len(data) == 0 {
        return nil, nil
    }
    sliceType := reflect.SliceOf(wrapper.bodyType)
    slicePtr := reflect.New(sliceType)
    if err := json.Unmarshal(data, slicePtr.Interface()); err != nil {
        return nil, fmt.Errorf("failed to unmarshal to slice of %v, %v: %w", wrapper.bodyType, string(data), err)
    }
    return slicePtr.Elem().Interface(), nil // This is []T
}

type ConnectorHubBatch[T any] struct {
    Batch   []T
    Headers http.Header
}

type LoggingData struct {
    ID          string            `json:"id"`
    Source      string            `json:"source"`
    SpecVersion string            `json:"specversion"`
    Subject     string            `json:"subject"`
    Type        string            `json:"type"`
    Data        map[string]string `json:"data"`
    Oracle      map[string]string `json:"oracle"`
    Time        time.Time         `json:"time"`
}

type MetricData struct {
    Namespace     string            `json:"namespace"`
    ResourceGroup string            `json:"resourceGroup"`
    CompartmentId string            `json:"compartmentId"`
    Name          string            `json:"name"`
    Dimensions    map[string]string `json:"dimensions"`
    Metadata      map[string]string `json:"metadata"`
    Datapoints    []Datapoint       `json:"datapoints"`
}

type Datapoint struct {
    Timestamp int64    `json:"timestamp"`
    Value     *float64 `json:"value"`
    Count     *int     `json:"count"`
}

func (e *Datapoint) Time() time.Time {
    return time.UnixMilli(e.Timestamp)
}