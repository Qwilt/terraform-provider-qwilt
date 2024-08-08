package client

// Assuming "ga" as target environment
var targetEnvironment string = "ga"

// SiteClientFacade -
type SiteClientFacade struct {
	*PublishOpsClient
	*SiteClient
	*SiteConfigurationClient
	*SiteCertificatesClient
	ApiEndpoint string
}

// Decorator on top of Client type
func NewSiteFacadeClient(target string, client *Client) *SiteClientFacade {
	c := SiteClientFacade{
		PublishOpsClient:        NewPublishOpsClient(target, client),
		SiteClient:              NewSiteClient(target, client),
		SiteConfigurationClient: NewSiteConfigurationClient(target, client),
		SiteCertificatesClient:  NewSiteCertificatesClient(target, client),
	}
	return &c
}
