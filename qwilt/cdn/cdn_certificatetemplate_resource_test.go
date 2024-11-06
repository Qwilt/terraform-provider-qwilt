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
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/stretchr/testify/assert"
)

func TestCertificateTemplateResource(t *testing.T) {

	t.Logf("Starting TestCertificateTemplateResource test")

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
	orgName := "qwilt"
	sans := []string{"test1.com", "test2.com"}

	terraformBuilder := NewTerraformConfigBuilder().CertificateTemplateResource("test", domain, orgName, sans, true)
	terraformConfig := terraformBuilder.Build()

	//t.Logf("config: %s", terraformConfig)
	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err := tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(state.Values.RootModule.Resources))

	certificateTemplateState := findStateResource(state, "qwilt_cdn_certificate_template", "test")
	assert.NotNil(t, certificateTemplateState)

	id := certificateTemplateState.AttributeValues["certificate_template_id"]
	t.Logf("certificate_template_id: %s", id)
	assert.Equal(t, domain, certificateTemplateState.AttributeValues["common_name"])
	assert.Equal(t, orgName, certificateTemplateState.AttributeValues["organization_name"])
	assert.Equal(t, sans, certificateTemplateState.AttributeValues["sans"])

	//check that plan gives no diff - this actually checks the refresh and that all attributes in the state are the same as in the configuration
	plan, err := tf.Plan(context.Background())
	assert.Equal(t, nil, err)
	assert.False(t, plan) //no diff

	//update just the domain
	domain = "newdomain.com"
	terraformBuilder.CertificateTemplateResource("test", domain, orgName, sans, true)
	terraformConfig = terraformBuilder.Build()

	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(state.Values.RootModule.Resources))

	certificateTemplateState = findStateResource(state, "qwilt_cdn_certificate_template", "test")
	assert.NotNil(t, certificateTemplateState)

	assert.NotEqual(t, id, certificateTemplateState.AttributeValues["certificate_template_id"]) //make sure it changed
	assert.Equal(t, domain, certificateTemplateState.AttributeValues["common_name"])
	assert.Equal(t, orgName, certificateTemplateState.AttributeValues["organization_name"])
	assert.Equal(t, sans, certificateTemplateState.AttributeValues["sans"])

	//check that plan gives no diff - this actually checks the refresh and that all attributes in the state are the same as in the configuration
	plan, err = tf.Plan(context.Background())
	assert.Equal(t, nil, err)
	assert.False(t, plan) //no diff

	//prepare for import
	err = tf.StateRm(context.Background(), "qwilt_cdn_certificate_template.test")
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Nil(t, state.Values)

	err = tf.Import(context.Background(), "qwilt_cdn_certificate_template.test", fmt.Sprintf("%s", id))
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(state.Values.RootModule.Resources))

	certificateTemplateState = findStateResource(state, "qwilt_cdn_certificate_template", "test")
	assert.NotNil(t, certificateTemplateState)

	assert.NotEqual(t, id, certificateTemplateState.AttributeValues["certificate_template_id"])
	assert.Equal(t, domain, certificateTemplateState.AttributeValues["common_name"])
	assert.Equal(t, orgName, certificateTemplateState.AttributeValues["organization_name"])
	assert.Equal(t, sans, certificateTemplateState.AttributeValues["sans"])

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
