package main

import (
    "fmt"
    fn_events "github.com/fnproject/fn-events"
    "log"
)

type Employee struct {
    Name string `json:"name"`
}

type StreamService interface {
    ReadStream(streamingData fn_events.StreamingData[Employee]) error
}

type RealStreamService struct{}

func (es *RealStreamService) ReadStream(streamingData fn_events.StreamingData[Employee]) error {
    log.Printf("ReadContent: %v", streamingData)

    if streamingData.Stream == "" {
        return fmt.Errorf("streamingData.Stream is empty")
    }
    if streamingData.Partition == "" {
        return fmt.Errorf("streamingData.Partition is empty")
    }
    employee := streamingData.Value
    if isEmptyEmployee(employee) {
        return fmt.Errorf("streamingData.Value (Employee) is empty")
    }
    if streamingData.Timestamp == 0 {
        return fmt.Errorf("streamingData.Timestamp is zero")
    }

    log.Printf("streamingData.Value: %v", streamingData.Value)
    return nil
}

func isEmptyEmployee(emp Employee) bool {
    return emp.Name == ""
}