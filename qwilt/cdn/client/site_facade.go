// Package client
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.

package client

// Assuming "ga" as target environment
var targetEnvironment string = "ga"

// SiteClientFacade -
type SiteClientFacade struct {
	*PublishOpsClient
	*SiteClient
	*SiteConfigurationClient
	*SiteCertificatesClient
	*CertificateTemplateClient
	ApiEndpoint string
}

// Decorator on top of Client type
func NewSiteFacadeClient(target string, client *Client) *SiteClientFacade {
	c := SiteClientFacade{
		PublishOpsClient:          NewPublishOpsClient(target, client),
		SiteClient:                NewSiteClient(target, client),
		SiteConfigurationClient:   NewSiteConfigurationClient(target, client),
		SiteCertificatesClient:    NewSiteCertificatesClient(target, client),
		CertificateTemplateClient: NewCertificateTemplateClient(client),
	}
	return &c
}
