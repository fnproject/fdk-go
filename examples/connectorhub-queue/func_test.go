package main

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    fn_events "github.com/fnproject/fn-events"
)

// MockQueueService provides a mock for QueueService.
type MockQueueService struct {
    ReadContentFunc func(data Employee) error
    Called          int
}

func (m *MockQueueService) ReadContent(queueContent Employee) error {
    m.Called++
    return m.ReadContentFunc(queueContent)
}

func TestQueueHandlerCallsQueueServiceForEachItem(t *testing.T) {
    mockService := &MockQueueService{
        ReadContentFunc: func(content Employee) error {
            return nil
        },
    }
    handler := &ExampleConnectorHubHandler{
        queueService: mockService,
    }
    queueContent := Employee{Name: "foo"}
    queueContent2 := Employee{Name: "bar"}
    batch := fn_events.ConnectorHubBatch[Employee]{Batch: []Employee{queueContent, queueContent2}}

    handler.Serve(context.Background(), batch)

    assert.Equal(t, 2, mockService.Called)
}