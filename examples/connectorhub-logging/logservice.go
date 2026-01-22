package main

import (
    "fmt"
    fn_events "github.com/fnproject/fn-events"
    "log"
)

type LogService interface {
    ReadLog(loggingData fn_events.LoggingData) error
}

type RealLogService struct{}

func (es *RealLogService) ReadLog(loggingData fn_events.LoggingData) error {
    log.Printf("loggingData: %v", loggingData)

    if len(loggingData.Data) < 1 {
        return fmt.Errorf("loggingData.Data is empty")
    }
    if loggingData.ID == "" {
        return fmt.Errorf("loggingData.ID is empty")
    }
    if len(loggingData.Oracle) < 1 {
        return fmt.Errorf("loggingData.Oracle is empty")
    }
    if loggingData.Source == "" {
        return fmt.Errorf("loggingData.Source is empty")
    }
    if loggingData.SpecVersion == "" {
        return fmt.Errorf("loggingData.Specversion is empty")
    }
    if loggingData.Time.IsZero() {
        return fmt.Errorf("loggingData.Time is empty")
    }
    if loggingData.Type == "" {
        return fmt.Errorf("loggingData.Type is empty")
    }
    return nil
}