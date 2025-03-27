// Package client
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.

package client_test

import (
	"encoding/json"
	"fmt"
	"github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/api"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"

	cdnclient "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/client"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website.

// If you do not have terraform installed, you can remove the formatting command, but it is suggested to
// ensure the documentation is formatted properly.

// these will be set by the goreleaser configuration
// to appropriate values for the compiled binary.
var version string = "dev"

// goreleaser can pass other information to the main package, such as the specific commit
// https://goreleaser.com/cookbooks/using-main.version/

func getSites(t *testing.T, client *cdnclient.SiteClientFacade, includeActiveLastPub bool) {
	sites, err := client.GetSites(includeActiveLastPub, false)
	if err != nil {
		log.Fatal(err.Error())
	}

	assert.Greater(t, 100, len(sites))
	//for _, site := range sites {
	//	// fmt.Printf("%s: %s\n", site.SiteId, site.SiteName)
	//	sd, err := json.Marshal(site)
	//	if err != nil {
	//		return
	//	}
	//	fmt.Printf("%s\n\n", string(sd))
	//}
}

func getSiteIdByName(client *cdnclient.SiteClientFacade, name string) (string, error) {
	sites, err := client.GetSites(false, false)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, site := range sites {
		index := strings.Index(site.SiteName, name)
		if index != -1 {
			return site.SiteId, nil
		}
	}
	return "", fmt.Errorf("error site not found")
}

func getSite(client *cdnclient.SiteClientFacade, siteId string, includeActiveLastPub bool) {
	siteDetails, err := client.GetSite(siteId, cdnclient.TARGET_GA, includeActiveLastPub, false)
	if err != nil {
		log.Fatal(err.Error())
	}

	sd, err := json.Marshal(siteDetails)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", string(sd))
}

func createSite(client *cdnclient.SiteClientFacade, siteName string, routingMethod string) {
	var newSite api.SiteCreateRequest
	newSite.SiteName = siteName
	newSite.RoutingMethod = routingMethod

	// Create new site
	site, err := client.CreateSite(newSite)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%s: %s\n", site.SiteId, site.SiteName)
}

func updateSite(client *cdnclient.SiteClientFacade, siteId string, siteName string) {
	var siteDetails api.SiteUpdateRequest
	siteDetails.SiteName = siteName

	siteResults, err := client.UpdateSite(siteId, siteDetails)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%s: %s\n", siteResults.SiteId, siteResults.SiteName)
}

func deleteSite(client *cdnclient.SiteClientFacade, siteId string, name string) {
	err := client.DeleteSite(siteId)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Printf("Site %s deleted successfully\n", siteId)
	}
}

func getSiteConfigs(client *cdnclient.SiteClientFacade, siteId string, truncateHostIndex bool) {
	siteConfigs, err := client.GetSiteConfigs(siteId, truncateHostIndex)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, siteConfig := range siteConfigs {
		fmt.Printf("%s: %s %d\n", siteConfig.SiteId, siteConfig.RevisionId, siteConfig.RevisionNum)
		fmt.Printf("%s\n", siteConfig.HostIndex)
	}
}

func getSiteConfig(client *cdnclient.SiteClientFacade, siteId string, revisionId string) {
	siteConfigVersion, err := client.GetSiteConfig(siteId, revisionId, false)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%s: %s %d\n", siteConfigVersion.SiteId, siteConfigVersion.RevisionId, siteConfigVersion.RevisionNum)
	fmt.Printf("%s", siteConfigVersion.HostIndex)
}

func createSiteConfig(client *cdnclient.SiteClientFacade, siteId string, hostIndex json.RawMessage, changeDescription string) {
	var newSiteConfig api.SiteConfigAddRequest
	newSiteConfig.ChangeDescription = changeDescription
	newSiteConfig.HostIndex = hostIndex

	siteConfigResponse, err := client.CreateSiteConfig(siteId, newSiteConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%s: %s %d\n", siteConfigResponse.SiteId, siteConfigResponse.RevisionId, siteConfigResponse.RevisionNum)
}

func findLatestPubOp(client *cdnclient.SiteClientFacade, siteId string, revisionId string) {
	pubOp, err := client.FindLatestPubOp(siteId, revisionId)
	if err != nil {
		log.Fatal(err.Error())
	}
	if pubOp.PublishId != "" {
		fmt.Printf("%s: %s %s %s isActive=%t is latest\n", pubOp.PublishId, pubOp.OperationType, pubOp.PublishState, pubOp.PublishStatus, pubOp.IsActive)
	} else {
		fmt.Println("No latest pubOp found")
	}
}

func getPubOp(client *cdnclient.SiteClientFacade, siteId string, publishId string) {
	pubOp, err := client.GetPubOp(siteId, publishId)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%s: %s %s %s isActive=%t\n", pubOp.PublishId, pubOp.OperationType, pubOp.PublishState, pubOp.PublishStatus, pubOp.IsActive)
}

func getPubOps(client *cdnclient.SiteClientFacade, siteId string, isActive bool, publishState string) {
	pubOps, err := client.GetPubOps(siteId, isActive, publishState)
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, pubOp := range pubOps {
		fmt.Printf("%s: %s %s %s %s isActive=%t\n", pubOp.PublishId, pubOp.RevisionId, pubOp.OperationType, pubOp.PublishState, pubOp.PublishStatus, pubOp.IsActive)
	}
}

func getPubStatus(client *cdnclient.SiteClientFacade, siteId string) {
	currentStatus, transitionStatus, err := client.GetSitePubStatus(siteId)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Status: %s %s\n", currentStatus, transitionStatus)
}

func publish(client *cdnclient.SiteClientFacade, siteId string, revisionId string) {
	target := cdnclient.TARGET_GA
	pubOp, err := client.Publish(siteId, revisionId, target)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%s: %s %s\n", pubOp.PublishId, pubOp.PublishState, pubOp.PublishStatus)
}

func unpublish(client *cdnclient.SiteClientFacade, siteId string) {
	target := cdnclient.TARGET_GA
	pubOp, err := client.Unpublish(siteId, target)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%s: %s %s\n", pubOp.PublishId, pubOp.PublishState, pubOp.PublishStatus)
}

func auth(envType string, username string, password string, token string) (client *cdnclient.Client) {
	client, err := cdnclient.NewClient(envType, username, password, token)
	if err != nil {
		log.Fatal(err.Error())
	}

	if token == "" {
		fmt.Printf("export QCDN_API_KEY='%s'\n", client.Token)
	}

	return client
}

//func TestClient(t *testing.T) {
//	var debug bool
//
//	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
//	flag.Parse()
//
//	// fmt.Println("Qwilt API test")
//
//	envtype := "dev" //os.Getenv("QCDN_ENVTYPE")
//	username := os.Getenv("QCDN_USERNAME")
//	password := "" // Only doing token authentication; eval $(go run token.go)
//	token := os.Getenv("QCDN_API_KEY")
//
//	// siteId := "65ea13f8bda5005ce98807e5"
//	// revisionId := "65ea1411bda5005ce98807e6"
//	// publishId := "5f542265-a012-4c07-a719-45922c129da1"
//	// siteId := "65ea1a1832cd0d7df30dc1d2" //"65e76a1f32cd0d7df30dc1c9"
//	// revisionId := "65ea1a19bda5005ce98807e7" //"65e76a1f32cd0d7df30dc1ca"
//	// publishId := "6468c307-5d96-4b1a-b86d-3a461e1da89c"
//	siteId := "65fc907e1c17716a2c46883d"
//	//revisionId := "65f9db4a4a1c9c72079eb7b6"
//	//publishId := "90fef06f-4d7e-48ef-9950-81919504f11e"
//	//siteName := "jayr-API-test-site9"
//	//routingMethod := "DNS"
//	hostIndex := json.RawMessage(`{
//	"hosts": [
//		  {
//			"host": "www.examplesite-manual1.com",
//			"host-metadata": {
//			  "metadata": [
//				{
//				  "generic-metadata-type": "MI.SourceMetadataExtended",
//				  "generic-metadata-value": {
//					"sources": [
//					  {
//						"protocol": "https/1.1",
//						"endpoints": [
//						  "www.example-origin-host.com"
//						]
//					  }
//					]
//				  }
//				}
//			  ],
//			  "paths": []
//			}
//		  }
//		]
//	}`)
//	changeDescription := "Testing out the API with a Go-based library"
//
//	client := auth(envtype, username, password, token)
//	siteClient := cdnclient.NewSiteFacadeClient(api.SITES_HOSTNAME, client)
//
//	/* Select API calls to test here */
//	//var curSiteName string
//	//generateSiteName(curSiteName)
//	siteName := "test-client"
//	siteId, err := getSiteIdByName(siteClient, siteName)
//	if err == nil {
//		t.Logf("found sidte id %s for name %s", siteId, siteName)
//		deleteSite(siteClient, siteId, siteName)
//	}
//	createSite(siteClient, siteName, "DNS")
//	siteId, err = getSiteIdByName(siteClient, siteName)
//	assert.Nil(t, err)
//
//	createSiteConfig(siteClient, siteId, hostIndex, changeDescription)
//	getSiteConfigs(siteClient, siteId, false)
//
//	getSites(t, siteClient, true)
//	//getSite(siteClient, siteId, true)
//	// createSite(client, siteName, routingMethod)
//	// updateSite(client, siteId, siteName + " DELETED")
//	// deleteSite(client, siteId)
//	// nukeSite(client, siteId)
//	// getSiteConfigs(client, siteId, false)
//	// getSiteConfig(client, siteId, revisionId)
//	// createSiteConfig(client, siteId, hostIndex, changeDescription)
//	// getPubOps(client, siteId, "", "")
//	// findLatestPubOp(client, siteId, revisionId)
//	// getPubOp(client, siteId, publishId)
//	// getPubStatus(client, siteId)
//	// publish(client, siteId, revisionId)
//	// unpublish(client, siteId)
//	// republish(client, siteId)
//
//	/* Site clean-up */
//	/*
//		siteIdList := []string{"65f20a3d2628f55e6aae2ab3","65f21a7f2628f55e6aae2ab9"}
//		for _, siteToNuke := range siteIdList {
//			fmt.Printf("Cleaning up siteId %s\n", siteToNuke)
//			nukeSite(client, siteToNuke)
//		}
//	*/
//
//	/* Qwell unused var warnings */
//	//fmt.Println(
//	//	siteId,
//	//	revisionId,
//	//	publishId,
//	//	siteName,
//	//	routingMethod,
//	//	hostIndex,
//	//	changeDescription,
//	//)
//}
