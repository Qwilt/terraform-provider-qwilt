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
	"net/url"
	"strings"
	"time"
)

const TARGET_GA = "ga"
const TARGET_STAGING = "staging"
const ACCEPTANCE_STATUS_PENDING = "Pending"
const ACCEPTANCE_STATUS_INVALID = "Invalid"
const ACCEPTANCE_STATUS_DISMISSED = "Dismissed"
const ACCEPTANCE_TIMEOUT = 180 * time.Second

type PublishOpsClient struct {
	*Client
	apiEndpoint string
}

func NewPublishOpsClient(target string, client *Client) *PublishOpsClient {
	c := PublishOpsClient{
		Client:      client,
		apiEndpoint: client.endpointBuilder.Build(target),
	}
	return &c
}

// FindPubOp - Returns latest publishing operation for site
func (c *PublishOpsClient) FindLatestPubOp(siteId string, revisionId string) (*api.PubOp, error) {
	if siteId == "" || revisionId == "" {
		return nil, fmt.Errorf("Invalid input, siteId=%s revisionId=%s", siteId, revisionId)
	}

	pubOp := api.PubOp{}

	// Get list of publish ops
	pubOps, err := c.GetPubOps(siteId, false, "")
	if err != nil || len(pubOps) == 0 {
		// No error, but no jobs found
		return &pubOp, nil
	}

	// Traverse publishops to find latest operation.  Criteria:
	// 1. Matching revisionId
	// 2. isActive or InProgress
	// 3. Prefer InProgress to isActive
	for _, pubOpCandidate := range pubOps {
		if pubOpCandidate.RevisionId != revisionId {
			continue
		}
		if pubOpCandidate.PublishStatus == "InProgress" {
			pubOp = pubOpCandidate
			break
		}
		if pubOpCandidate.IsActive {
			pubOp = pubOpCandidate
			// Not breaking since there might be a better candidate
		}
	}

	return &pubOp, nil
}

// GetSitePubStatus - Get a site's publishing status
func (c *PublishOpsClient) GetSitePubStatus(siteId string) (string, string, error) {
	pubOps, err := c.GetPubOps(siteId, false, "")
	if err != nil {
		return "", "", err
	}

	currentStatus := "Unpublished"
	var transitionStatus string

	// No records = Unpublished
	if len(pubOps) != 0 {
		for _, pubOp := range pubOps {
			// Return Published or Unpublished if isActive
			if pubOp.IsActive == true {
				currentStatus = pubOp.OperationType + "ed"
			}
			// Return Publishing or Unpublishing if isActive
			if pubOp.PublishStatus == "InProgress" {
				transitionStatus = pubOp.OperationType + "ing"
			}
		}
	}

	// No transition found
	if transitionStatus == "" {
		transitionStatus = currentStatus
	}

	return currentStatus, transitionStatus, err
}

// GetPubOps - Returns list of publishing operations
func (c *PublishOpsClient) GetPubOps(siteId string, isActive bool, publishState string) ([]api.PubOp, error) {
	if siteId == "" {
		return nil, fmt.Errorf("siteId is empty")
	}

	// Build URL with query string
	base, err := url.Parse(fmt.Sprintf("%s/api/v2/sites/%s/publishing-operations", c.apiEndpoint, siteId))
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	if isActive == true {
		params.Add("isActive", "isActive")
	}
	if publishState != "" {
		params.Add("publishState", publishState)
	}
	base.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", base.String(), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	pubOps := []api.PubOp{}
	err = json.Unmarshal(body, &pubOps)
	if err != nil {
		return nil, err
	}

	return pubOps, nil
}

// GetPubOp - Returns details on a publishing operation
func (c *PublishOpsClient) GetPubOp(siteId string, publishId string) (*api.PubOp, error) {
	if siteId == "" || publishId == "" {
		return nil, fmt.Errorf("Invalid input, siteId=%s publishId=%s", siteId, publishId)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/sites/%s/publishing-operations/%s", c.apiEndpoint, siteId, publishId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	pubOp := api.PubOp{}
	err = json.Unmarshal(body, &pubOp)
	if err != nil {
		return nil, err
	}

	return &pubOp, nil
}

// GetAndWaitForPubOpAcceptance - Returns details about a publishing operation after waiting for it to complete validation step.
func (c *PublishOpsClient) GetAndWaitForPubOpAcceptance(siteId string, publishId string, timeout time.Duration) (*api.PubOp, error) {
	if siteId == "" || publishId == "" {
		return nil, fmt.Errorf("Invalid input, siteId=%s publishId=%s", siteId, publishId)
	}

	start := time.Now()
	var pubOpGetResp *api.PubOp
	var err error

	for time.Since(start) < timeout {
		pubOpGetResp, err = c.GetPubOp(siteId, publishId)

		if err != nil {
			return nil, err
		}

		if pubOpGetResp.PublishAcceptanceStatus != ACCEPTANCE_STATUS_PENDING {
			return pubOpGetResp, nil
		}
		time.Sleep(3 * time.Second) // Wait for few seconds before checking again
	}
	return pubOpGetResp, fmt.Errorf("Publish Operation TimedOut waiting for acceptance status for siteId=%s publishId=%s", siteId, publishId)
}

// Publish - Publish a site
func (c *PublishOpsClient) Publish(siteId string, revisionId string, target string) (*api.PubOp, error) {
	if siteId == "" || revisionId == "" || target == "" {
		return nil, fmt.Errorf("Invalid input, siteId=%s revisionId=%s target=%s", siteId, revisionId, target)
	}

	pubReq := api.PubRequest{}
	pubReq.RevisionId = revisionId
	pubReq.Target = target

	rb, err := json.Marshal(pubReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v2/sites/%s/publishing-operations", c.apiEndpoint, siteId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	pubResp := api.PubOp{}
	err = json.Unmarshal(body, &pubResp)
	if err != nil {
		return nil, err
	}

	return &pubResp, nil
}

// Unpublish - Unpublish a site
func (c *PublishOpsClient) Unpublish(siteId string, target string) (*api.PubOp, error) {
	if siteId == "" || target == "" {
		return nil, fmt.Errorf("Invalid input, siteId=%s target=%s", siteId, target)
	}

	unpubReq := api.UnpubRequest{}
	unpubReq.Target = target

	rb, err := json.Marshal(unpubReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v2/sites/%s/publishing-operations/actions/un-publish", c.apiEndpoint, siteId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	unpubResp := api.PubOp{}
	err = json.Unmarshal(body, &unpubResp)
	if err != nil {
		return nil, err
	}

	return &unpubResp, nil
}

// Republish - Republish the active configuration version
func (c *PublishOpsClient) Republish(siteId string, target string) (*api.PubOp, error) {
	if siteId == "" || target == "" {
		return nil, fmt.Errorf("Invalid input, siteId=%s target=%s", siteId, target)
	}

	repubReq := api.RepubRequest{}
	repubReq.Target = target

	rb, err := json.Marshal(repubReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v2/sites/%s/publishing-operations/actions/republish", c.apiEndpoint, siteId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	repubResp := api.PubOp{}
	err = json.Unmarshal(body, &repubResp)
	if err != nil {
		return nil, err
	}

	return &repubResp, nil
}

// Cancel - Cancel an ongoing publish operation
func (c *PublishOpsClient) Cancel(siteId string, publishId string) error {
	if siteId == "" || publishId == "" {
		return fmt.Errorf("Invalid input, siteId=%s publishIdf=%s", siteId, publishId)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v2/sites/%s/publishing-operations/%s/actions/cancel", c.apiEndpoint, siteId, publishId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
