package qwiltcdn

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

func TestSiteActivationResource(t *testing.T) {

	t.Logf("Starting TestSiteActivationResource test")

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

	var curSiteName string
	var curHostName string
	generateSiteName(&curSiteName)
	generateHostName(&curHostName)

	var changeDesc = fmt.Sprintf("Terraform plugin unit testing description for site %s", curSiteName)

	// Create a temporary directory to hold the Terraform configuration
	tempDir, err = os.MkdirTemp("", "tf-exec-example-2")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory after the test

	terraformBuilder := NewTerraformConfigBuilder()
	terraformBuilder.SiteResource("test", generateSiteName(&curSiteName))
	terraformBuilder.SiteConfigResource("test", curHostName, changeDesc)
	terraformBuilder.SiteActivationResource("test")
	terraformConfig := terraformBuilder.Build()

	//t.Logf("config: %s", terraformConfig)
	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err := tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(state.Values.RootModule.Resources))

	siteState := findStateResource(state, "qwiltcdn_site", "test")
	siteConfigState := findStateResource(state, "qwiltcdn_site_configuration", "test")
	siteActivationState := findStateResource(state, "qwiltcdn_site_activation", "test")
	assert.NotNil(t, siteState)
	assert.NotNil(t, siteConfigState)
	assert.NotNil(t, siteActivationState)

	//get id's to test it later with import
	siteId := siteState.AttributeValues["site_id"]
	revisionId := siteConfigState.AttributeValues["revision_id"]
	publishId := siteActivationState.AttributeValues["publish_id"]

	assert.Equal(t, revisionId, siteActivationState.AttributeValues["revision_id"])
	assert.Equal(t, siteId, siteActivationState.AttributeValues["site_id"])
	assert.Equal(t, "ga", siteActivationState.AttributeValues["target"])
	assert.Equal(t, "null", siteActivationState.AttributeValues["validators_err_details"])
	assert.Equal(t, "Publish", siteActivationState.AttributeValues["operation_type"])

	//check that plan gives no diff - this actually checks the refresh and that all attributes in the state are the same as in the configuration
	plan, err := tf.Plan(context.Background())
	assert.Equal(t, nil, err)
	assert.False(t, plan) //no diff

	//wait for activation to complete
	start := time.Now()
	publishCompleted := false
	for time.Since(start) < 60*time.Second {
		tf.Refresh(context.Background())
		siteActivationState = findStateResource(state, "qwiltcdn_site_activation", "test")
		if siteActivationState.AttributeValues["publish_status"] != "InProgress" {
			publishCompleted = true
			t.Logf("publish operation %s completed, status %s", publishId, siteActivationState.AttributeValues["publish_status"])
			break
		}
		t.Logf("wait for publish operation %s completion", publishId)
		time.Sleep(3 * time.Second) // Wait for few seconds before checking again
	}
	assert.True(t, publishCompleted)
	assert.Equal(t, "Success", siteActivationState.AttributeValues["publish_status"])
	assert.Equal(t, "Accepted", siteActivationState.AttributeValues["publish_acceptance_status"])
	assert.Equal(t, "null", siteActivationState.AttributeValues["validators_err_details"])

	//prepare for import, remove existing resources from state
	err = tf.StateRm(context.Background(), "qwiltcdn_site_activation.test")
	assert.Equal(t, nil, err)

	err = tf.StateRm(context.Background(), "qwiltcdn_site_configuration.test")
	assert.Equal(t, nil, err)

	err = tf.StateRm(context.Background(), "qwiltcdn_site.test")
	assert.Equal(t, nil, err)

	//import resources using only site_id
	t.Logf("implicitly importing resources for site %s", siteId)
	err = tf.Import(context.Background(), "qwiltcdn_site.test", fmt.Sprintf("%s", siteId))
	assert.Equal(t, nil, err)

	err = tf.Import(context.Background(), "qwiltcdn_site_configuration.test", fmt.Sprintf("%s", siteId))
	assert.Equal(t, nil, err)

	err = tf.Import(context.Background(), "qwiltcdn_site_activation.test", fmt.Sprintf("%s", siteId))
	assert.Equal(t, nil, err)

	//prepare for import, remove existing resources from state
	err = tf.StateRm(context.Background(), "qwiltcdn_site_activation.test")
	assert.Equal(t, nil, err)

	err = tf.StateRm(context.Background(), "qwiltcdn_site_configuration.test")
	assert.Equal(t, nil, err)

	//import resources using explicit identifiers
	t.Logf("explicitly importing resources for site %s", siteId)
	err = tf.Import(context.Background(), "qwiltcdn_site_configuration.test", fmt.Sprintf("%s:%s", siteId, revisionId))
	assert.Equal(t, nil, err)

	err = tf.Import(context.Background(), "qwiltcdn_site_activation.test", fmt.Sprintf("%s:%s", siteId, publishId))
	assert.Equal(t, nil, err)

	//check that plan gives no diff - this actually checks the refresh and that all attributes in the state are the same as in the configuration
	plan, err = tf.Plan(context.Background())
	assert.Equal(t, nil, err)
	assert.False(t, plan) //no diff

	//remove site_activation first and site_configuration and wait
	t.Logf("removing site_activation and site_configuration resources for site %s", siteId)
	terraformBuilder.DelSiteActivationResource("test")
	terraformBuilder.DelSiteCfgResource("test")
	terraformConfig = terraformBuilder.Build()

	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	//we expect here to delete site_activation and site_configuration, site will fail to delete until it is unpublished
	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(state.Values.RootModule.Resources))

	siteState = findStateResource(state, "qwiltcdn_site", "test")
	siteConfigState = findStateResource(state, "qwiltcdn_site_configuration", "test")
	siteActivationState = findStateResource(state, "qwiltcdn_site_activation", "test")
	assert.NotNil(t, siteState)
	assert.Nil(t, siteConfigState)
	assert.Nil(t, siteActivationState)

	t.Logf("wait for un-publish operation %s completion", publishId)
	time.Sleep(10 * time.Second) // Wait for few seconds before checking again

	//finally, remove site now that there it is unpublished
	t.Logf("removing site resource for site %s", siteId)
	terraformConfig = QwiltCdnFullProviderConfig

	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Nil(t, err)
	assert.Nil(t, state.Values)

	//use data source to query the data source
	//terraformConfig = QwiltCdnFullProviderConfig + `data "qwiltcdn_sites" "test" {}`
	//err = ioutil.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	//err = tf.Apply(context.Background())
	//assert.Equal(t, nil, err)

	//// (Optional) Run Terraform destroy to clean up resources
	//if err := tf.Destroy(context.Background()); err != nil {
	//	log.Fatalf("Error destroying Terraform-managed infrastructure: %s", err)
	//}
}
