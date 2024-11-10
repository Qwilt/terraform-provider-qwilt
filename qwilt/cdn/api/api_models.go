// Package api
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.

package api

import (
	"encoding/json"
)

const SITES_HOSTNAME = "media-sites"
const OPERATION_TYPE_UNPUBLISH = "Unpublish"

// SiteConfiguration - Model for the SiteConfiguration object
type SiteConfiguration struct {
	RevisionId                string          `json:"revisionId"`
	RevisionNum               int             `json:"revisionNum"`
	ConfigCreationTimeMilli   int             `json:"configCreationTimeMilli"`
	ConfigLastUpdateTimeMilli int             `json:"configLastUpdateTimeMilli"`
	ConfigCreatedUser         string          `json:"configCreatedUser"`
	HostIndex                 json.RawMessage `json:"hostIndex"`
	ChangeDescription         string          `json:"changeDescription"`
}

// SiteConfigChange - Model for updating the SiteConfiguration object
type SiteConfigChange struct {
	SiteId                  string          `tfsdk:"site_id"`
	OwnerOrgId              string          `json:"ownerOrgId"`
	SiteCreationTimeMilli   int             `json:"siteCreationTimeMilli"`
	SiteLastUpdateTimeMilli int             `json:"siteLastUpdateTimeMilli"`
	SiteCreatedUser         string          `json:"siteCreatedUser"`
	HostIndex               json.RawMessage `json:"hostIndex"`
	ChangeDescription       string          `json:"changeDescription"`
}

// Site - Model for the Site object
type Site struct {
	SiteId                           string         `json:"siteId"`
	OwnerOrgId                       string         `json:"ownerOrgId"`
	CreationTimeMilli                int            `json:"creationTimeMilli"`
	LastUpdateTimeMilli              int            `json:"lastUpdateTimeMilli"`
	CreatedUser                      string         `json:"createdUser"`
	LastUpdatedUser                  string         `json:"lastUpdatedUser"`
	SiteDnsCnameDelegationTarget     string         `json:"siteDnsCnameDelegationTarget"`
	SiteName                         string         `json:"siteName"`
	ApiVersion                       string         `json:"apiVersion"`
	ActiveAndLastPublishingOperation *ActiveLastPub `json:"activeAndLastPublishingOperation"`
	ServiceType                      string         `json:"serviceType"`
	RoutingMethod                    string         `json:"routingMethod"`
	ShouldProvisionToThirdPartyCdn   bool           `json:"shouldProvisionToThirdPartyCdn"`
	ServiceId                        string         `json:"serviceId"`
	IsSelfServiceBlocked             bool           `json:"isSelfServiceBlocked"`
	IsDeleted                        bool           `json:"IsDeleted"`
}

// ActiveLastPub - Contains values for activeAndLastPublishingOperation
type ActiveLastPub struct {
	Last   *PubOp `json:"last"`
	Active *PubOp `json:"active"`
}

// SiteDetailsRequest - Model for requesting Site details
type SiteDetailsRequest struct {
	SiteId string `json:"siteId"`
}

// SiteCreateRequest - Model for creating a new Site
type SiteCreateRequest struct {
	SiteName      string `json:"siteName"`
	RoutingMethod string `json:"routingMethod,omitempty"`
}

// SiteUpdateRequest - Model for updating an existing Site
type SiteUpdateRequest struct {
	SiteName string `json:"siteName"`
}

// SiteConfigVersion - Model for the config revision for a Site
type SiteConfigVersion struct {
	SiteId              string          `json:"siteId"`
	RevisionId          string          `json:"revisionId"`
	RevisionNum         int             `json:"revisionNum"`
	OwnerOrgId          string          `json:"ownerOrgId"`
	CreationTimeMilli   int             `json:"creationTimeMilli"`
	LastUpdateTimeMilli int             `json:"lastUpdateTimeMilli"`
	CreatedUser         string          `json:"createdUser"`
	HostIndex           json.RawMessage `json:"hostIndex"`
	ChangeDescription   string          `json:"changeDescription"`
}

// SiteCertificateResponse - Model for the config revision for a Site
type SiteCertificateResponse struct {
	CertificateId   string `json:"certificateId"`
	CertificateType string `json:"certificateType"`
	Target          string `json:"target"`
	State           string `json:"state"`
}

// SiteConfigAddRequest - Model for creating a new SiteConfig
type SiteConfigAddRequest struct {
	HostIndex         json.RawMessage `json:"hostIndex"`
	ChangeDescription string          `json:"changeDescription"`
}

// PubOp - Model for the publishing operation of a Site
type PubOp struct {
	PublishId                   string          `json:"publishId"`
	CreationTimeMilli           int             `json:"creationTimeMilli"`
	OwnerOrgId                  string          `json:"ownerOrgId"`
	LastUpdateTimeMilli         int             `json:"lastUpdateTimeMilli"`
	RevisionId                  string          `json:"revisionId"`
	Target                      string          `json:"target"`
	Username                    string          `json:"username"`
	PublishState                string          `json:"publishState"`
	PublishStatus               string          `json:"publishStatus"`
	PublishAcceptanceStatus     string          `json:"publishAcceptanceStatus"`
	PublishHidden               bool            `json:"publishHidden"`
	PublishMode                 string          `json:"publishMode"`
	OperationType               string          `json:"operationType"`
	StatusLine                  []string        `json:"statusLine"`
	ConfigLastModifiedTimeMilli int             `json:"configLastModifiedTimeMilli"`
	IsActive                    bool            `json:"isActive"`
	ValidatorsErrDetails        json.RawMessage `json:"validatorsErrDetails"`
}

// PubRequest - Model for requesting a new Publish operation
type PubRequest struct {
	RevisionId string `json:"revisionId"`
	Target     string `json:"target,omitempty"`
}

// UnpubRequest - Model for requesting a new Unpublish operation
type UnpubRequest struct {
	Target string `json:"target"`
}

// RepubRequest - Model for requesting a new Republish operation
type RepubRequest struct {
	Target string `json:"target"`
}

// Certificate - Model for the Certificate object
type Certificate struct {
	CertId           int64  `json:"certId"`
	Certificate      string `json:"certificate"`
	CertificateChain string `json:"certificateChain"`
	//PrivateKey       string `json:"privateKey"`
	//Email            string `json:"email"`
	Description string `json:"description"`
	PkHash      string `json:"pkHash"`
	Tenant      string `json:"tenant"`
	Domain      string `json:"domain"`
	Status      string `json:"status"`
	Type        string `json:"type"`
}

// CertificateCreateRequest - Model for creating a new Certificate
type CertificateCreateRequest struct {
	Certificate      string `json:"certificate"`
	CertificateChain string `json:"certificateChain"`
	PrivateKey       string `json:"privateKey"`
	//Email            string `json:"email"`
	Description string `json:"description"`
}

// CertificateUpdateRequest - Model for updating an existing Certificate
type CertificateUpdateRequest struct {
	Certificate      string `json:"certificate"`
	CertificateChain string `json:"certificateChain"`
	PrivateKey       string `json:"privateKey"`
	//Email            string `json:"email"`
	Description string `json:"description"`
}

type CertificateTemplate struct {
	CertificateTemplateID          int64    `json:"certificateTemplateId"`
	Country                        *string  `json:"country"`
	Tenant                         string   `json:"tenant"`
	State                          *string  `json:"state"`
	Locality                       *string  `json:"locality"`
	OrganizationName               *string  `json:"organizationName"`
	CommonName                     string   `json:"commonName"`
	SANs                           []string `json:"sans"`
	AutoManagedCertificateTemplate bool     `json:"autoManagedCertificateTemplate"`
	LastCertificateID              *int64   `json:"lastCertificateId"`
	CsrIds                         []int64  `json:"csrIds"`
}

type CertificateTemplateCreateRequest struct {
	Country                        *string  `json:"country,omitempty"`
	State                          *string  `json:"state,omitempty"`
	Locality                       *string  `json:"locality,omitempty"`
	OrganizationName               *string  `json:"organizationName,omitempty"`
	CommonName                     string   `json:"commonName"`
	SANs                           []string `json:"sans,omitempty"`
	AutoManagedCertificateTemplate bool     `json:"autoManagedCertificateTemplate"`
}

// SiteCertificateLinkRequest - Model for requesting a new Link Request
type SiteCertificateLinkRequest struct {
	CertificateId string `json:"certificateId"`
}

// IP ALLOW List model
type NetworkDeviceIpsModel struct {
	Ipv4 []string `json:"ipv4"`
	Ipv6 []string `json:"ipv6"`
}

type DeviceIpsModel struct {
	Md5              string                           `json:"md5"`
	CreateTimeMillis int                              `json:"createTimeMillis"`
	IpData           map[string]NetworkDeviceIpsModel `json:"ipData"`
}
