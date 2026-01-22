package main

import (
    "log"
)

type Employee struct {
    Name string `json:"name"`
}

type QueueService interface {
    ReadContent(content Employee) error
}

type RealQueueService struct{}

func (es *RealQueueService) ReadContent(content Employee) error {
    log.Printf("ReadContent: %v", content)
    return nil
}
