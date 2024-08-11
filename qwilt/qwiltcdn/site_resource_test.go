// Package qwiltcdn
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.
package qwiltcdn

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestSiteResource(t *testing.T) {

	t.Logf("Starting TestSiteResource test DEBUG: ")

	//os.Setenv("TF_CLI_CONFIG_FILE", "/Users/efrats/.terraformrc")

	tfBinaryPath := "terraform"
	var curSiteName, curSiteName2 string

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

	terraformBuilder := NewTerraformConfigBuilder().SiteResource("test", generateSiteName(&curSiteName))
	terraformConfig := terraformBuilder.Build()

	//t.Logf("config: %s", terraformConfig)
	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err := tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(state.Values.RootModule.Resources))

	siteState := findStateResource(state, "qwiltcdn_site", "test")
	assert.NotNil(t, siteState)

	t.Logf("site: %s", siteState.AttributeValues)
	assert.Equal(t, curSiteName, siteState.AttributeValues["site_name"])
	assert.Equal(t, "devorg", siteState.AttributeValues["owner_org_id"])
	assert.Equal(t, "DNS", siteState.AttributeValues["routing_method"])
	return
	//assert.Equal(t, "edge.ds-c7409-devorg.global.kan11-c.tc-rnd.cqloud.com", siteState.AttributeValues["site_dns_cname_delegation_target"])

	//check that plan gives no diff - this actually checks the refresh and that all attributes in the state are the same as in the configuration
	plan, err := tf.Plan(context.Background())
	assert.Equal(t, nil, err)
	assert.False(t, plan) //no diff

	//get the siteid to test it later with import
	siteId := siteState.AttributeValues["site_id"]

	//update site name
	terraformBuilder.SiteResource("test", generateSiteName(&curSiteName2))
	terraformConfig = terraformBuilder.Build()

	//t.Logf("config: %s", terraformConfig)
	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(state.Values.RootModule.Resources))

	siteState = findStateResource(state, "qwiltcdn_site", "test")
	assert.NotNil(t, siteState)

	assert.Equal(t, siteId, siteState.AttributeValues["site_id"])
	assert.Equal(t, curSiteName2, siteState.AttributeValues["site_name"])
	assert.Equal(t, "devorg", siteState.AttributeValues["owner_org_id"])
	assert.Equal(t, "DNS", siteState.AttributeValues["routing_method"])

	err = tf.StateRm(context.Background(), "qwiltcdn_site.test")
	assert.Equal(t, nil, err)

	err = tf.Import(context.Background(), "qwiltcdn_site.test", fmt.Sprintf("%s", siteId))
	assert.Equal(t, nil, err)

	siteState = findStateResource(state, "qwiltcdn_site", "test")
	assert.NotNil(t, siteState)

	assert.Equal(t, curSiteName2, siteState.AttributeValues["site_name"])
	assert.Equal(t, "devorg", siteState.AttributeValues["owner_org_id"])
	assert.Equal(t, "DNS", siteState.AttributeValues["routing_method"])

	//remove the configuration and check that it is destroyed
	terraformBuilder.DelSiteResource("test")
	terraformConfig = terraformBuilder.Build()

	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Nil(t, state.Values)
}
