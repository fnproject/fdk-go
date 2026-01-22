# Example Fn Go FDK : Service Connector Hub - Queue

This example provides a Function to use as a service connector hub target.
The function accepts a typed event containing a batch of source events.

## Source
[ConnectorHubBatch](../../fn-events/connectorhub.go)

## Dependencies
* [fn-events] for ConnectorHubQueueHandler fdk handler.

## Demonstrated FDK features
This example showcases how to use the fn-event ConnectorHubQueueHandler to
use a Function as the target for Queue source.

## Step by step

Set up the connector hub with Queue source and Function target:
* [Setup default policies](https://docs.oracle.com/en-us/iaas/Content/connector-hub/overview.htm#Authenti__default-policies)
* [create connector hub](https://docs.oracle.com/en-us/iaas/Content/connector-hub/create-service-connector-queue-source.htm)

The Function fdk.handler accepts the event handler.
```go
    handler := fn_events.ConnectorHubQueueHandler(
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
    queueService QueueService
}

func NewExampleConnectorHubHandler() *ExampleConnectorHubHandler {
    return &ExampleConnectorHubHandler{
        queueService: &RealQueueService{},
    }
}

// Serve processes a batch of Employee queue messages.
func (h *ExampleConnectorHubHandler) Serve(ctx context.Context, batch fn_events.ConnectorHubBatch[Employee]) {
    for _, content := range batch.Batch {
        if err := h.queueService.ReadContent(content); err != nil {
            log.Fatalf("error queue content: %v", err)
        }
    }
}

func main() {
    log.SetFlags(0) // Removes redundant timestamps from logs
    handler := fn_events.ConnectorHubQueueHandler(
        NewExampleConnectorHubHandler(),
        reflect.TypeOf(Employee{}),
    )
    fdk.Handle(handler)
}
```

The [ConnectorHubBatch](../../fn-events/connectorhub.go)
`batch` contains a list of messages from Queue as
specified in [Batch Settings](https://docs.oracle.com/en-us/iaas/Content/connector-hub/overview.htm#batch-settings).

The Employee struct is the content of each message in the Queue.

To return an error response, exit the Function e.g log.Fatalf().
Doing so will cause the Function to return a 502 [Retry policy](https://docs.oracle.com/en-us/iaas/Content/connector-hub/overview.htm#deactivate)
