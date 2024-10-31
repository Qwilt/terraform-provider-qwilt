// Package cdn
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.
package cdn

import (
	"context"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestIpAllowListDataResource(t *testing.T) {

	t.Logf("Starting TestIpAllowListDataResource test DEBUG")

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

	terraformBuilder := NewTerraformConfigBuilder().IpAllowListDataSource("detail")
	terraformConfig := terraformBuilder.Build()

	t.Logf("ncjhdbcj: %s", terraformConfig)
	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	// Initialize a new Terraform instance
	tf, err := tfexec.NewTerraform(tempDir, tfBinaryPath)
	assert.Equal(t, nil, err)

	//varOption := tfexec.Var(fmt.Sprintf("site_id=%s", siteId))
	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	// Read the output value
	_, err = tf.Output(context.Background())
	assert.Equal(t, nil, err)

	//t.Logf("%s", output)

	// Assert that the output matches the expected value
	//var data map[string]interface{}
	//err = json.Unmarshal(output["cert_detail"].Value, &data)
	//assert.Equal(t, nil, err)
}
