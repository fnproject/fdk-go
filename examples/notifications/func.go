package main

import (
    "context"
    "github.com/fnproject/fdk-go"
    fn_events "github.com/fnproject/fn-events"
    "log"
    "reflect"
)

type ExampleNotificationsHandler struct {
    notificationService NotificationService
}

func NewNotificationsHandler() *ExampleNotificationsHandler {
    return &ExampleNotificationsHandler{
        notificationService: &RealNotificationService{},
    }
}

func (h *ExampleNotificationsHandler) Serve(ctx context.Context, notification fn_events.NotificationMessage[Employee]) {
    h.notificationService.ReadNotification(notification)
}

func main() {
    log.SetFlags(0) // Removes redundant timestamps from logs
    handler := fn_events.NotificationsHandler(NewNotificationsHandler(),
        reflect.TypeOf(Employee{}),
    )
    fdk.Handle(handler)
}
