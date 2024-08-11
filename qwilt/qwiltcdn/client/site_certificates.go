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

// SiteCertificatesClient -
type SiteCertificatesClient struct {
	*Client
	apiEndpoint string
}

func NewSiteCertificatesClient(target string, client *Client) *SiteCertificatesClient {
	c := SiteCertificatesClient{
		Client:      client,
		apiEndpoint: client.endpointBuilder.Build(target),
	}
	return &c
}

// GetSiteCertficates - Returns list of site certificates
func (c *SiteCertificatesClient) GetSiteCertificates(siteId string, revisionId string) ([]api.SiteCertificateResponse, error) {
	if siteId == "" {
		return nil, fmt.Errorf("siteId is empty")
	}

	querystring := ""
	if revisionId != "" {
		querystring = fmt.Sprintf("?siteRevisionId=%s", revisionId)
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/sites/%s/certificates%s", c.apiEndpoint, siteId, querystring), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	siteCertificates := []api.SiteCertificateResponse{}
	err = json.Unmarshal(body, &siteCertificates)
	if err != nil {
		return nil, err
	}

	return siteCertificates, nil
}

// LinkSiteCertificate - links a site to a certificate
func (c *SiteCertificatesClient) LinkSiteCertificate(siteId string, certId string) (*api.SiteCertificateResponse, error) {
	if siteId == "" || certId == "" {
		return nil, fmt.Errorf("Invalid input, siteId=%s certId=%s", siteId, certId)
	}

	certsResp, err := c.GetSiteCertificates(siteId, "")
	if err != nil {
		return nil, err
	}

	if len(certsResp) > 0 {
		//for now QC supports only one certificate. unlink the previously linked one
		err := c.UnLinkSiteCertificate(siteId, certsResp[0].CertificateId)
		if err != nil {
			return nil, err
		}
	}

	linkReq := api.SiteCertificateLinkRequest{}
	linkReq.CertificateId = certId

	rb, err := json.Marshal(linkReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v2/sites/%s/certificates", c.apiEndpoint, siteId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	linkResponse := []api.SiteCertificateResponse{}
	err = json.Unmarshal(body, &linkResponse)
	if err != nil {
		return nil, err
	}

	return &linkResponse[0], nil
}

// LinkSiteCertificate - links asite to a certificate
func (c *SiteCertificatesClient) UnLinkSiteCertificate(siteId string, certId string) error {
	if siteId == "" || certId == "" {
		return fmt.Errorf("Invalid input, siteId=%s certId=%s", siteId, certId)
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v2/sites/%s/certificates/%s", c.apiEndpoint, siteId, certId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
