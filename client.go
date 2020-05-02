package main

import (
	"context"
	"fmt"
	runtimeclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"os"
)

const pageSize int64 = 50

var apiToken string
var netboxHost string

func init() {
	apiToken = os.Getenv("NETBOX_API_TOKEN")
	netboxHost = os.Getenv("NETBOX_HOST")

	if apiToken == "" {
		fmt.Println("NETBOX_API_TOKEN environment variable not set but required")
		os.Exit(1)
	}
	if netboxHost == "" {
		fmt.Println("NETBOX_HOST environment variable not set but required")
		os.Exit(1)
	}
}

func newNetboxClient() *client.NetBox {
	t := runtimeclient.New(netboxHost, client.DefaultBasePath, client.DefaultSchemes)
	authHeader := "Token " + apiToken
	t.DefaultAuthentication = runtimeclient.APIKeyAuth("Authorization", "header", authHeader)

	if *verbose {
		t.SetDebug(true)
	}

	client := client.New(t, strfmt.Default)
	return client
}

// Wrapper around DcimDevicesList to deal with pagination
func queryDevices(client *client.NetBox) []string {
	currentOffset := int64(0)

	// cannot address pointer of const apparently
	ourLimit := pageSize

	params := dcim.DcimDevicesListParams{
		Context:      context.Background(),
		Site:         site,
		Tenant:       tenant,
		Role:         role,
		Manufacturer: manufacturer,
		Status:       status,
		Limit:        &ourLimit,
		Offset:       &currentOffset,
	}

	deviceArray := make([]string, 0)
	hasMoreResults := true
	for hasMoreResults == true {
		response, requestErr := client.Dcim.DcimDevicesList(&params, nil)

		if requestErr != nil {
			panic(requestErr)
		}

		for _, device := range response.Payload.Results {
			if device.Name != nil {
				deviceArray = append(deviceArray, *device.Name)

			}
		}

		if response.Payload.Next != nil {
			*params.Offset += ourLimit
		} else {
			hasMoreResults = false
		}
	}

	return deviceArray
}
