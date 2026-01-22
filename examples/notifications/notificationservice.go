package main

import (
    fn_events "github.com/fnproject/fn-events"
    "log"
)

type Employee struct {
    Name string `json:"name"`
}

type NotificationService interface {
    ReadNotification(message fn_events.NotificationMessage[Employee])
}

type RealNotificationService struct{}

func (es *RealNotificationService) ReadNotification(message fn_events.NotificationMessage[Employee]) {
    log.Printf("ReadNotification: %v", message)
    employee := message.Content
    log.Printf("message content: %v", employee)
}
