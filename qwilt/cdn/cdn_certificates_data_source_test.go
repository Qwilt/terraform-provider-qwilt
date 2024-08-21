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
	"strconv"
	"testing"
)

func TestCertificatesDataResource(t *testing.T) {

	t.Logf("Starting TestCertificatesDataResource test DEBUG")

	//os.Setenv("TF_CLI_CONFIG_FILE", "/Users/efrats/.terraformrc")

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

	terraformBuilder := NewTerraformConfigBuilder().CertResource("test", test_domain_key, test_domain_crt, "ccc")
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
	certId := fmt.Sprintf("%s", certState.AttributeValues["cert_id"])

	tempDir2, err := os.MkdirTemp("", "tf-exec-example-data-sources")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir2) // Clean up the temporary directory after the test

	terraformBuilder2 := NewTerraformConfigBuilder().CertsDataSource("detail", certId)
	terraformConfig2 := terraformBuilder2.Build()

	t.Logf("ncjhdbcj: %s", terraformConfig2)
	tfFilePath2 := tempDir2 + "/main.tf"
	err = os.WriteFile(tfFilePath2, []byte(terraformConfig2), 0644)
	assert.Equal(t, nil, err)

	// Initialize a new Terraform instance
	tf2, err := tfexec.NewTerraform(tempDir2, tfBinaryPath)
	assert.Equal(t, nil, err)

	//varOption := tfexec.Var(fmt.Sprintf("site_id=%s", siteId))
	err = tf2.Apply(context.Background())
	assert.Equal(t, nil, err)

	// Read the output value
	output, err := tf2.Output(context.Background())
	assert.Equal(t, nil, err)

	//t.Logf("type: %s", output["site_detail"].Type)

	// Assert that the output matches the expected value
	var data map[string]interface{}
	err = json.Unmarshal(output["cert_detail"].Value, &data)
	assert.Equal(t, nil, err)

	assert.Equal(t, certId, strconv.Itoa(int(data["cert_id"].(float64))))

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
