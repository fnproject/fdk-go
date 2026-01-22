# Example Fn Go FDK : Notifications

This example provides a Function to consume a Notification.
The function accepts a typed message event.

## Source
[NotificationMessage](../../fn-events/notifications.go)

## Dependencies
* [fn-events] for NotificationsHandler fdk handler.

## Demonstrated FDK features
This example showcases how to use the fn-event NotificationsHandler to
use a Function as the target for Notifications.

## Step by step

Set up the Notification Topic with Function Subscription:
* [Setup default policies](https://docs.oracle.com/en-us/iaas/Content/Security/Reference/notifications_security.htm#iam-policies__subs)
* [create subscription function](https://docs.oracle.com/en-us/iaas/Content/Notification/Tasks/create-subscription-function.htm)

The Function fdk.handler accepts the event handler.
```go
    handler := fn_events.NotificationsHandler(NewNotificationsHandler(),
  reflect.TypeOf(Employee{}),
  )
  fdk.Handle(handler)
```

Full example:
```go
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

```

The [NotificationMessage](../../fn-events/notifications.go)
`Content` contains the message from the Notification.

The Employee struct is used to unmarshalled the message. Alternatively, provide a string type for basic message.
```go
    handler := NotificationsHandler(
    MyNotificationHandler,
    reflect.TypeOf(""),
    )
```
