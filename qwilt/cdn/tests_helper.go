// Package cdn
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.
package cdn

import (
	b64 "encoding/base64"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"math/rand"
	"time"
)

type TerraformConfigBuilder struct {
	siteResources                  map[string]string
	certResources                  map[string]string
	siteCfgResources               map[string]string
	siteActivationResources        map[string]string
	siteActivationStagingResources map[string]string
	siteDataSources                map[string]string
	Host                           string
}

func NewTerraformConfigBuilder() *TerraformConfigBuilder {
	b := TerraformConfigBuilder{}
	b.siteResources = make(map[string]string, 0)
	b.certResources = make(map[string]string, 0)
	b.siteCfgResources = make(map[string]string, 0)
	b.siteActivationResources = make(map[string]string, 0)
	b.siteActivationStagingResources = make(map[string]string, 0)
	b.siteDataSources = make(map[string]string, 0)
	return &b
}
func (b *TerraformConfigBuilder) SitesDataSource(name, siteId string) *TerraformConfigBuilder {
	dataCfg := fmt.Sprintf(`
data "qwilt_cdn_sites" "%s" {
	filter = {
		site_id             = "%s"
		revision_id         = "all"
		publish_id          = "all"
		truncate_host_index = false
	}
}
output "site_detail" {
	value = data.qwilt_cdn_sites.%s.site[0]
}`, name, siteId, name)

	b.siteDataSources[name] = dataCfg
	return b
}
func (b *TerraformConfigBuilder) CertsDataSource(name, certId string) *TerraformConfigBuilder {
	dataCfg := fmt.Sprintf(`
data "qwilt_cdn_certificates" "%s" {
	filter = {
		cert_id             = "%s"
	}
}
output "cert_detail" {
	value = data.qwilt_cdn_certificates.%s.certificate[0]
}`, name, certId, name)

	b.siteDataSources[name] = dataCfg
	return b
}
func (b *TerraformConfigBuilder) CertResource(name, pk, cert, desc string) *TerraformConfigBuilder {
	certBase64Encode := b64.URLEncoding.EncodeToString([]byte(cert))
	pkBase64Encode := b64.URLEncoding.EncodeToString([]byte(pk))
	certCfg := fmt.Sprintf(`resource "qwilt_cdn_certificate" "%s" {
certificate       = <<EOF
%s
EOF
certificate_chain = <<EOF
%s
EOF
private_key       = <<EOF
%s
EOF
description = "%s"

lifecycle {
    ignore_changes = [
      private_key
    ]
  }

}`, name, certBase64Encode, certBase64Encode, pkBase64Encode, desc)
	b.certResources[name] = certCfg
	return b
}
func (b *TerraformConfigBuilder) SiteResource(name, siteName string) *TerraformConfigBuilder {
	siteCfg := fmt.Sprintf(`
resource "qwilt_cdn_site" "%s" {
	site_name = "%s"
}`, name, siteName)
	b.siteResources[name] = siteCfg
	return b
}
func (b *TerraformConfigBuilder) SiteConfigResource(name string, host string, changeDesc string) *TerraformConfigBuilder {
	siteConfigCfg := fmt.Sprintf(`
		resource "qwilt_cdn_site_configuration" "%s" {
			site_id = qwilt_cdn_site.%s.site_id
			host_index = <<-EOT
			{
				"hosts": [
					{
						"host": "%s",
						"host-metadata": {
							"metadata": [
								{
									"generic-metadata-type": "MI.SourceMetadataExtended",
									"generic-metadata-value": {
										"sources": [
											{
												"protocol": "https/1.1",
												"endpoints": [
													"www.example-origin-Host.com"
												]
											}
										]
									}
								}
							],
							"paths": []
						}
					}
				]
			}
			EOT
			  change_description = "%s"
			}`, name, name, host, changeDesc)
	b.siteCfgResources[name] = siteConfigCfg
	b.Host = host
	return b
}
func (b *TerraformConfigBuilder) SiteActivationResource(name string) *TerraformConfigBuilder {
	cfg := fmt.Sprintf(`
resource "qwilt_cdn_site_activation" "%s" {
		site_id = qwilt_cdn_site_configuration.%s.site_id
		revision_id = qwilt_cdn_site_configuration.%s.revision_id
	}`, name, name, name)
	b.siteActivationResources[name] = cfg
	return b
}
func (b *TerraformConfigBuilder) SiteActivationStagingResource(name string) *TerraformConfigBuilder {
	cfg := fmt.Sprintf(`
resource "qwilt_cdn_site_activation_staging" "%s" {
		site_id = qwilt_cdn_site_configuration.%s.site_id
		revision_id = qwilt_cdn_site_configuration.%s.revision_id
	}`, name, name, name)
	b.siteActivationStagingResources[name] = cfg
	return b
}
func (b *TerraformConfigBuilder) DelSiteCfgResource(name string) *TerraformConfigBuilder {
	delete(b.siteCfgResources, name)
	return b
}
func (b *TerraformConfigBuilder) DelSiteActivationResource(name string) *TerraformConfigBuilder {
	delete(b.siteActivationResources, name)
	return b
}
func (b *TerraformConfigBuilder) DelSiteResource(name string) *TerraformConfigBuilder {
	delete(b.siteResources, name)
	return b
}
func (b *TerraformConfigBuilder) DelCertResource(name string) *TerraformConfigBuilder {
	delete(b.certResources, name)
	return b
}
func (b *TerraformConfigBuilder) Build() string {
	terraformConfig := QwiltCdnFullProviderConfig
	for _, cfg := range b.certResources {
		terraformConfig += cfg + "\n"
	}
	for _, cfg := range b.siteResources {
		terraformConfig += cfg + "\n"
	}
	for _, cfg := range b.siteCfgResources {
		terraformConfig += cfg + "\n"
	}
	for _, cfg := range b.siteActivationResources {
		terraformConfig += cfg + "\n"
	}
	for _, cfg := range b.siteActivationStagingResources {
		terraformConfig += cfg + "\n"
	}
	for _, cfg := range b.siteDataSources {
		terraformConfig += cfg + "\n"
	}
	//t.Logf("config: %s", terraformConfig)
	return terraformConfig
}

func findStateResource(state *tfjson.State, resourceType string, name string) *tfjson.StateResource {
	if state.Values != nil && state.Values.RootModule != nil && state.Values.RootModule.Resources != nil {
		for _, stResource := range state.Values.RootModule.Resources {
			if stResource.Type == resourceType && stResource.Name == name {
				return stResource
			}
		}
	}
	return nil
}

//func readDataSourceForSite(state *tfjson.State, site_id string) *tfjson.StateOutput {
//	if state.Values != nil && state.Values.Outputs != nil {
//		for _, output := range state.Values.Outputs {
//			if output. {
//				return output
//			}
//		}
//	}
//	return nil
//}

func randString(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	for i := 0; i < length; i++ {
		result += string(charSet[rand.Intn(len(charSet))])
	}
	return result
}
func generateSiteName(generatedSiteName *string) string {
	randomStr := randString(16)

	// Create the final string with prefix
	siteName := fmt.Sprintf("unit-test resource site-%s", randomStr)
	*generatedSiteName = siteName
	return siteName
}
func generateHostName(generatedHostName *string) string {
	randomStr := randString(16)

	// Create the final string with prefix
	hostName := fmt.Sprintf("www.terraform-unit-test-%s.com", randomStr)
	*generatedHostName = hostName
	return hostName
}

//func VerifySite(st *terraform.State) error {
//
//	debugLog := log.New(log.Writer(), "xzjbcsjkbksjdbdskjbskj ", log.LstdFlags)
//	site_id := st.Modules[0].Resources["qwilt_cdn_site.test"].Primary.Attributes["site_id"]
//	debugLog.Println("****** state site_id: " + site_id)
//	return nil
//}
