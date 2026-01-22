package main

import (
    "context"
    "fmt"
    fn_events "github.com/fnproject/fn-events"
    "github.com/stretchr/testify/assert"
    "net/http"
    "testing"
)

func TestServe_ReturnsCreatedWithHeadersAndBody(t *testing.T) {
    handler := NewApiGatewayHandler()

    request := fn_events.APIGatewayRequestEvent[RequestEmployee]{
        Body: RequestEmployee{Name: "foo"},
        RequestURL: "/1?id=3",
        QueryParameters: map[string][]string{"id": {"3"}},
    }

    resp := handler.Serve(context.Background(), request)

    assert.Equal(t, http.StatusCreated, resp.StatusCode)
    assert.Equal(t, []string{"HeaderValue"}, resp.Headers["X-Custom-Header"])
    assert.Equal(t, []string{"HeaderValue2"}, resp.Headers["X-Custom-Header-2"])
    responseEmployee, ok := resp.Body.(*ResponseEmployee)

    assert.True(t, ok, "body should be of type *ResponseEmployee")
    assert.Equal(t, 3, responseEmployee.Id)
    assert.Equal(t, "foo", responseEmployee.Name)
}

func TestServe_ReturnsErrorWhenNoId(t *testing.T) {
    handler := NewApiGatewayHandler()

    request := fn_events.APIGatewayRequestEvent[RequestEmployee]{
        Body:       RequestEmployee{Name: "foo"},
        RequestURL: "/1",
    }

    resp := handler.Serve(context.Background(), request)

    assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    assert.Equal(t, fmt.Errorf("id must not be empty"), resp.Body)
}
