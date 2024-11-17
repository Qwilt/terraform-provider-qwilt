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
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tfjson "github.com/hashicorp/terraform-json"
)

type TerraformConfigBuilder struct {
	siteResources                  map[string]string
	certificateTemplateResources   map[string]string
	certResources                  map[string]string
	siteCfgResources               map[string]string
	siteActivationResources        map[string]string
	siteActivationStagingResources map[string]string
	siteDataSources                map[string]string
	Host                           string
}

func SetDevOverrides() {
	//set this after running script scripts/generate_dev_overrides.sh
	os.Setenv("TF_CLI_CONFIG_FILE", "/Users/yuvalw/Qwilt/projects/qc-microservices/fork/terraform-provider-qwilt/examples/playground/qwilt-config-certificate-template/bin/developer_overrides.tfrc")

}
func NewTerraformConfigBuilder() *TerraformConfigBuilder {
	b := TerraformConfigBuilder{}
	b.siteResources = make(map[string]string, 0)
	b.certResources = make(map[string]string, 0)
	b.certificateTemplateResources = make(map[string]string, 0)
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
func (b *TerraformConfigBuilder) CertificateTemplateDataSource(name, id string) *TerraformConfigBuilder {
	dataCfg := fmt.Sprintf(`
data "qwilt_cdn_certificate_templates" "%s" {
	filter = {
		id             = "%s"
	}
}
output "certificate_template" {
	value = data.qwilt_cdn_certificate_templates.%s.certificate_template[0]
}`, name, id, name)

	b.certificateTemplateResources[name] = dataCfg
	return b
}
func (b *TerraformConfigBuilder) CertificateTemplateResource(name, commonName, orgName string, sans []string, autoManaged bool) *TerraformConfigBuilder {
	certificateTemplateConfig := fmt.Sprintf(`resource "qwilt_cdn_certificate_template" "%s" {
	common_name = "%s"
	sans = %s
	organization_name = "%s"
	auto_managed_certificate_template = %t

}`, name, commonName, "[\""+strings.Join(sans, "\",\"")+"\"]", orgName, autoManaged)
	b.certificateTemplateResources[name] = certificateTemplateConfig
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
													"www.example-origin-host.com"
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
func (b *TerraformConfigBuilder) SiteConfigResourceWithTabs(name string, host string, changeDesc string) *TerraformConfigBuilder {
	siteConfigCfg := fmt.Sprintf(`
		resource "qwilt_cdn_site_configuration" "%s" {
			site_id = qwilt_cdn_site.%s.site_id
			host_index = <<-EOT
			{
				"hosts":		 [
					{
						"host": 	"%s",
						"host-metadata": 	{
							"metadata": 	[
								{
									"generic-metadata-type": "MI.SourceMetadataExtended",
									"generic-metadata-value": {
										"sources": [
											{
												"protocol": "https/1.1",
												"endpoints": [
													"www.example-origin-host.com"
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
func (b *TerraformConfigBuilder) SiteActivationResourceWithCert(name string, cert_id, csr_id *int) *TerraformConfigBuilder {
	var cert_id_str string = "null"
	var csr_id_str string = "null"
	if cert_id != nil {
		cert_id_str = strconv.Itoa(*cert_id)
	}
	if csr_id != nil {
		csr_id_str = strconv.Itoa(*csr_id)
	}
	cfg := fmt.Sprintf(`
resource "qwilt_cdn_site_activation" "%s" {
		site_id = qwilt_cdn_site_configuration.%s.site_id
		revision_id = qwilt_cdn_site_configuration.%s.revision_id
		certificate_id = %s
		csr_id = %s
	}`, name, name, name, cert_id_str, csr_id_str)
	b.siteActivationResources[name] = cfg
	return b
}
func (b *TerraformConfigBuilder) SiteActivationResourceWithCertRef(name string, cert_ref_name string) *TerraformConfigBuilder {
	cfg := fmt.Sprintf(`
resource "qwilt_cdn_site_activation" "%s" {
		site_id = qwilt_cdn_site_configuration.%s.site_id
		revision_id = qwilt_cdn_site_configuration.%s.revision_id
		certificate_id = qwilt_cdn_certificate.%s.cert_id
	}`, name, name, name, cert_ref_name)
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
func (b *TerraformConfigBuilder) DelCertificateTemplateResource(name string) *TerraformConfigBuilder {
	delete(b.certificateTemplateResources, name)
	return b
}
func (b *TerraformConfigBuilder) DelCertResource(name string) *TerraformConfigBuilder {
	delete(b.certResources, name)
	return b
}
func (b *TerraformConfigBuilder) OriginAllowListDataSource(name string) *TerraformConfigBuilder {
	dataCfg := fmt.Sprintf(`
data "qwilt_cdn_origin_allow_list" "%s" {
}
output "origin_allow_list_detail" {
	value = data.qwilt_cdn_origin_allow_list.%s
}`, name, name)

	b.siteDataSources[name] = dataCfg
	return b
}
func (b *TerraformConfigBuilder) Build() string {
	terraformConfig := QwiltCdnFullProviderConfig
	for _, cfg := range b.certificateTemplateResources {
		terraformConfig += cfg + "\n"
	}
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

	charSet := "abcdefghijklmnopqrstuvwxyz0123456789"
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
	hostName := fmt.Sprintf("www.%s-terraform-unit-test.com", randomStr)
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
