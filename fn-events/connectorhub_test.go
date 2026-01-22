package fn_events

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "reflect"
    "testing"
    "time"

    "github.com/fnproject/fdk-go"
    "github.com/stretchr/testify/assert"
)

type ConnectorHubSpy[T any] struct {
    CallCount   int
    LastCtx     context.Context
    LastRequest ConnectorHubBatch[T]
}

func (s *ConnectorHubSpy[T]) Serve(ctx context.Context, batch ConnectorHubBatch[T]) {
    s.CallCount++
    s.LastCtx = ctx
    s.LastRequest = batch
}

func structConnectorHubWrapperHandler[T any](body string, wrapper *ConnectorHubHandleWrapper[T]) {
    var out bytes.Buffer
    ctx := &httpContext{
        config: map[string]string{},
        header: map[string][]string{
            "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
        },
    }
    wrapper.Serve(fdk.WithContext(context.Background(), ctx), bytes.NewReader([]byte(body)), &out)
}

func TestHandleLoggingDataIsParsedCorrectly(t *testing.T) {
    expectedTime, _ := time.Parse(time.RFC3339, "2025-10-24T15:06:17.000Z")
    expected := LoggingData{
        ID:          "abc-zyx",
        Source:      "your-log",
        SpecVersion: "1.0",
        Subject:     "schedule",
        Type:        "com.oraclecloud.functions.application.functioninvoke",
        Data: map[string]string{
            "applicationId": "ocid1.fnapp.oc1.abc",
            "containerId":   "n/a",
            "functionId":    "ocid1.fnfunc.oc1.abc",
            "message":       "Received function invocation request",
            "opcRequestId":  "/abc/def",
            "requestId":     "/def/abc",
            "src":           "stdout",
        },
        Oracle: map[string]string{
            "compartmentid": "ocid1.tenancy.oc1..xyz",
            "ingestedtime":  "2025-10-23T15:45:19.457Z",
            "loggroupid":    "ocid1.loggroup.oc1.abc",
            "logid":         "ocid1.log.oc1.def",
            "tenantid":      "ocid1.tenancy.oc1..xyz",
        },
        Time: expectedTime,
    }
    requestBody := `[{
        "data": {
            "applicationId": "ocid1.fnapp.oc1.abc",
            "containerId": "n/a",
            "functionId": "ocid1.fnfunc.oc1.abc",
            "message": "Received function invocation request",
            "opcRequestId": "/abc/def",
            "requestId": "/def/abc",
            "src": "stdout"
        },
        "id": "abc-zyx",
        "oracle": {
            "compartmentid": "ocid1.tenancy.oc1..xyz",
            "ingestedtime": "2025-10-23T15:45:19.457Z",
            "loggroupid": "ocid1.loggroup.oc1.abc",
            "logid": "ocid1.log.oc1.def",
            "tenantid": "ocid1.tenancy.oc1..xyz"
        },
        "source": "your-log",
        "specversion": "1.0",
        "subject": "schedule",
        "time": "2025-10-24T15:06:17.000Z",
        "type": "com.oraclecloud.functions.application.functioninvoke"
    }, {
        "data": {
            "applicationId": "ocid1.fnapp.oc1.def",
            "containerId": "n/a",
            "functionId": "ocid1.fnfunc.oc1.def",
            "message": "Received function invocation request",
            "opcRequestId": "/def/xyz",
            "requestId": "/foo/bar",
            "src": "stdout"
        },
        "id": "foo-zyx",
        "oracle": {
            "compartmentid": "ocid1.tenancy.oc1..xyz",
            "ingestedtime": "2025-11-23T15:45:19.457Z",
            "loggroupid": "ocid1.loggroup.oc1.def",
            "logid": "ocid1.log.oc1.xyz",
            "tenantid": "ocid1.tenancy.oc1..xyz"
        },
        "source": "your-log",
        "specversion": "1.0",
        "subject": "schedule",
        "time": "2025-11-23T15:45:17.239Z",
        "type": "com.oraclecloud.functions.application.functioninvoke"
    }]`
    spy := &ConnectorHubSpy[LoggingData]{}
    wrapper := &ConnectorHubHandleWrapper[LoggingData]{
        handler:  spy,
        bodyType: reflect.TypeOf(LoggingData{}),
    }
    structConnectorHubWrapperHandler(requestBody, wrapper)
    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, 2, len(spy.LastRequest.Batch))
    assert.Equal(t, expected, spy.LastRequest.Batch[0])
    assert.Equal(t, http.Header{
        "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
    }, spy.LastRequest.Headers)
}

func TestHandleLoggingDataFailsWhenNotListFromConnectorHub(t *testing.T) {
    requestBody := `{
        "data": {
            "applicationId": "ocid1.fnapp.oc1.abc",
            "containerId": "n/a",
            "functionId": "ocid1.fnfunc.oc1.abc",
            "message": "Received function invocation request",
            "opcRequestId": "/abc/def",
            "requestId": "/def/abc",
            "src": "stdout"
        },
        "id": "abc-zyx",
        "oracle": {
            "compartmentid": "ocid1.tenancy.oc1..xyz",
            "ingestedtime": "2025-10-23T15:45:19.457Z",
            "loggroupid": "ocid1.loggroup.oc1.abc",
            "logid": "ocid1.log.oc1.def",
            "tenantid": "ocid1.tenancy.oc1..xyz"
        },
        "source": "your-log",
        "specversion": "1.0",
        "subject": "schedule",
        "time": "2025-10-24T15:06:17.000Z",
        "type": "com.oraclecloud.functions.application.functioninvoke"
    }`
    var capturedErr string
    spy := &ConnectorHubSpy[LoggingData]{}
    wrapper := &ConnectorHubHandleWrapper[LoggingData]{
        handler:  spy,
        bodyType: reflect.TypeOf(LoggingData{}),
        FatalFuncWrapper: FatalFuncWrapper{
            fatal: func(args ...interface{}) {
                capturedErr = fmt.Sprint(args...)
            },
        },
    }
    structConnectorHubWrapperHandler(requestBody, wrapper)
    expectedErr := fmt.Sprintf("failed to unmarshal to slice of fn_events.LoggingData, %s: json: cannot unmarshal object into Go value of type []fn_events.LoggingData", requestBody)
    assert.Equal(t, 0, spy.CallCount)
    assert.Equal(t, expectedErr, capturedErr)
}

func TestHandleLoggingDataSucceedsWithEmptyListFromConnectorHub(t *testing.T) {
    requestBody := "[]"
    var errors int
    spy := &ConnectorHubSpy[LoggingData]{}
    wrapper := &ConnectorHubHandleWrapper[LoggingData]{
        handler:  spy,
        bodyType: reflect.TypeOf(LoggingData{}),
        FatalFuncWrapper: FatalFuncWrapper{
            fatal: func(args ...interface{}) {
                errors++
            },
        },
    }
    structConnectorHubWrapperHandler(requestBody, wrapper)
    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, 0, errors)
}

func TestHandleMetricDataIsParsedCorrectly(t *testing.T) {
    expected := MetricData{
        Namespace:     "oci_objectstorage",
        ResourceGroup: "nullable",
        CompartmentId: "ocid1.tenancy.oc1..xyz",
        Name:          "PutRequests",
        Dimensions: map[string]string{
            "resourceID":          "ocid1.bucket.oc1.uk-london-1.xyz",
            "resourceDisplayName": "foo",
        },
        Metadata: map[string]string{
            "displayName": "PutObject Request Count",
            "unit":        "count",
        },
        Datapoints: []Datapoint{
            {
                Timestamp: 1761318377414,
                Value:     float64Ptr(1.0),
                Count:     intPtr(1),
            },
        },
    }
    requestBody := `[{
        "namespace": "oci_objectstorage",
        "resourceGroup": "nullable",
        "compartmentId": "ocid1.tenancy.oc1..xyz",
        "name": "PutRequests",
        "dimensions": {
            "resourceID": "ocid1.bucket.oc1.uk-london-1.xyz",
            "resourceDisplayName": "foo"
        },
        "metadata": {
            "displayName": "PutObject Request Count",
            "unit": "count"
        },
        "datapoints": [
            {
                "timestamp": 1761318377414,
                "value": 1.0,
                "count": 1
            }
        ]
    }, {
        "namespace": "oci_objectstorage",
        "resourceGroup": null,
        "compartmentId": "ocid1.tenancy.oc1..abc",
        "name": "PutRequests",
        "dimensions": {
            "resourceID": "ocid1.bucket.oc1.uk-london-1.abc",
            "resourceDisplayName": "bar"
        },
        "metadata": {
            "displayName": "PutObject Request Count",
            "unit": "count"
        },
        "datapoints": [
            {
                "timestamp": 1761318377414,
                "value": 1.0,
                "count": 1
            },
            {
                "timestamp": 1761318377614,
                "value": 2.0,
                "count": 1
            }
        ]
    }]`
    spy := &ConnectorHubSpy[MetricData]{}
    wrapper := &ConnectorHubHandleWrapper[MetricData]{
        handler:  spy,
        bodyType: reflect.TypeOf(MetricData{}),
    }
    structConnectorHubWrapperHandler(requestBody, wrapper)
    assert.Equal(t, 1, spy.CallCount)
    metric := spy.LastRequest.Batch[0]
    assert.Equal(t, expected, metric)
    assert.Equal(t, time.UnixMilli(expected.Datapoints[0].Timestamp), metric.Datapoints[0].Time())
    assert.Equal(t, http.Header{
        "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
    }, spy.LastRequest.Headers)
}

func TestHandleMetricDataSucceedsWithEmptyListFromConnectorHub(t *testing.T) {
    requestBody := "[]"
    var errors int
    spy := &ConnectorHubSpy[MetricData]{}
    wrapper := &ConnectorHubHandleWrapper[MetricData]{
        handler:  spy,
        bodyType: reflect.TypeOf(MetricData{}),
        FatalFuncWrapper: FatalFuncWrapper{
            fatal: func(args ...interface{}) {
                errors++
            },
        },
    }
    structConnectorHubWrapperHandler(requestBody, wrapper)
    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, 0, len(spy.LastRequest.Batch))
    assert.Equal(t, 0, errors)
    assert.Equal(t, http.Header{
        "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
    }, spy.LastRequest.Headers)
}

func TestHandleQueueSucceedsWithStructFromConnectorHub(t *testing.T) {
    employee := Employee{Name: "a string name"}
    employeeJSON, _ := json.Marshal(employee)

    requestBody := fmt.Sprintf("[%s]", employeeJSON)
    var errors int
    spy := &ConnectorHubSpy[Employee]{}
    wrapper := &ConnectorHubHandleWrapper[Employee]{
        handler:  spy,
        bodyType: reflect.TypeOf(Employee{}),
        FatalFuncWrapper: FatalFuncWrapper{
            fatal: func(args ...interface{}) {
                errors++
            },
        },
    }
    structConnectorHubWrapperHandler(requestBody, wrapper)
    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, 1, len(spy.LastRequest.Batch))
    assert.Equal(t, employee, spy.LastRequest.Batch[0])
    assert.Equal(t, 0, errors)
    assert.Equal(t, http.Header{
        "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
    }, spy.LastRequest.Headers)
}

func TestHandleQueueSucceedsWithStringFromConnectorHub(t *testing.T) {
    requestBody := `["a plain string"]`
    var errors int
    spy := &ConnectorHubSpy[string]{}
    wrapper := &ConnectorHubHandleWrapper[string]{
        handler:  spy,
        bodyType: reflect.TypeOf(""),
        FatalFuncWrapper: FatalFuncWrapper{
            fatal: func(args ...interface{}) {
                errors++
            },
        },
    }
    structConnectorHubWrapperHandler(requestBody, wrapper)
    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, 1, len(spy.LastRequest.Batch))
    assert.Equal(t, "a plain string", spy.LastRequest.Batch[0])
    assert.Equal(t, 0, errors)
    assert.Equal(t, http.Header{
        "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
    }, spy.LastRequest.Headers)
}

func TestHandleQueueSucceedsWithEmptyListFromConnectorHub(t *testing.T) {
    requestBody := "[]"
    var errors int
    spy := &ConnectorHubSpy[MetricData]{}
    wrapper := &ConnectorHubHandleWrapper[MetricData]{
        handler:  spy,
        bodyType: reflect.TypeOf(MetricData{}),
        FatalFuncWrapper: FatalFuncWrapper{
            fatal: func(args ...interface{}) {
                errors++
            },
        },
    }
    structConnectorHubWrapperHandler(requestBody, wrapper)
    assert.Equal(t, 1, spy.CallCount)
    assert.Equal(t, 0, len(spy.LastRequest.Batch))
    assert.Equal(t, 0, errors)
    assert.Equal(t, http.Header{
        "Oci-Subject-Tenancy-Id": {"ocid1.tenancy.oc1..abc"},
    }, spy.LastRequest.Headers)
}

func float64Ptr(f float64) *float64 { return &f }
func intPtr(i int) *int             { return &i }
