package qwiltcdn

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Qwilt/terraform-provider-qwilt/qwilt/qwiltcdn/api"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCoffeesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: QwiltCdnProviderConfig + `data "qwiltcdn_sites" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Ideally, I would like to see this:
					// 1. Use the resource to provision a site
					// 2. Look up the site and revision IDs of the site
					// 3. Verify that the attributes exist for that site and are correct
					// 4. Remove that site
					// But, that will require some more work.
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.#"),
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.0.created_user"),
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.0.creation_time_milli"),
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.0.is_self_service_blocked"),
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.0.last_update_time_milli"),
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.0.last_updated_user"),
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.0.owner_org_id"),
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.0.routing_method"),
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.0.service_id"),
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.0.service_type"),
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.0.should_provision_to_third_party_cdn"),
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.0.site_dns_cname_delegation_target"),
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.0.site_id"),
					resource.TestCheckResourceAttrSet("data.qwiltcdn_sites.test", "site.0.site_name"),
				),
			},
		},
	})
}

func TestSitesDataResource(t *testing.T) {

	t.Logf("Starting TestSitesDataResource test DEBUG")

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

	var curSiteName string

	terraformBuilder := NewTerraformConfigBuilder().SiteResource("test", generateSiteName(&curSiteName))
	terraformConfig := terraformBuilder.Build()
	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	// Initialize a new Terraform instance
	tf, err := tfexec.NewTerraform(tempDir, tfBinaryPath)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	//get the site_id to test it later with data-source
	state, err := tf.Show(context.Background())
	siteState := findStateResource(state, "qwiltcdn_site", "test")
	siteId := fmt.Sprintf("%s", siteState.AttributeValues["site_id"])

	tempDir2, err := os.MkdirTemp("", "tf-exec-example-data-sources")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir2) // Clean up the temporary directory after the test

	terraformBuilder2 := NewTerraformConfigBuilder().SitesDataResource("detail", siteId)
	terraformConfig2 := terraformBuilder2.Build()

	t.Logf("ncjhdbcj: %s", terraformConfig2)
	tfFilePath2 := tempDir2 + "/main.tf"
	err = os.WriteFile(tfFilePath2, []byte(terraformConfig2), 0644)
	assert.Equal(t, nil, err)

	// Initialize a new Terraform instance
	tf2, err := tfexec.NewTerraform(tempDir2, tfBinaryPath)
	assert.Equal(t, nil, err)

	err = tf2.Apply(context.Background())
	assert.Equal(t, nil, err)

	// Read the output value
	output, err := tf2.Output(context.Background())

	assert.Equal(t, nil, err)

	// Assert that the output matches the expected value
	site := api.Site{}
	err = json.Unmarshal(output["site_detail"].Value, &site)
	assert.Equal(t, nil, err)

	//t.Logf("site: %s", site.SiteName)
	//assert.Equal(t, siteId, site.SiteId)

}
