package main

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    fn_events "github.com/fnproject/fn-events"
)

type MockStreamingService struct {
    ReadStreamFunc func(data fn_events.StreamingData[Employee]) error
    Called         int
}

func (m *MockStreamingService) ReadStream(streamingData fn_events.StreamingData[Employee]) error {
    m.Called++
    return m.ReadStreamFunc(streamingData)
}

func TestStreamHandlerCallsStreamServiceForEachItem(t *testing.T) {
    mockService := &MockStreamingService{
        ReadStreamFunc: func(data fn_events.StreamingData[Employee]) error {
            // Simulate a successful ReadStream call
            return nil
        },
    }
    handler := &ExampleConnectorHubHandler{
        streamService: mockService,
    }
    batch := fn_events.ConnectorHubBatch[fn_events.StreamingData[Employee]]{
        Batch: []fn_events.StreamingData[Employee]{
            {Value: Employee{Name: "foo"}},
            {Value: Employee{Name: "Bar"}},
        },
    }

    handler.Serve(context.Background(), batch)
    assert.Equal(t, 2, mockService.Called)
}
