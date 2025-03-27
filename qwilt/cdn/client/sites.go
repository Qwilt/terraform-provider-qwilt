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
	"errors"
	"fmt"
	"github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/api"
	"log"
	"net/http"
	"strings"
)

type SiteClient struct {
	*Client
	apiEndpoint string
}

func NewSiteClient(target string, client *Client) *SiteClient {
	c := SiteClient{
		Client:      client,
		apiEndpoint: client.endpointBuilder.Build(target),
	}
	return &c
}

// GetSites - Returns list of sites
func (c *SiteClient) GetSites(includeActiveLastPub bool, includeDeletedSites bool) ([]api.Site, error) {
	var param string
	if includeActiveLastPub {
		param = "?includePublishDetails=true"
	} else {
		param = ""
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/sites%s", c.apiEndpoint, param), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sites := []api.Site{}
	err = json.Unmarshal(body, &sites)
	if err != nil {
		return nil, err
	}

	if includeDeletedSites {
		return sites, nil
	}

	//filter deleted sites
	nonDeletedSites := []api.Site{}
	for _, site := range sites {
		if site.IsDeleted {
			continue
		}
		nonDeletedSites = append(nonDeletedSites, site)
	}
	return nonDeletedSites, nil
}

// GetSite - Returns site details
func (c *SiteClient) GetSite(siteId string, target string, includeActiveLastPub bool, includeDeletedSites bool) (*api.Site, error) {
	if siteId == "" {
		return nil, fmt.Errorf("siteId is empty")
	}

	var queryParams string
	if includeActiveLastPub && target != "" {
		queryParams = fmt.Sprintf("?includePublishDetails=true&publishTarget=%s", target)
	} else if includeActiveLastPub && target == "" {
		queryParams = fmt.Sprintf("?includePublishDetails=true&publishTarget=ga")
	} else if !includeActiveLastPub && target != "" {
		queryParams = fmt.Sprintf("?publishTarget=%s", target)
	} else {
		queryParams = ""
	}

	url := fmt.Sprintf("%s/api/v2/sites/%s%s", c.apiEndpoint, siteId, queryParams)
	log.Printf("***** %s *****", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	siteDetail := api.Site{}
	err = json.Unmarshal(body, &siteDetail)
	if err != nil {
		return nil, err
	}

	if !includeDeletedSites && siteDetail.IsDeleted {
		return nil, errors.New("Site was found but marked for deletion")
	}
	return &siteDetail, nil
}

// CreateSite - Create new site
func (c *SiteClient) CreateSite(site api.SiteCreateRequest) (*api.Site, error) {
	rb, err := json.Marshal(site)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v2/sites", c.apiEndpoint), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newSite := api.Site{}
	err = json.Unmarshal(body, &newSite)
	if err != nil {
		return nil, err
	}

	return &newSite, nil
}

// UpdateSite - Update site details
func (c *SiteClient) UpdateSite(siteId string, site api.SiteUpdateRequest) (*api.Site, error) {
	if siteId == "" {
		return nil, fmt.Errorf("siteId is empty")
	}

	rb, err := json.Marshal(site)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v2/sites/%s", c.apiEndpoint, siteId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedSite := api.Site{}
	err = json.Unmarshal(body, &updatedSite)
	if err != nil {
		return nil, err
	}

	return &updatedSite, nil
}

// DeleteSite - Deletes a site
func (c *SiteClient) DeleteSite(siteId string) error {
	if siteId == "" {
		return fmt.Errorf("siteId is empty")
	}

	url := fmt.Sprintf("%s/api/v2/sites/%s?permanent=true", c.apiEndpoint, siteId)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	_ = body
	if err != nil {
		return err
	}

	return nil
}
