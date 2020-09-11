package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const pageSize int = 50

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

type DRFResponse struct {
	Count    uint   `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	// This is kinda dirty but we only want a single key/value pair out of this
	// this map, it's not worth the effort to define all the structs that could
	// be defined here
	Results []map[string]interface{} `json:results`
}

// Tried using go-netbox here but if the NetBox instance is using HTTPS, it'll try HTTP
// first, receive a 301 redirect, and the redirect is followed but will not include the
// auth header when following redirect. Since this is very focused on a single endpoint
// with a limited query param set, we'll just handle this ourselves.
func queryDevices() []string {
	client := &http.Client{Timeout: time.Second * 10}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/dcim/devices/", netboxHost), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", apiToken))

	q := url.Values{}
	if *status != "" {
		q.Add("status", *status)
	}
	for _, value := range *site {
		q.Add("site", value)
	}
	for _, value := range *tenant {
		q.Add("tenant", value)
	}
	for _, value := range *role {
		q.Add("role", value)
	}
	for _, value := range *customfield {
		q.Add(fmt.Sprintf("cf_%s", value.Key), value.Value)
	}
	q.Add("limit", strconv.Itoa(pageSize))

	// Allows user to use custom field to determine FQDN of device
	usingCustomField := false
	if *namefield != "" {
		usingCustomField = true
	}

	deviceSet := make(map[string]bool)
	hasMoreResults := true
	currentOffset := 0
	for hasMoreResults == true {
		var payload DRFResponse

		q.Set("offset", strconv.Itoa(currentOffset))
		req.URL.RawQuery = q.Encode()
		resp, requestErr := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)

		if requestErr != nil {
			panic(requestErr)
		}

		if resp.StatusCode != 200 {
			panic(fmt.Errorf("%s %s", resp.Status, body))
		}

		json.Unmarshal(body, &payload)
		for _, device := range payload.Results {
			name := getDeviceString(device, usingCustomField)
			if _, ok := deviceSet[name]; !ok {
				deviceSet[name] = true
			}
		}

		if payload.Next != "" {
			currentOffset += pageSize
		} else {
			hasMoreResults = false
		}
	}

	deviceArray := make([]string, 0)
	for device, _ := range deviceSet {
		deviceArray = append(deviceArray, device)
	}

	return deviceArray
}

func getDeviceString(device map[string]interface{}, customField bool) string {
	if customField {
		custom_fields := device["custom_fields"].(map[string]interface{})
		if val, _ := custom_fields[*namefield]; val != nil {
			return val.(string)
		}
	}

	// name will always exist on device object
	name := device["name"].(string)

	return name
}
