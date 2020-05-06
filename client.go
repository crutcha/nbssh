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
	Count    uint                     `json:"count"`
	Next     string                   `json:"next"`
	Previous string                   `json:"previous"`
    // This is kinda dirty but we only want a single key/value pair out of this
    // this map, it's not worth the effort to define all the structs that could
    // be defined here
	Results  []map[string]interface{} `json:results`
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
    if *site != "" {
        q.Add("site", *site)
    }
    if *tenant != "" {
        q.Add("tenant", *tenant)
    }
    if *role != "" {
        q.Add("role", *role)
    }
	q.Add("limit", strconv.Itoa(pageSize))

	deviceArray := make([]string, 0)
	hasMoreResults := true
	currentOffset := 0
	for hasMoreResults == true {
		var payload DRFResponse

        q.Set("offset", strconv.Itoa(currentOffset))
		req.URL.RawQuery = q.Encode()
		resp, requestErr := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &payload)

		if requestErr != nil {
			panic(requestErr)
		}

		for _, device := range payload.Results {
			if device["name"] != nil {
				deviceArray = append(deviceArray, device["name"].(string))

			}
		}

		if payload.Next != "" {
            currentOffset += pageSize
		} else {
            hasMoreResults = false
		}
	}

	return deviceArray
}
