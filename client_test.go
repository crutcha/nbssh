package main

import (
	"testing"
    "sort"

	"github.com/stretchr/testify/assert"
    "github.com/jarcoal/httpmock"
)

func TestClientGetDevices(t *testing.T) {
	assert := assert.New(t)

    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    
    test_payload := `
    {
        "count": 3,
        "next": null,
        "previous": null,
        "results": [
            {
                "name": "test-device-1"
            },
            {
                "name": "test-device-1"
            },
            {
                "name": "test-device-2"
            }
        ]
    }
    `
    mock_response := httpmock.NewStringResponder(200, test_payload)
    httpmock.RegisterResponder("GET", "http://localhost/api/dcim/devices/", mock_response)

    devices := queryDevices()

    // since hash type is used during query, no guarentee it will be ordered
    sort.Slice(devices, func(i, j int) bool { return devices[i] < devices[j] })

    assert.Equal(2, len(devices), "2 devices are returned after de-dup")
    assert.Equal("test-device-1", devices[0], "test-device-1 in return value")
    assert.Equal("test-device-2", devices[1], "test-device-1 in return value")
}
