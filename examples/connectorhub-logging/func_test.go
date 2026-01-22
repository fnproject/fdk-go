package main

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"

    fn_events "github.com/fnproject/fn-events"
)

type MockLogService struct {
    ReadLogFunc func(fn_events.LoggingData) error
    Called      int
}

func (m *MockLogService) ReadLog(log fn_events.LoggingData) error {
    m.Called++
    return m.ReadLogFunc(log)
}

func TestLoggingHandlerCallsLogServiceForEachItem(t *testing.T) {
    mockService := &MockLogService{
        ReadLogFunc: func(log fn_events.LoggingData) error {
            // Simulate successful ReadLog call
            return nil
        },
    }
    handler := &ExampleConnectorHubHandler{
        logService: mockService,
    }
    logData := fn_events.LoggingData{ID: "foo"}
    logData2 := fn_events.LoggingData{ID: "bar"}
    batch := fn_events.ConnectorHubBatch[fn_events.LoggingData]{
        Batch: []fn_events.LoggingData{logData, logData2},
    }

    handler.Serve(context.Background(), batch)

    assert.Equal(t, 2, mockService.Called)
}