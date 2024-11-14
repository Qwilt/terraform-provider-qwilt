// Package client
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.

package client

import (
	"encoding/json"
	"fmt"
	"github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/api"
	"net/http"
)

type DeviceIpsClient struct {
	*Client
	apiEndpoint string
}

func NewDeviceIpsClient(client *Client) *DeviceIpsClient {
	c := DeviceIpsClient{
		Client:      client,
		apiEndpoint: client.endpointBuilder.Build("device-ip"),
	}
	return &c
}

// GetCertificates - Returns list of certificates
func (c *DeviceIpsClient) GetOriginAllowList() (*api.DeviceIpsModel, error) {

	URL := fmt.Sprintf("%s/api/1.0/network/device-ip", c.apiEndpoint)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	originAllowList := api.DeviceIpsModel{}
	err = json.Unmarshal(body, &originAllowList)
	if err != nil {
		return nil, err
	}

	return &originAllowList, nil
}
