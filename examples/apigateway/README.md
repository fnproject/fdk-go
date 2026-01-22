# Example Fn Go FDK : API Gateway

This example provides a Function to use as an API Gateway backend.
The function accepts a typed request for easy object handling and returns
http response.

## Dependencies

* [fn-events] for fn_events.APIGatewayHandler handlers.

## Demonstrated FDK features

This example showcases how to use the fn-event APIGatewayFunction to
use a Function as the backend.

## Step by step

Set the API Gateway and Function
backend [Adding a Function in OCI Functions as an API Gateway Back End](https://docs.oracle.com/en-us/iaas/Content/APIGateway/Tasks/apigatewayusingfunctionsbackend.htm)

The Function entrypoint accepts the `fn_events.APIGatewayHandler` handler.

[func.go](func.go)
```go
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

```
The APIGatewayRequestEvent `RequestURL` is relative to the deployment
path prefix [see API Gateway using HTTP backend](https://docs.oracle.com/en-us/iaas/Content/APIGateway/Tasks/apigatewayusinghttpbackend.htm#usingjson)


The struct RequestEmployee is the body and Serve is always provided with the Pointer.
[ResponseEmployee](employee.go) is the response body.

To return an error response construct it. 
```go
        builder := fn_events.NewAPIGatewayResponseEventBuilder[any]()
        return builder.Body(err).StatusCode(http.StatusBadRequest).Build()
```
