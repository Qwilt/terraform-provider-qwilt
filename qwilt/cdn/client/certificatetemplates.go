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
	"net/http"
	"strings"

	"github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/api"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ApiRoot = "/api/v2/certificate-templates"

type CertificateTemplatesClient struct {
	*Client
	apiEndpoint string
}

func NewCertificateTemplateClient(client *Client) *CertificateTemplatesClient {
	c := CertificateTemplatesClient{
		Client:      client,
		apiEndpoint: client.endpointBuilder.Build("cert-manager"),
	}
	return &c
}

// GetCertificateTemplates - Returns list of Certificate Templates
func (c *CertificateTemplatesClient) GetCertificateTemplates() ([]api.CertificateTemplate, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.apiEndpoint, ApiRoot), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var certs []api.CertificateTemplate
	err = json.Unmarshal(body, &certs)
	if err != nil {
		return nil, err
	}

	return certs, nil
}

// GetCertificateTemplate - Returns Certificate Template details
func (c *CertificateTemplatesClient) GetCertificateTemplate(id types.Int64) (*api.CertificateTemplate, error) {
	if id.IsNull() {
		return nil, fmt.Errorf("certificate template id is empty")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.apiEndpoint, ApiRoot, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	certDetail := api.CertificateTemplate{}
	err = json.Unmarshal(body, &certDetail)
	if err != nil {
		return nil, err
	}

	return &certDetail, nil
}

// CreateCertificateTemplate - Create new Certificate Template
func (c *CertificateTemplatesClient) CreateCertificateTemplate(cert api.CertificateTemplateCreateRequest) (*api.CertificateTemplate, error) {
	rb, err := json.Marshal(cert)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.apiEndpoint, ApiRoot), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCert := api.CertificateTemplate{}
	err = json.Unmarshal(body, &newCert)
	if err != nil {
		return nil, err
	}

	return &newCert, nil
}

// DeleteCertificateTemplate - Deletes a certificate Template
func (c *CertificateTemplatesClient) DeleteCertificateTemplate(id types.Int64) error {
	if id.IsNull() {
		return fmt.Errorf("certificate template id is empty")
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.apiEndpoint, ApiRoot, id), nil)
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
