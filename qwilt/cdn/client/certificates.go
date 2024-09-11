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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"strings"
)

type CertificatesClient struct {
	*Client
	apiEndpoint string
}

func NewCertificatesClient(client *Client) *CertificatesClient {
	c := CertificatesClient{
		Client:      client,
		apiEndpoint: client.endpointBuilder.Build("cert-manager"),
	}
	return &c
}

// GetCertificates - Returns list of certificates
func (c *CertificatesClient) GetCertificates(detailed bool) ([]api.Certificate, error) {

	querystring := ""
	if detailed == true {
		querystring = "?detailed=true"
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/certificates%s", c.apiEndpoint, querystring), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	certs := []api.Certificate{}
	err = json.Unmarshal(body, &certs)
	if err != nil {
		return nil, err
	}

	return certs, nil
}

// GetCertificate - Returns certificate details
func (c *CertificatesClient) GetCertificate(certId types.Int64, detailed bool) (*api.Certificate, error) {
	if certId.IsNull() {
		return nil, fmt.Errorf("certId is empty")
	}

	querystring := ""
	if detailed == true {
		querystring = "?detailed=true"
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/certificates/%s%s", c.apiEndpoint, certId, querystring), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	certDetail := api.Certificate{}
	err = json.Unmarshal(body, &certDetail)
	if err != nil {
		return nil, err
	}

	return &certDetail, nil
}

// CreateCertificate - Create new certificate
func (c *CertificatesClient) CreateCertificate(cert api.CertificateCreateRequest) (*api.Certificate, error) {
	rb, err := json.Marshal(cert)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v2/certificates", c.apiEndpoint), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCert := api.Certificate{}
	err = json.Unmarshal(body, &newCert)
	if err != nil {
		return nil, err
	}

	return &newCert, nil
}

// UpdateCertificate - Update cert details
func (c *CertificatesClient) UpdateCertificate(certId int64, site api.CertificateUpdateRequest) (*api.Certificate, error) {

	rb, err := json.Marshal(site)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v2/certificates/%d", c.apiEndpoint, certId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedCert := api.Certificate{}
	err = json.Unmarshal(body, &updatedCert)
	if err != nil {
		return nil, err
	}

	return &updatedCert, nil
}

// DeleteCertificate - Deletes a certificate
func (c *CertificatesClient) DeleteCertificate(certId types.Int64) error {
	if certId.IsNull() {
		return fmt.Errorf("certId is empty")
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v2/certificates/%s", c.apiEndpoint, certId), nil)
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
