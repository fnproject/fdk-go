package main

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"

    fn_events "github.com/fnproject/fn-events"
)

type MockMetricService struct {
    ReadMetricFunc func(data fn_events.MetricData) error
    Called         int
}

func (m *MockMetricService) ReadMetric(metricData fn_events.MetricData) error {
    m.Called++
    return m.ReadMetricFunc(metricData)
}

func TestMonitoringHandlerCallsMetricServiceForEachItem(t *testing.T) {
    mockService := &MockMetricService{
        ReadMetricFunc: func(metric fn_events.MetricData) error {
            // Simulate successful ReadMetric
            return nil
        },
    }
    handler := &ExampleConnectorHubHandler{
        metricService: mockService,
    }
    metricData := fn_events.MetricData{Name: "foo"}
    metricData2 := fn_events.MetricData{Name: "bar"}
    batch := fn_events.ConnectorHubBatch[fn_events.MetricData]{
        Batch: []fn_events.MetricData{metricData, metricData2},
    }

    handler.Serve(context.Background(), batch)

    assert.Equal(t, 2, mockService.Called)
}