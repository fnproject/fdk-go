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