package fn_events

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "github.com/stretchr/testify/assert"
    "io"
    "log"
    "net/http"
    "reflect"
    "testing"
    "time"
    "github.com/fnproject/fdk-go"
)

type httpContext struct {
    header         http.Header
    config         map[string]string
    callID         string
    tracingContext fdk.TracingContext
}

func (c *httpContext) Deadline() (deadline time.Time, ok bool) { return }
func (c *httpContext) Done() <-chan struct{}                   { return nil }
func (c *httpContext) Err() error                              { return fmt.Errorf("") }
func (c *httpContext) Value(key any) any                       { return "" }
func (c *httpContext) Config() map[string]string               { return c.config }
func (c *httpContext) Header() http.Header                     { return c.header }
func (c *httpContext) ContentType() string                     { return c.header.Get("Content-Type") }
func (c *httpContext) CallID() string                          { return c.callID }
func (c *httpContext) AppID() string                           { return "" }
func (c *httpContext) FnID() string                            { return "" }
func (c *httpContext) AppName() string                         { return "" }
func (c *httpContext) RequestMethod() string                   { return c.header.Get("Fn-Http-Method") }
func (c *httpContext) RequestURL() string                      { return c.header.Get("Fn-Http-Request-Url") }
func (c *httpContext) FnName() string                          { return "" }
func (c *httpContext) TracingContextData() fdk.TracingContext  { return c.tracingContext }

type otherContext struct {
    header         http.Header
    config         map[string]string
    callID         string
    tracingContext fdk.TracingContext
}

func (c *otherContext) FnName() string                         { return "" }
func (c *otherContext) TracingContextData() fdk.TracingContext { return c.tracingContext }
func (c *otherContext) Config() map[string]string              { return c.config }
func (c *otherContext) Header() http.Header                    { return c.header }
func (c *otherContext) ContentType() string                    { return c.header.Get("Content-Type") }
func (c *otherContext) CallID() string                         { return c.callID }
func (c *otherContext) AppID() string                          { return "" }
func (c *otherContext) FnID() string                           { return "" }
func (c *otherContext) AppName() string                        { return "" }

type Employee struct {
    Name string `json:"name"`
}

type ApiGatewaySpy[T any] struct {
    CallCount   int
    LastCtx     context.Context
    LastRequest APIGatewayRequestEvent[T]
}

func (s *ApiGatewaySpy[T]) Serve(ctx context.Context, req APIGatewayRequestEvent[T]) *APIGatewayResponseEvent[any] {
    s.CallCount++
    s.LastCtx = ctx
    s.LastRequest = req
    builder := NewAPIGatewayResponseEventBuilder[any]()
    return builder.Body(req.Body).StatusCode(200).Headers(Headers{"X-Test": {"abc"}}).Build()
}

func structWrapperHandler[T any](body T) Employee {
    input, _ := json.Marshal(body)
    wrapper := &APIGatewayHandleWrapper[T]{
        handler:  &ApiGatewaySpy[T]{},
        bodyType: reflect.TypeOf(Employee{}),
    }
    var out bytes.Buffer
    ctx := &httpContext{
        header: http.Header{
            "Content-Type":        []string{"application/json"},
            "Fn-Http-Request-Url": []string{"/v1?param1=value%20with%20spaces"},
            "Fn-Http-Method":      []string{"POST"},
        },
        config: map[string]string{},
    }
    wrapper.Serve(fdk.WithContext(context.Background(), ctx), bytes.NewReader(input), &out)
    var result Employee
    _ = json.Unmarshal(out.Bytes(), &result)
    return result
}

func stringWrapperHandler(body any, wrapper *APIGatewayHandleWrapper[string]) string {
    input, _ := json.Marshal(body)
    wrapper.handler = &ApiGatewaySpy[string]{}
    var out bytes.Buffer
    ctx := &httpContext{
        header: http.Header{
            "Content-Type":        []string{"application/json"},
            "Fn-Http-Request-Url": []string{"/v1?param1=value%20with%20spaces"},
            "Fn-Http-Method":      []string{"POST"},
        },
        config: map[string]string{},
    }
    wrapper.Serve(fdk.WithContext(context.Background(), ctx), bytes.NewReader(input), &out)
    var result string
    _ = json.Unmarshal(out.Bytes(), &result)
    return result
}

func wrapperHandlerNonHttpContext[T any](body T) string {
    input, _ := json.Marshal(body)
    var capturedErr string
    wrapper := &APIGatewayHandleWrapper[T]{
        bodyType: reflect.TypeOf(""),
        fatal: func(args ...interface{}) {
            capturedErr = fmt.Sprint(args...)
        },
    }
    wrapper.handler = &ApiGatewaySpy[T]{}
    var out bytes.Buffer
    ctx := &otherContext{
        header: http.Header{},
        config: map[string]string{},
    }
    wrapper.Serve(fdk.WithContext(context.Background(), ctx), bytes.NewReader(input), &out)
    return capturedErr
}

func TestHandleRequestWithStringRequestBodyShouldReturnValidResponse(t *testing.T) {
    requestBody := "input data"
    wrapper := &APIGatewayHandleWrapper[string]{
        bodyType: reflect.TypeOf(""),
    }
    responseBody := stringWrapperHandler(&requestBody, wrapper)
    assert.Equal(t, "input data", responseBody)
}

func TestHandleResponse_WithHeaders(t *testing.T) {
    calledHeaders := map[string][]string{}
    addHeaderFunc = func(out io.Writer, k, v string) {
        calledHeaders[k] = append(calledHeaders[k], v)
    }
    defer func() { addHeaderFunc = fdk.AddHeader }()
    requestBody := "input data"
    wrapper := &APIGatewayHandleWrapper[string]{
        bodyType: reflect.TypeOf(""),
    }
    stringWrapperHandler(&requestBody, wrapper)
    assert.Equal(t, 1, len(calledHeaders))
    assert.Equal(t, "abc", calledHeaders["X-Test"][0])
}

func TestRequestAttributesArePassedToServe(t *testing.T) {
    spy := &ApiGatewaySpy[string]{}
    wrapper := &APIGatewayHandleWrapper[string]{
        handler:  spy,
        bodyType: reflect.TypeOf(""),
    }
    ctx := &httpContext{
        header: http.Header{
            "Content-Type":        []string{"application/json"},
            "Fn-Http-Request-Url": []string{"/v1/hello?id=123"},
            "Fn-Http-Method":      []string{"POST"},
            "Foo":                 []string{"bar", "hello"},
        },
        config: map[string]string{},
    }
    var out bytes.Buffer
    wrapper.Serve(fdk.WithContext(context.Background(), ctx), bytes.NewReader([]byte("")), &out)

    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, "POST", spy.LastRequest.Method)
    assert.Equal(t, "/v1/hello?id=123", spy.LastRequest.RequestURL)
    assert.Equal(t, "123", spy.LastRequest.QueryParameters.Get("id"))
    assert.Equal(t, "application/json", spy.LastRequest.Headers.Get("Content-Type"))
    assert.Equal(t, []string{"bar", "hello"}, spy.LastRequest.Headers.Values("Foo"))
}

func TestRequestNoQueryParameters(t *testing.T) {
    spy := &ApiGatewaySpy[string]{}
    wrapper := &APIGatewayHandleWrapper[string]{
        handler:  spy,
        bodyType: reflect.TypeOf(""),
    }
    ctx := &httpContext{
        header: http.Header{
            "Content-Type":        []string{"application/json"},
            "Fn-Http-Request-Url": []string{"/v1/hello"},
            "Fn-Http-Method":      []string{"POST"},
            "Foo":                 []string{"bar", "hello"},
        },
        config: map[string]string{},
    }
    var out bytes.Buffer
    wrapper.Serve(fdk.WithContext(context.Background(), ctx), bytes.NewReader([]byte("")), &out)

    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, "/v1/hello", spy.LastRequest.RequestURL)
    assert.NotNil(t, spy.LastRequest.QueryParameters)
}

func TestRequestDuplicateQueryParametersShouldUseFirst(t *testing.T) {
    spy := &ApiGatewaySpy[string]{}
    wrapper := &APIGatewayHandleWrapper[string]{
        handler:  spy,
        bodyType: reflect.TypeOf(""),
    }
    ctx := &httpContext{
        header: http.Header{
            "Content-Type":        []string{"application/json"},
            "Fn-Http-Request-Url": []string{"/v1/hello?id=123&id=321"},
            "Fn-Http-Method":      []string{"POST"},
        },
        config: map[string]string{},
    }
    var out bytes.Buffer
    wrapper.Serve(fdk.WithContext(context.Background(), ctx), bytes.NewReader([]byte("")), &out)

    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, "/v1/hello?id=123&id=321", spy.LastRequest.RequestURL)
    assert.Equal(t, "123", spy.LastRequest.QueryParameters.Get("id"))
}

func TestRequestCustomHeaders(t *testing.T) {
    spy := &ApiGatewaySpy[string]{}
    wrapper := &APIGatewayHandleWrapper[string]{
        handler:  spy,
        bodyType: reflect.TypeOf(""),
        fatal: func(args ...interface{}) {
            log.Fatal(args...)
        },
    }
    ctx := &httpContext{
        header: http.Header{
            "Fn-Intent":           []string{"httprequest"},
            "Content-Type":        []string{"application/json"},
            "Fn-Http-Request-Url": []string{"/v1/hello"},
            "Fn-Http-Method":      []string{"POST"},
            "Fn-Http-H-":          []string{"ignored"},
            "Fn-Http-H-Foo":       []string{"c"},
            "Fn-Http-H-Bar":       []string{"a", "b"},
        },
        config: map[string]string{},
    }
    var out bytes.Buffer
    wrapper.Serve(fdk.WithContext(context.Background(), ctx), bytes.NewReader([]byte("")), &out)

    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, "c", spy.LastRequest.Headers.Get("Fn-Http-H-Foo"))
    assert.Equal(t, []string{"a", "b"}, spy.LastRequest.Headers.Values("Fn-Http-H-Bar"))
}

func TestHandleRequest_WithStructRequestBodyShouldReturnValidResponse(t *testing.T) {
    requestBody := Employee{Name: "foo"}
    responseBody := structWrapperHandler(requestBody)
    assert.Equal(t, requestBody, responseBody)
}

/*
This test is equivalent to calling the Function directly, thus bypassing the OCI API Gateway.
*/
func TestHandleRequest_WithInvalidRequestShouldReturnError(t *testing.T) {
    requestBody := "input data"
    fatalLog := wrapperHandlerNonHttpContext(&requestBody)
    assert.Equal(t, "failed to read HTTP Context. Ensure the Function is invoked via OCI API Gateway and try again", fatalLog)
}

func TestReadFnRequestString(t *testing.T) {
    wrapper := &APIGatewayHandleWrapper[string]{
        bodyType: reflect.TypeOf(""),
    }
    body, err := wrapper.readFnInput(bytes.NewReader([]byte("a plain old string, input")))
    assert.Nil(t, err)
    assert.Equal(t, "a plain old string, input", body)
}

func TestReadFnRequestEmptyString(t *testing.T) {
    wrapper := &APIGatewayHandleWrapper[string]{
        bodyType: reflect.TypeOf(""),
    }
    body, err := wrapper.readFnInput(bytes.NewReader([]byte{}))
    assert.Nil(t, err)
    assert.Empty(t, body)
}

func TestReadFnRequestStruct(t *testing.T) {
    wrapper := &APIGatewayHandleWrapper[Employee]{
        bodyType: reflect.TypeOf(Employee{}),
    }
    body, err := wrapper.readFnInput(bytes.NewReader([]byte(`{"name":"foo"}`)))
    assert.Nil(t, err)
    assert.Equal(t, Employee{Name: "foo"}, body)
}

func TestReadFnRequestStructPointer(t *testing.T) {
    wrapper := &APIGatewayHandleWrapper[Employee]{
        bodyType: reflect.TypeOf(&Employee{}),
    }
    body, err := wrapper.readFnInput(bytes.NewReader([]byte(`{"name":"foo"}`)))
    assert.Nil(t, err)
    assert.Equal(t, &Employee{Name: "foo"}, body)
}

func TestReadFnRequestEmptyStruct(t *testing.T) {
    wrapper := &APIGatewayHandleWrapper[Employee]{
        bodyType: reflect.TypeOf(Employee{}),
    }
    body, err := wrapper.readFnInput(bytes.NewReader([]byte{}))
    assert.Nil(t, err)
    assert.Equal(t, Employee{}, body)
}
