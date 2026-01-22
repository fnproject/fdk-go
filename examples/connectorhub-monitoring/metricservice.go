package main

import (
    "fmt"
    fn_events "github.com/fnproject/fn-events"
    "log"
)

type MetricService interface {
    ReadMetric(metricData fn_events.MetricData) error
}

type RealMetricService struct{}

func (es *RealMetricService) ReadMetric(metricData fn_events.MetricData) error {
    log.Printf("metricData: %v", metricData)

    if len(metricData.Datapoints) < 1 {
        return fmt.Errorf("metricData.Data is empty")
    }
    if metricData.CompartmentId == "" {
        return fmt.Errorf("metricData.CompartmentId is empty")
    }
    if len(metricData.Dimensions) < 1 {
        return fmt.Errorf("metricData.Dimensions is empty")
    }
    if len(metricData.Metadata) < 1 {
        return fmt.Errorf("metricData.Metadata is empty")
    }
    if metricData.Name == "" {
        return fmt.Errorf("metricData.Name is empty")
    }
    if metricData.Namespace == "" {
        return fmt.Errorf("metricData.Namespace is empty")
    }
    return nil
}