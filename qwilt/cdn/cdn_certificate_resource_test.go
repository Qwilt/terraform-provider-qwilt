// Package qwilt_cdn
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.
package cdn

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strings"
	"testing"
)

func TestCertificateResource(t *testing.T) {

	t.Logf("Starting TestSiteResource test")

	//set this after running script generate_dev_overrides.sh
	SetDevOverrides()

	tfBinaryPath := "terraform"

	// Create a temporary directory to hold the Terraform configuration
	tempDir, err := os.MkdirTemp("", "tf-exec-example")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %s", err)
	}

	defer os.RemoveAll(tempDir) // Clean up the temporary directory after the test

	// Write the Terraform configuration to a file in the temporary directory
	tfFilePath := tempDir + "/main.tf"

	// Initialize a new Terraform instance
	tf, err := tfexec.NewTerraform(tempDir, tfBinaryPath)
	assert.Equal(t, nil, err)

	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)

	var domain string
	generateHostName(&domain)
	certGen := NewSelfSignedCertGenerator()
	certGen.generate(domain)
	//t.Logf("pk: %s", certGen.PK)
	//t.Logf("cert: %s", certGen.Crt)

	terraformBuilder := NewTerraformConfigBuilder().CertResource("test", certGen.PK, certGen.Crt, "ccc")
	terraformConfig := terraformBuilder.Build()

	//t.Logf("config: %s", terraformConfig)
	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err := tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(state.Values.RootModule.Resources))

	certState := findStateResource(state, "qwilt_cdn_certificate", "test")
	assert.NotNil(t, certState)

	certId := certState.AttributeValues["cert_id"]
	t.Logf("cert_id: %s", certId)
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(certGen.Crt)), strings.TrimSpace(certState.AttributeValues["certificate"].(string)))
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(certGen.Crt)), strings.TrimSpace(certState.AttributeValues["certificate_chain"].(string)))
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(certGen.PK)), strings.TrimSpace(certState.AttributeValues["private_key"].(string)))
	assert.Equal(t, "ccc", certState.AttributeValues["description"])
	assert.Equal(t, domain, certState.AttributeValues["domain"])
	assert.Equal(t, "devorg", certState.AttributeValues["tenant"])

	//check that plan gives no diff - this actually checks the refresh and that all attributes in the state are the same as in the configuration
	plan, err := tf.Plan(context.Background())
	assert.Equal(t, nil, err)
	assert.False(t, plan) //no diff

	//update just the description
	terraformBuilder.CertResource("test", certGen.PK, certGen.Crt, "description-after-change")
	terraformConfig = terraformBuilder.Build()

	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(state.Values.RootModule.Resources))

	certState = findStateResource(state, "qwilt_cdn_certificate", "test")
	assert.NotNil(t, certState)

	assert.Equal(t, certId, certState.AttributeValues["cert_id"]) //make sure therewas no replace, just update
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(certGen.Crt)), strings.TrimSpace(certState.AttributeValues["certificate"].(string)))
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(certGen.Crt)), strings.TrimSpace(certState.AttributeValues["certificate_chain"].(string)))
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(certGen.PK)), strings.TrimSpace(certState.AttributeValues["private_key"].(string)))
	assert.Equal(t, "description-after-change", certState.AttributeValues["description"])
	assert.Equal(t, domain, certState.AttributeValues["domain"])
	assert.Equal(t, "devorg", certState.AttributeValues["tenant"])

	//check that plan gives no diff - this actually checks the refresh and that all attributes in the state are the same as in the configuration
	plan, err = tf.Plan(context.Background())
	assert.Equal(t, nil, err)
	assert.False(t, plan) //no diff

	//prepare for import
	err = tf.StateRm(context.Background(), "qwilt_cdn_certificate.test")
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Nil(t, state.Values)

	err = tf.Import(context.Background(), "qwilt_cdn_certificate.test", fmt.Sprintf("%s", certId))
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(state.Values.RootModule.Resources))

	certState = findStateResource(state, "qwilt_cdn_certificate", "test")
	assert.NotNil(t, certState)

	assert.Equal(t, certId, certState.AttributeValues["cert_id"])
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(certGen.Crt)), strings.TrimSpace(certState.AttributeValues["certificate"].(string)))
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(certGen.Crt)), strings.TrimSpace(certState.AttributeValues["certificate_chain"].(string)))
	//assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(domain_key)), strings.TrimSpace(certState.AttributeValues["private_key"].(string)))
	assert.Equal(t, "description-after-change", certState.AttributeValues["description"])
	assert.Equal(t, domain, certState.AttributeValues["domain"])
	assert.Equal(t, "devorg", certState.AttributeValues["tenant"])

	//check that plan gives no diff - this actually checks the refresh and that all attributes in the state are the same as in the configuration
	plan, err = tf.Plan(context.Background())
	assert.Equal(t, nil, err)
	assert.False(t, plan) //no diff

	terraformBuilder.DelCertResource("test")
	terraformConfig = terraformBuilder.Build()

	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Nil(t, state.Values)
}
