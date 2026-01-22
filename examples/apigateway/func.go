package main

import (
    "context"
    "github.com/fnproject/fdk-go"
    fn_events "github.com/fnproject/fn-events"
    "net/http"
    "reflect"
)

type RequestEmployee struct {
    Name string `json:"name"`
}

type ExampleAPIGatewayHandler struct {
    employeeService EmployeeService
}

func NewApiGatewayHandler() *ExampleAPIGatewayHandler {
    return &ExampleAPIGatewayHandler{
        employeeService: &RealEmployeeService{},
    }
}

func (h *ExampleAPIGatewayHandler) Serve(ctx context.Context, requestEvent fn_events.APIGatewayRequestEvent[RequestEmployee]) *fn_events.APIGatewayResponseEvent[any] {
    requestEmployee := requestEvent.Body

    id := requestEvent.QueryParameters.Get("id")

    responseEmployee, err := h.employeeService.CreateEmployee(&requestEmployee, id)
    if err != nil {
        builder := fn_events.NewAPIGatewayResponseEventBuilder[any]()
        return builder.Body(err).StatusCode(http.StatusBadRequest).Build()
    }
    headers := fn_events.Headers{
        "X-Custom-Header":   {"HeaderValue"},
        "X-Custom-Header-2": {"HeaderValue2"},
        "Content-Type":      {"application/json"},
    }

    builder := fn_events.NewAPIGatewayResponseEventBuilder[any]()
    return builder.Body(responseEmployee).StatusCode(http.StatusCreated).Headers(headers).Build()
}

func main() {
    handler := fn_events.APIGatewayHandler(NewApiGatewayHandler(), reflect.TypeOf(RequestEmployee{}))
    fdk.Handle(handler)
}
