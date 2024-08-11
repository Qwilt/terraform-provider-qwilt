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
	"github.com/Qwilt/terraform-provider-qwilt/qwilt/qwiltcdn/api"
	"net/http"
	"strings"
)

// SiteConfigurationClient -
type SiteConfigurationClient struct {
	*Client
	apiEndpoint string
}

func NewSiteConfigurationClient(target string, client *Client) *SiteConfigurationClient {
	c := SiteConfigurationClient{
		Client:      client,
		apiEndpoint: client.endpointBuilder.Build(target),
	}
	return &c
}

// GetSiteConfigs - Returns list of site configurations
func (c *SiteConfigurationClient) GetSiteConfigs(siteId string, truncateHostIndex bool) ([]api.SiteConfigVersion, error) {
	if siteId == "" {
		return nil, fmt.Errorf("siteId is empty")
	}

	querystring := ""
	if truncateHostIndex == true {
		querystring = "?truncateHostIndex=true"
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/1/sites/%s/configurations%s", c.apiEndpoint, siteId, querystring), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	siteConfigVersions := []api.SiteConfigVersion{}
	err = json.Unmarshal(body, &siteConfigVersions)
	if err != nil {
		return nil, err
	}

	return siteConfigVersions, nil
}

// GetSiteConfig - Returns site details
func (c *SiteConfigurationClient) GetSiteConfig(siteId string, revisionId string, truncateHostIndex bool) (*api.SiteConfigVersion, error) {
	if siteId == "" || revisionId == "" {
		return nil, fmt.Errorf("Invalid input, siteId=%s revisionId=%s", siteId, revisionId)
	}

	querystring := ""
	if truncateHostIndex == true {
		querystring = "?truncateHostIndex=true"
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/1/sites/%s/configurations/%s%s", c.apiEndpoint, siteId, revisionId, querystring), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	siteConfigVersion := api.SiteConfigVersion{}
	err = json.Unmarshal(body, &siteConfigVersion)
	if err != nil {
		return nil, err
	}

	return &siteConfigVersion, nil
}

// CreateSiteConfig - Add a new site configuration
func (c *SiteConfigurationClient) CreateSiteConfig(siteId string, siteConfigVersion api.SiteConfigAddRequest) (*api.SiteConfigVersion, error) {
	if siteId == "" {
		return nil, fmt.Errorf("siteId is empty")
	}

	rb, err := json.Marshal(siteConfigVersion)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/1/sites/%s/configurations", c.apiEndpoint, siteId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	siteConfigVersionResponse := api.SiteConfigVersion{}
	err = json.Unmarshal(body, &siteConfigVersionResponse)
	if err != nil {
		return nil, err
	}
	return &siteConfigVersionResponse, nil
}
