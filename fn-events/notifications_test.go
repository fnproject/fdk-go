package fn_events

import (
    "bytes"
    "context"
    "encoding/json"
    "github.com/stretchr/testify/assert"
    "net/http"
    "reflect"
    "testing"
    "github.com/fnproject/fdk-go"
)

type NotificationsSpy[T any] struct {
    CallCount   int
    LastCtx     context.Context
    LastRequest NotificationMessage[T]
}

func (s *NotificationsSpy[T]) Serve(ctx context.Context, batch NotificationMessage[T]) {
    s.CallCount++
    s.LastCtx = ctx
    s.LastRequest = batch
}

func notificationsWrapperHandler[T any](body string, wrapper *NotificationHandleWrapper[T]) {
    var out bytes.Buffer
    ctx := &httpContext{
        config: map[string]string{},
        header: map[string][]string{
            "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
        },
    }
    wrapper.Serve(fdk.WithContext(context.Background(), ctx), bytes.NewReader([]byte(body)), &out)
}

func TestHandleNotificationStringIsParsedCorrectly(t *testing.T) {
    requestBody := "a plain string"
    spy := &NotificationsSpy[string]{}
    wrapper := &NotificationHandleWrapper[string]{
        handler:  spy,
        bodyType: reflect.TypeOf(""),
    }
    notificationsWrapperHandler(requestBody, wrapper)

    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, requestBody, spy.LastRequest.Content)
    assert.Equal(t, http.Header{
        "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
    }, spy.LastRequest.Headers)
}

func TestHandleNotificationStringEmpty(t *testing.T) {
    requestBody := ""
    spy := &NotificationsSpy[string]{}
    wrapper := &NotificationHandleWrapper[string]{
        handler:  spy,
        bodyType: reflect.TypeOf(""),
    }
    notificationsWrapperHandler(requestBody, wrapper)

    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, requestBody, spy.LastRequest.Content)
    assert.Equal(t, http.Header{
        "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
    }, spy.LastRequest.Headers)
}

func TestHandleNotificationStruct(t *testing.T) {
    employee := Employee{
        Name: "foo",
    }
    spy := &NotificationsSpy[Employee]{}
    requestBody, _ := json.Marshal(employee)
    wrapper := &NotificationHandleWrapper[Employee]{
        handler:  spy,
        bodyType: reflect.TypeOf(Employee{}),
    }
    notificationsWrapperHandler(string(requestBody), wrapper)

    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, employee, spy.LastRequest.Content)
    assert.Equal(t, http.Header{
        "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
    }, spy.LastRequest.Headers)
}

func TestHandleNotificationStructPointer(t *testing.T) {
    employee := &Employee{
        Name: "foo",
    }
    spy := &NotificationsSpy[*Employee]{}
    requestBody, _ := json.Marshal(employee)
    wrapper := &NotificationHandleWrapper[*Employee]{
        handler:  spy,
        bodyType: reflect.TypeOf(&Employee{}),
    }
    notificationsWrapperHandler(string(requestBody), wrapper)

    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, employee, spy.LastRequest.Content)
    assert.Equal(t, http.Header{
        "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
    }, spy.LastRequest.Headers)
}
