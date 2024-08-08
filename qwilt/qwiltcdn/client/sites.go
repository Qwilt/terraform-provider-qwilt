package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Qwilt/terraform-provider-qwilt/qwilt/qwiltcdn/api"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SiteClient struct {
	*Client
	apiEndpoint string
}

func NewSiteClient(target string, client *Client) *SiteClient {
	c := SiteClient{
		Client:      client,
		apiEndpoint: client.endpointBuilder.Build(target),
	}
	return &c
}

// GetSites - Returns list of sites
func (c *SiteClient) GetSites(includeActiveLastPub bool, includeDeletedSites bool) ([]api.Site, error) {
	var param string
	if includeActiveLastPub {
		param = "?includePublishDetails=true"
	} else {
		param = ""
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/sites%s", c.apiEndpoint, param), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sites := []api.Site{}
	err = json.Unmarshal(body, &sites)
	if err != nil {
		return nil, err
	}

	if includeDeletedSites {
		return sites, nil
	}

	//filter deleted sites
	nonDeletedSites := []api.Site{}
	for _, site := range sites {
		if site.IsDeleted {
			continue
		}
		nonDeletedSites = append(nonDeletedSites, site)
	}
	return nonDeletedSites, nil
}

// GetSite - Returns site details
func (c *SiteClient) GetSite(siteId string, target string, includeActiveLastPub bool, includeDeletedSites bool) (*api.Site, error) {
	if siteId == "" {
		return nil, fmt.Errorf("siteId is empty")
	}

	var queryParams string
	if includeActiveLastPub && target != "" {
		queryParams = fmt.Sprintf("?includePublishDetails=true&publishTarget=%s", target)
	} else if includeActiveLastPub && target == "" {
		queryParams = fmt.Sprintf("?includePublishDetails=true&publishTarget=ga")
	} else if !includeActiveLastPub && target != "" {
		queryParams = fmt.Sprintf("?publishTarget=%s", target)
	} else {
		queryParams = ""
	}

	url := fmt.Sprintf("%s/api/v2/sites/%s%s", c.apiEndpoint, siteId, queryParams)
	log.Printf("***** %s *****", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	siteDetail := api.Site{}
	err = json.Unmarshal(body, &siteDetail)
	if err != nil {
		return nil, err
	}

	if !includeDeletedSites && siteDetail.IsDeleted {
		return nil, errors.New("Site was found but marked for deletion")
	}
	return &siteDetail, nil
}

// CreateSite - Create new site
func (c *SiteClient) CreateSite(site api.SiteCreateRequest) (*api.Site, error) {
	rb, err := json.Marshal(site)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v2/sites", c.apiEndpoint), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newSite := api.Site{}
	err = json.Unmarshal(body, &newSite)
	if err != nil {
		return nil, err
	}

	return &newSite, nil
}

// UpdateSite - Update site details
func (c *SiteClient) UpdateSite(siteId string, site api.SiteUpdateRequest) (*api.Site, error) {
	if siteId == "" {
		return nil, fmt.Errorf("siteId is empty")
	}

	rb, err := json.Marshal(site)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v2/sites/%s", c.apiEndpoint, siteId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedSite := api.Site{}
	err = json.Unmarshal(body, &updatedSite)
	if err != nil {
		return nil, err
	}

	return &updatedSite, nil
}

// DeleteSite - Deletes a site
func (c *SiteClient) DeleteSite(siteId string) error {
	if siteId == "" {
		return fmt.Errorf("siteId is empty")
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v2/sites/%s", c.apiEndpoint, siteId), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	_ = body
	if err != nil {
		return err
	}

	return nil
}

// DeleteAndRenameSite - Deletes a site so we can reuse the name later
//
//	This adds "DELETED <unixtimestamp>" to the end of the site name to make it unique
func (c *SiteClient) DeleteAndRenameSite(siteId string, siteName string) error {
	if siteId == "" || siteName == "" {
		return fmt.Errorf("DeleteAndRenameSite: Invalid input, siteId=%s siteName=%s", siteId, siteName)
	}

	err := c.DeleteSite(siteId)
	if err != nil {
		return fmt.Errorf("DeleteAndRenameSite delete API call failed: %s", err)
	}

	site := api.SiteUpdateRequest{
		SiteName: fmt.Sprintf("%s DELETED %s", siteName,
			strconv.FormatInt(time.Now().UTC().UnixNano(), 10)),
	}

	_, err = c.UpdateSite(siteId, site)
	if err != nil {
		return fmt.Errorf("DeleteAndRenameSite update API call failed: %s", err)
	}

	return nil
}
