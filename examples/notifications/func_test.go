package main

import (
    "context"
    "github.com/stretchr/testify/assert"
    "testing"

    fn_events "github.com/fnproject/fn-events"
)

type MockNotificationService struct {
    ReadNotificationFunc func(data fn_events.NotificationMessage[Employee])
    Called               int
}

func (m *MockNotificationService) ReadNotification(data fn_events.NotificationMessage[Employee]) {
    m.Called++
    m.ReadNotificationFunc(data)
}

func TestNotificationHandlerCallsNotificationService(t *testing.T) {
    mockService := &MockNotificationService{
        ReadNotificationFunc: func(data fn_events.NotificationMessage[Employee]) {
            // Simulate a successful ReadContent call
        },
    }
    handler := &ExampleNotificationsHandler{
        notificationService: mockService,
    }
    message := fn_events.NotificationMessage[Employee]{Content: Employee{Name: "foo"}}

    handler.Serve(context.Background(), message)
    assert.Equal(t, 1, mockService.Called)

}
