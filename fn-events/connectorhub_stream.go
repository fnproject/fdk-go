package fn_events

import (
    "context"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "github.com/fnproject/fdk-go"
    "io"
    "log"
    "reflect"
    "time"
)

type ConnectorHubHandleStreamWrapper[T any] struct {
    handler    ConnectorHubStreamFnHandler[T]
    sourceType reflect.Type
    valueType  reflect.Type
    FatalFuncWrapper
}

type ConnectorHubStreamFnHandler[T any] interface {
    Serve(ctx context.Context, batch ConnectorHubBatch[StreamingData[T]])
}

// ConnectorHubHandlerStreaming returns an fdk.Handler that wraps a customer-supplied streaming handler,
// providing type-safe, automatic unmarshaling and base64 decoding for each streaming event element.
//
// Parameters:
//   - fnHandler: An implementation of ConnectorHubStreamFnHandler[T]—your event handler—which processes batches
//     of streaming events, where each event's Value field is of type T.
//   - valueType: The reflect.Type of T, required so that event payloads are correctly decoded and unmarshaled
//     at runtime. For example, use reflect.TypeOf(MyStruct{}) or reflect.TypeOf("") for string payloads.
//
// Returns:
//   - fdk.Handler: An FDK handler that can be registered to process streaming events. It will handle
//     decoding, error reporting, and proper dispatch to your handler.
//
// Example usage:
//
//   type MyPayload struct { ... }
//   type MyStreamHandler struct {}
//   func (h MyStreamHandler) Serve(ctx context.Context, batch ConnectorHubBatch[StreamingData[MyPayload]]) { ... }
//
//   fdk.Handle(
//       ConnectorHubHandlerStreaming[MyPayload](&MyStreamHandler{}, reflect.TypeOf(MyPayload{})),
//   )
//
// For string payloads, use:
//   ConnectorHubHandlerStreaming[string](&MyHandler{}, reflect.TypeOf(""))
//
func ConnectorHubHandlerStreaming[T any](fnHandler ConnectorHubStreamFnHandler[T], valueType reflect.Type) fdk.Handler {
    return &ConnectorHubHandleStreamWrapper[T]{
        handler:    fnHandler,
        sourceType: reflect.TypeOf(StreamingData[T]{}), // the type for the main struct
        valueType:  valueType, // the type for the Value field
        FatalFuncWrapper: FatalFuncWrapper{
            fatal: func(args ...interface{}) {
                log.Fatal(args...)
            },
        },
    }
}

func (wrapper *ConnectorHubHandleStreamWrapper[T]) Serve(ctx context.Context, in io.Reader, out io.Writer) {
    items, err := wrapper.readFnInput(in)
    if err != nil {
        wrapper.HandleError(err)
        return
    }
    batch, _ := items.([]StreamingData[T])
    fdkContext := FetchContext(ctx)
    connectorHubBatch := ConnectorHubBatch[StreamingData[T]]{
        Batch:   batch,
        Headers: fdkContext.Header(),
    }
    wrapper.handler.Serve(ctx, connectorHubBatch)
}

func (wrapper *ConnectorHubHandleStreamWrapper[T]) readFnInput(in io.Reader) (any, error) {
    data, err := io.ReadAll(in)
    if err != nil {
        return nil, fmt.Errorf("failed to read input: %w", err)
    }
    if len(data) == 0 {
        return nil, nil
    }
    // Use a temporary struct for unmarshalling
    var raws []struct {
        Stream    string `json:"stream"`
        Partition string `json:"partition"`
        Key       *string `json:"key"`
        Value     string `json:"value"`
        Offset    int    `json:"offset"`
        Timestamp int64  `json:"timestamp"`
    }
    if err := json.Unmarshal(data, &raws); err != nil {
        return nil, err
    }

    // Build results as []StreamingData[T]
    result := make([]StreamingData[T], len(raws))
    for i, raw := range raws {
        decodedBytes, err := base64.StdEncoding.DecodeString(raw.Value)
        if err != nil && raw.Value != "" {
            return nil, fmt.Errorf("base64 decode: %w", err)
        }
        var value T
        switch wrapper.valueType.Kind() {
        case reflect.Struct, reflect.Map, reflect.Slice:
            valuePtr := reflect.New(wrapper.valueType)
            if err := json.Unmarshal(decodedBytes, valuePtr.Interface()); err != nil {
                return nil, fmt.Errorf("json decode on field 'value': %w", err)
            }
            value = valuePtr.Elem().Interface().(T)
        case reflect.String:
            v := string(decodedBytes)
            value = any(v).(T)
        default:
            // For []byte or other primitives
            value = any(decodedBytes).(T)
        }
        result[i] = StreamingData[T]{
            Stream:    raw.Stream,
            Partition: raw.Partition,
            Key:       raw.Key,
            Value:     value,
            Offset:    raw.Offset,
            Timestamp: raw.Timestamp,
        }
    }
    return result, nil
}

type StreamingData[T any] struct {
    Stream    string `json:"stream"`
    Partition string `json:"partition"`
    Key       *string `json:"key"`
    Offset    int    `json:"offset"`
    Timestamp int64  `json:"timestamp"`
    Value     T      `json:"value"`
}

func (sd *StreamingData[T]) Time() time.Time {
    return time.UnixMilli(sd.Timestamp)
}