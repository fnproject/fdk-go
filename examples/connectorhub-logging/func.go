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