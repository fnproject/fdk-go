package main

import (
    "errors"
    "strconv"
)

type ResponseEmployee struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
}

type EmployeeService interface {
    CreateEmployee(req *RequestEmployee, id string) (*ResponseEmployee, error)
}

type RealEmployeeService struct{}

func (es *RealEmployeeService) CreateEmployee(req *RequestEmployee, id string) (*ResponseEmployee, error) {
    if req == nil {
        return nil, errors.New("requestEmployee must not be nil")
    }
    if id == "" {
        return nil, errors.New("id must not be empty")
    }
    idInt, err := strconv.Atoi(id)
    if err != nil {
        return nil, errors.New("id must be an integer")
    }
    return &ResponseEmployee{
        Id:   idInt,
        Name: req.Name,
    }, nil
}
