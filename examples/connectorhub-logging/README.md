# Example Fn Go FDK : Service Connector Hub - Logging

This example provides a Function to use as a service connector hub target.
The function accepts a typed event containing a batch of source events.

## Source
[LoggingData](../../fn-events/connectorhub.go)

## Dependencies
* [fn-events] for ConnectorHubLoggingHandler fdk handler.

## Demonstrated FDK features
This example showcases how to use the fn-event ConnectorHubLoggingHandler to 
use a Function as the target for Logging source.

## Step by step

Set up the connector hub with Logging source and Function target:
* [Setup default policies](https://docs.oracle.com/en-us/iaas/Content/connector-hub/overview.htm#Authenti__default-policies)
* [create connector hub](https://docs.oracle.com/en-us/iaas/Content/connector-hub/create-service-connector-logging-source.htm)

The Function fdk.handler accepts the event handler.
```go
    handler := fn_events.ConnectorHubLoggingHandler(NewExampleConnectorHubHandler())
    fdk.Handle(handler)
```

Full example:
```go
package main

import (
    "context"
    "log"

    "github.com/fnproject/fdk-go"
    fn_events "github.com/fnproject/fn-events"
)

type ExampleConnectorHubHandler struct {
    logService LogService
}

func NewExampleConnectorHubHandler() *ExampleConnectorHubHandler {
    return &ExampleConnectorHubHandler{
        logService: &RealLogService{},
    }
}

func (h *ExampleConnectorHubHandler) Serve(ctx context.Context, batch fn_events.ConnectorHubBatch[fn_events.LoggingData]) {
    for _, item := range batch.Batch {
        if err := h.logService.ReadLog(item); err != nil {
            log.Fatalf("error: %v", err)
        }
    }
}

func main() {
    log.SetFlags(0) // Removes redundant timestamps from logs
    handler := fn_events.ConnectorHubLoggingHandler(NewExampleConnectorHubHandler())
    fdk.Handle(handler)
}
```

The [ConnectorHubBatch](../../fn-events/connectorhub.go)
`batch` contains a list of events from Logging as 
specified in [Batch Settings](https://docs.oracle.com/en-us/iaas/Content/connector-hub/overview.htm#batch-settings).

The LoggingData struct is 
each logging event.

To return an error response, exit the Function e.g log.Fatalf().
Doing so will cause the Function to return a 502 [Retry policy](https://docs.oracle.com/en-us/iaas/Content/connector-hub/overview.htm#deactivate)
