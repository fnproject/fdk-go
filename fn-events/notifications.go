package fn_events

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/fnproject/fdk-go"
    "io"
    "log"
    "net/http"
    "reflect"
)

type NotificationHandleWrapper[T any] struct {
    handler  NotificationsFnHandler[T]
    bodyType reflect.Type
    FatalFuncWrapper
}

type NotificationsFnHandler[T any] interface {
    Serve(ctx context.Context, notification NotificationMessage[T])
}

// NotificationsHandler returns an fdk.Handler that wraps a NotificationsFnHandler.
//
// Parameters:
//   - fnHandler: The user-supplied handler function, which will be called with each
//     NotificationMessage instance. Its generic type parameter (NotificationsFnHandler[any])
//     indicates that it accepts any type but you should instantiate this handler for your
//     customer-supplied struct type for strong typing and automatic JSON deserialization.
//   - sourceType: A reflect.Type representing the customer struct type you wish JSON input
//     to be automatically deserialized into. Typically, pass reflect.TypeOf(YourStruct{}) here.
//
// The NotificationHandleWrapper embeds a FatalFuncWrapper, which allows you to inject a custom
// fatal function for error handling and process termination behavior. By default, this uses log.Fatal.
// During unit testing, you can override the FatalFuncWrapper's behavior to prevent test interruption
// or to assert on fatal error occurrences.
//
// Example usage:
//
//    type MyCustomPayload struct {
//        Name string `json:"name"`
//        Age  int    `json:"age"`
//    }
//
//    func MyNotificationHandler(ctx context.Context, msg NotificationMessage[any]) error {
//        // Handle message.Content and message.Headers here
//        ...
//    }
//
//    handler := NotificationsHandler(
//        MyNotificationHandler,
//        reflect.TypeOf(MyCustomPayload{}),
//    )
//
func NotificationsHandler[T any](fnHandler NotificationsFnHandler[T], bodyType reflect.Type) fdk.Handler {
    return &NotificationHandleWrapper[T]{
        handler:  fnHandler,
        bodyType: bodyType,
        FatalFuncWrapper: FatalFuncWrapper{
            fatal: func(args ...interface{}) {
                log.Fatal(args...)
            },
        },
    }
}

func (wrapper *NotificationHandleWrapper[T]) Serve(ctx context.Context, in io.Reader, out io.Writer) {
    content, err := wrapper.readFnInput(in)
    if err != nil {
        wrapper.HandleError(err)
        return
    }

    fdkContext := FetchContext(ctx)

    notification := NotificationMessage[T]{
        Content: content.(T),
        Headers: fdkContext.Header(),
    }
    wrapper.handler.Serve(ctx, notification)
}

func (wrapper *NotificationHandleWrapper[T]) readFnInput(in io.Reader) (any, error) {
    data, err := io.ReadAll(in)
    if err != nil {
        return nil, fmt.Errorf("failed to read input from Reader %v", err)
    }
    if len(data) == 0 {
        return reflect.Zero(wrapper.bodyType).Interface(), nil
    }

    kind := wrapper.bodyType.Kind()
    switch kind {
    case reflect.String:
        return string(data), nil

    case reflect.Ptr:
        // bodyType is a pointer to a struct
        ptrVal := reflect.New(wrapper.bodyType.Elem())
        if err := json.Unmarshal(data, ptrVal.Interface()); err != nil {
            return nil, fmt.Errorf("failed to coerce event to user function parameter type %v, %v", wrapper.bodyType, err)
        }
        return ptrVal.Interface(), nil

    case reflect.Struct, reflect.Map:
        val := reflect.New(wrapper.bodyType) // *T
        if err := json.Unmarshal(data, val.Interface()); err != nil {
            return nil, fmt.Errorf("failed to coerce event to user function parameter type %v, %v", wrapper.bodyType, err)
        }
        return val.Elem().Interface(), nil // dereference to value type

    default:
        return data, nil
    }
}

type NotificationMessage[T any] struct {
    Content T
    Headers http.Header
}
