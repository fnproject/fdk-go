package fn_events

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/fnproject/fdk-go"
    "io"
    "log"
    "net/http"
    "net/url"
    "reflect"
)

// Testing shims
var addHeaderFunc = fdk.AddHeader

type ApiGatewayFnHandler[T any] interface {
    Serve(ctx context.Context, request APIGatewayRequestEvent[T]) *APIGatewayResponseEvent[any]
}

type APIGatewayHandleWrapper[T any] struct {
    handler  ApiGatewayFnHandler[T]
    bodyType reflect.Type
    fatal    func(args ...interface{})
}

// APIGatewayHandler Called by customer to inject their handler logic
func APIGatewayHandler[T any](fnHandler ApiGatewayFnHandler[T], bodyType reflect.Type) fdk.Handler {
    return &APIGatewayHandleWrapper[T]{
        handler:  fnHandler,
        bodyType: bodyType,
        fatal: func(args ...interface{}) {
            log.Fatal(args...)
        },
    }
}

func (wrapper *APIGatewayHandleWrapper[T]) Serve(ctx context.Context, in io.Reader, out io.Writer) {
    body, err := wrapper.readFnInput(in)
    if err != nil {
        wrapper.handleError(err)
        return
    }

    httpCtx, err := wrapper.fetchHttpContext(ctx)
    if err != nil {
        wrapper.handleError(err)
        return // Return is needed for test because of overridden fatal os exit
    }
    queryParameters, err := url.Parse(httpCtx.RequestURL())
    if err != nil {
        wrapper.handleError(err)
        return
    }
    request := APIGatewayRequestEvent[T]{
        Body:            body.(T),
        Method:          httpCtx.RequestMethod(),
        RequestURL:      httpCtx.RequestURL(),
        Headers:         httpCtx.Header(),
        QueryParameters: queryParameters.Query(),
    }

    response := wrapper.handler.Serve(ctx, request)

    if response != nil {
        for k, vs := range response.Headers {
            for _, v := range vs {
                addHeaderFunc(out, k, v)
            }
        }

        fdk.WriteStatus(out, response.StatusCode)
        log.Printf("Writing body %v", response.Body)
        if err := wrapper.writeBody(response.Body, out); err != nil {
            wrapper.handleError(err)
            return
        }
        return
    }
}

func (wrapper *APIGatewayHandleWrapper[T]) writeBody(body any, out io.Writer) error {
    if str, ok := body.(string); ok {
        _, err := out.Write([]byte(str))
        if err != nil {
            return err
        }
    } else {
        err := json.NewEncoder(out).Encode(body)
        if err != nil {
            return err
        }
    }
    return nil
}

func (wrapper *APIGatewayHandleWrapper[T]) handleError(err error) {
    log.SetFlags(0) // Hides redundant timestamp in log content
    wrapper.fatal(err)
}

func (wrapper *APIGatewayHandleWrapper[T]) fetchHttpContext(ctx context.Context) (fdk.HTTPContext, error) {
    httpContext, ok := fdk.GetContext(ctx).(fdk.HTTPContext)
    if !ok {
        return nil, fmt.Errorf("failed to read HTTP Context. Ensure the Function is invoked via OCI API Gateway and try again")
    }

    return httpContext, nil
}

func (wrapper *APIGatewayHandleWrapper[T]) readFnInput(in io.Reader) (any, error) {
    data, err := io.ReadAll(in)
    if err != nil {
        return nil, fmt.Errorf("failed to read input from Reader %v", err)
    }
    // Handle empty-input case (zero value for the expected type)
    if len(data) == 0 {
        return reflect.Zero(wrapper.bodyType).Interface(), nil
    }

    kind := wrapper.bodyType.Kind()
    switch kind {
    case reflect.String:
        return string(data), nil

    case reflect.Ptr:
        // bodyType is a pointer to struct or map
        ptrVal := reflect.New(wrapper.bodyType.Elem())
        if err := json.Unmarshal(data, ptrVal.Interface()); err != nil {
            return nil, fmt.Errorf("failed to coerce event to user function parameter type %w", err)
        }
        return ptrVal.Interface(), nil

    case reflect.Struct, reflect.Map:
        val := reflect.New(wrapper.bodyType) // Always a pointer to value
        if err := json.Unmarshal(data, val.Interface()); err != nil {
            return nil, fmt.Errorf("failed to coerce event to user function parameter type %w", err)
        }
        return val.Elem().Interface(), nil

    default:
        return data, nil
    }
}


type APIGatewayRequestEvent[T any] struct {
    Body            T
    Method          string
    RequestURL      string
    Headers         http.Header
    QueryParameters url.Values
}

type APIGatewayResponseEvent[T any] struct {
    Body       T       `json:"body"`
    StatusCode int     `json:"statusCode"`
    Headers    Headers `json:"headers"`
}

type APIGatewayResponseEventBuilder[T any] struct {
    body       T
    statusCode int
    headers    Headers
}

func NewAPIGatewayResponseEventBuilder[T any]() *APIGatewayResponseEventBuilder[T] {
    return &APIGatewayResponseEventBuilder[T]{}
}

func (b *APIGatewayResponseEventBuilder[T]) Body(body T) *APIGatewayResponseEventBuilder[T] {
    b.body = body
    return b
}

func (b *APIGatewayResponseEventBuilder[T]) StatusCode(code int) *APIGatewayResponseEventBuilder[T] {
    b.statusCode = code
    return b
}

func (b *APIGatewayResponseEventBuilder[T]) Headers(headers Headers) *APIGatewayResponseEventBuilder[T] {
    b.headers = headers
    return b
}

func (b *APIGatewayResponseEventBuilder[T]) Build() *APIGatewayResponseEvent[T] {
    return &APIGatewayResponseEvent[T]{
        Body:       b.body,
        StatusCode: b.statusCode,
        Headers:    b.headers,
    }
}
