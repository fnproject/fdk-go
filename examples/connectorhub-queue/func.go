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
