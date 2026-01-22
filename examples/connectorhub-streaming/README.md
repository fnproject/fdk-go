# Example Fn Go FDK : Service Connector Hub - Streaming

This example provides a Function to use as a service connector hub target.
The function accepts a typed event containing a batch of source events.

## Source
[StreamingData](../../fn-events/connectorhub_stream.go)

## Dependencies
* [fn-events] for ConnectorHubHandlerStreaming fdk handler.

## Demonstrated FDK features
This example showcases how to use the fn-event ConnectorHubHandlerStreaming to
use a Function as the target for Streaming source.

## Step by step

Set up the connector hub with Streaming source and Function target:
* [Setup default policies](https://docs.oracle.com/en-us/iaas/Content/connector-hub/overview.htm#Authenti__default-policies)
* [create connector hub](https://docs.oracle.com/en-us/iaas/Content/connector-hub/create-service-connector-streaming-source.htm)

The Function fdk.handler accepts the event handler.
```go
    handler := fn_events.ConnectorHubHandlerStreaming(
    NewExampleConnectorHubHandler(),
    reflect.TypeOf(Employee{}),
    )
    fdk.Handle(handler)
```

Full example:
```go
package main

import (
    "context"
    "log"
    "reflect"

    "github.com/fnproject/fdk-go"
    fn_events "github.com/fnproject/fn-events"
)

type ExampleConnectorHubHandler struct {
    streamService StreamService
}

func NewExampleConnectorHubHandler() *ExampleConnectorHubHandler {
    return &ExampleConnectorHubHandler{
        streamService: &RealStreamService{},
    }
}

func (h *ExampleConnectorHubHandler) Serve(
    ctx context.Context,
    batch fn_events.ConnectorHubBatch[fn_events.StreamingData[Employee]],
) {
    for _, streamingData := range batch.Batch {
        if err := h.streamService.ReadStream(streamingData); err != nil {
            log.Fatalf("error reading stream: %v", err)
        }
    }
}

func main() {
    log.SetFlags(0) // Removes redundant timestamps from logs
    handler := fn_events.ConnectorHubHandlerStreaming(
        NewExampleConnectorHubHandler(),
        reflect.TypeOf(Employee{}),
    )
    fdk.Handle(handler)
}
```

The [ConnectorHubBatch](../../fn-events/connectorhub.go)
`batch` contains a list of events from Streaming as
specified in [Batch Settings](https://docs.oracle.com/en-us/iaas/Content/connector-hub/overview.htm#batch-settings).

The Employee struct is the Value from each Stream event.
The Value sent to the Function is base64 encoded, the convenience handler automatically decodes and
Unmarshals the value into your specified struct or string.

To return an error response, exit the Function e.g log.Fatalf().
Doing so will cause the Function to return a 502 [Retry policy](https://docs.oracle.com/en-us/iaas/Content/connector-hub/overview.htm#deactivate)
