package fn_events

import (
    "bytes"
    "context"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "net/http"
    "reflect"
    "testing"
    "time"

    "github.com/fnproject/fdk-go"
    "github.com/stretchr/testify/assert"
)

type ConnectorHubStreamingSpy[T any] struct {
    CallCount   int
    LastCtx     context.Context
    LastRequest ConnectorHubBatch[StreamingData[T]]
}

func (s *ConnectorHubStreamingSpy[T]) Serve(ctx context.Context, batch ConnectorHubBatch[StreamingData[T]]) {
    s.CallCount++
    s.LastCtx = ctx
    s.LastRequest = batch
}

func structConnectorHubStreamWrapperHandler[T any](body string, wrapper *ConnectorHubHandleStreamWrapper[T]) {
    var out bytes.Buffer
    ctx := &httpContext{
        config: map[string]string{},
        header: map[string][]string{
            "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
        },
    }
    wrapper.Serve(fdk.WithContext(context.Background(), ctx), bytes.NewReader([]byte(body)), &out)
}
func TestHandleStreamDataStringIsParsedCorrectly(t *testing.T) {
    valueBase64 := "U2VudCBhIHBsYWluIG1lc3NhZ2U="
    decodedValue, _ := base64.StdEncoding.DecodeString(valueBase64)
    expected := StreamingData[string]{
        Stream:    "stream-name",
        Partition: "0",
        Key:       nil,
        Value:     string(decodedValue),
        Offset:    3,
        Timestamp: 1761223385480,
    }
    requestBody := `[
    {"stream":"stream-name",
     "partition":"0",
     "key":null,
     "value":"U2VudCBhIHBsYWluIG1lc3NhZ2U=",
     "offset":3,
     "timestamp":1761223385480
    },
    {"stream":"stream-name",
     "partition":"0",
     "key":null,
     "value":"U2VudCBhIHBsYWluIG1lc3NhZ2U=",
     "offset":3,
     "timestamp":1761223385480
    }]`
    spy := &ConnectorHubStreamingSpy[string]{}
    wrapper := &ConnectorHubHandleStreamWrapper[string]{
        handler:   spy,
        valueType: reflect.TypeOf(""),
    }
    structConnectorHubStreamWrapperHandler(requestBody, wrapper)
    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, expected, spy.LastRequest.Batch[0])
    stream := spy.LastRequest.Batch[0]
    assert.Equal(t, time.UnixMilli(expected.Timestamp), stream.Time())
    assert.Equal(t, "Sent a plain message", stream.Value)
    assert.Equal(t, http.Header{
        "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
    }, spy.LastRequest.Headers)
}

func TestHandleStreamDataStructIsParsedCorrectly(t *testing.T) {
    employee := Employee{Name: "a string name"}
    employeeJSON, err := json.Marshal(employee)
    if err != nil {
        t.Fatalf("failed to marshal Employee: %v", err)
    }
    valueBase64 := base64.StdEncoding.EncodeToString(employeeJSON)
    expected := StreamingData[Employee]{
        Stream:    "stream-name",
        Partition: "0",
        Key:       nil,
        Value:     employee,
        Offset:    3,
        Timestamp: 1761223385480,
    }
    requestBodyTemplate := `[
    {"stream":"stream-name",
     "partition":"0",
     "key":null,
     "value":"%s",
     "offset":3,
     "timestamp":1761223385480
    },
    {"stream":"stream-name",
     "partition":"0",
     "key":null,
     "value":"%s",
     "offset":3,
     "timestamp":1761223385480
    }]`
    requestBody := fmt.Sprintf(requestBodyTemplate, valueBase64, valueBase64)
    spy := &ConnectorHubStreamingSpy[Employee]{}
    wrapper := &ConnectorHubHandleStreamWrapper[Employee]{
        handler:   spy,
        valueType: reflect.TypeOf(Employee{}),
    }
    structConnectorHubStreamWrapperHandler(requestBody, wrapper)

    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, expected, spy.LastRequest.Batch[0])
    stream := spy.LastRequest.Batch[0]
    assert.Equal(t, time.UnixMilli(expected.Timestamp), stream.Time())
    value := stream.Value
    assert.Equal(t, employee.Name, value.Name)
    assert.Equal(t, http.Header{
        "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
    }, spy.LastRequest.Headers)
}

func TestHandleStreamDataEmptyIsParsedCorrectly(t *testing.T) {
    requestBody := "[]"
    spy := &ConnectorHubStreamingSpy[Employee]{}
    wrapper := &ConnectorHubHandleStreamWrapper[Employee]{
        handler:   spy,
        valueType: reflect.TypeOf(Employee{}),
    }
    structConnectorHubStreamWrapperHandler(requestBody, wrapper)
    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, 0, len(spy.LastRequest.Batch))
    assert.Equal(t, http.Header{
        "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
    }, spy.LastRequest.Headers)
}
