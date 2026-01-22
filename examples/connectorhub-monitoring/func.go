package main

import (
    "context"
    "log"

    "github.com/fnproject/fdk-go"
    fn_events "github.com/fnproject/fn-events"
)

type ExampleConnectorHubHandler struct {
    metricService MetricService
}

func NewConnectorHubHandler() *ExampleConnectorHubHandler {
    return &ExampleConnectorHubHandler{
        metricService: &RealMetricService{},
    }
}

func (h *ExampleConnectorHubHandler) Serve(ctx context.Context, batch fn_events.ConnectorHubBatch[fn_events.MetricData]) {
    for _, item := range batch.Batch {
        if err := h.metricService.ReadMetric(item); err != nil {
            log.Fatalf("error reading metric: %v", err)
        }
    }
}

func main() {
    log.SetFlags(0) // Removes redundant timestamps from logs
    handler := fn_events.ConnectorHubMonitoringHandler(NewConnectorHubHandler())
    fdk.Handle(handler)
}