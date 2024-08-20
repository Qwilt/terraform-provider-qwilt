// Package qwiltcdn
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.
package cdn

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestSiteConfigResource(t *testing.T) {

	t.Logf("Starting TestSiteConfigResource test DEBUG: ")

	//os.Setenv("TF_CLI_CONFIG_FILE", "/Users/efrats/.terraformrc")
	os.Setenv("TF_LOG", "INFO")
	os.Setenv("TF_LOG_PROVIDER", "INFO")

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

	//tf.SetStdout(os.Stdout)
	//tf.SetStderr(os.Stderr)

	var curSiteName, curHostName string
	generateHostName(&curHostName)
	var changeDesc = fmt.Sprintf("Terraform plugin unit testing description for site %s", curSiteName)

	terraformBuilder := NewTerraformConfigBuilder()
	terraformBuilder.SiteResource("test", generateSiteName(&curSiteName))
	terraformBuilder.SiteConfigResource("test", curHostName, changeDesc)
	terraformConfig := terraformBuilder.Build()

	//t.Logf("config: %s", terraformConfig)
	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err := tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(state.Values.RootModule.Resources))

	siteCfgState := findStateResource(state, "qwilt_cdn_site_configuration", "test")
	assert.NotNil(t, siteCfgState)

	//get the site_id and revision_id to test it later with import
	siteId := siteCfgState.AttributeValues["site_id"]
	revisionId1 := siteCfgState.AttributeValues["revision_id"]
	t.Logf("siteId: %s", siteId)
	t.Logf("revisionId1: %s", revisionId1)

	revisionNum1 := siteCfgState.AttributeValues["revision_num"].(json.Number)
	assert.Equal(t, "1", revisionNum1.String())
	assert.Equal(t, changeDesc, siteCfgState.AttributeValues["change_description"])

	//check that plan gives no diff - this actually checks the refresh and that all attributes in the state are the same as in the configuration
	plan, err := tf.Plan(context.Background())
	assert.Equal(t, nil, err)
	assert.False(t, plan) //no diff

	//update host and desc
	terraformBuilder.SiteConfigResource("test", "www.unitests-2.com", "yyy")
	terraformConfig = terraformBuilder.Build()

	//t.Logf("config: %s", terraformConfig)
	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(state.Values.RootModule.Resources))

	siteCfgState = findStateResource(state, "qwilt_cdn_site_configuration", "test")
	assert.NotNil(t, siteCfgState)

	revisionNum2 := siteCfgState.AttributeValues["revision_num"].(json.Number)
	revisionId2 := siteCfgState.AttributeValues["revision_id"]

	t.Logf("revisionId2: %s", revisionId2)
	assert.Equal(t, "2", revisionNum2.String())
	assert.Equal(t, "yyy", siteCfgState.AttributeValues["change_description"])
	assert.NotEqual(t, revisionId1, revisionId2)

	//check that plan gives no diff - this actually checks the refresh and that all attributes in the state are the same as in the configuration
	plan, err = tf.Plan(context.Background())
	assert.Equal(t, nil, err)
	assert.False(t, plan) //no diff

	err = tf.StateRm(context.Background(), "qwilt_cdn_site_configuration.test")
	assert.Equal(t, nil, err)

	//import  with default revision (revision 2, latest)
	err = tf.Import(context.Background(), "qwilt_cdn_site_configuration.test", fmt.Sprintf("%s", siteId))
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(state.Values.RootModule.Resources))

	siteCfgState = findStateResource(state, "qwilt_cdn_site_configuration", "test")
	assert.NotNil(t, siteCfgState)

	revisionIdAfterImport := siteCfgState.AttributeValues["revision_id"]
	assert.Equal(t, revisionId2, revisionIdAfterImport)

	err = tf.StateRm(context.Background(), "qwilt_cdn_site_configuration.test")
	assert.Equal(t, nil, err)

	//import with explicit revision_id
	err = tf.Import(context.Background(), "qwilt_cdn_site_configuration.test", fmt.Sprintf("%s:%s", siteId, revisionId1))
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(state.Values.RootModule.Resources))

	siteCfgState = findStateResource(state, "qwilt_cdn_site_configuration", "test")
	assert.NotNil(t, siteCfgState)

	importedRevisionNum := siteCfgState.AttributeValues["revision_num"].(json.Number)
	assert.Equal(t, "1", importedRevisionNum.String())

	revisionIdAfterImport = siteCfgState.AttributeValues["revision_id"]
	assert.Equal(t, revisionId1, revisionIdAfterImport)

	//remove the configuration and check that it is destroyed
	terraformBuilder.DelSiteCfgResource("test")
	terraformConfig = terraformBuilder.Build()

	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(state.Values.RootModule.Resources))
}
