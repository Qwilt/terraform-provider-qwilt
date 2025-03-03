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
	"fmt"
	"strconv"
	"strings"

	"github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/api"
	cdnclient "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/client"
	cdnmodel "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/model"
	"github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/validators"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &siteActivationResource{}
	_ resource.ResourceWithConfigure   = &siteActivationResource{}
	_ resource.ResourceWithImportState = &siteActivationResource{}
)

// NewSiteActivationResource is a helper function to simplify the provider implementation.
func NewSiteActivationResource() resource.Resource {
	return &siteActivationResource{
		client: nil,
		target: cdnclient.TARGET_GA,
	}
}

// siteActivationResource is the resource implementation.
type siteActivationResource struct {
	client *cdnclient.SiteClientFacade
	target string
}

// Metadata returns the resource type name.
func (r *siteActivationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_site_activation"
}

// Schema defines the schema for the resource.
func (r *siteActivationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Qwilt CDN site activation and certificate assignment.<br><br>" +
			"Notes:<br>" +
			" - This resource takes a long time to fully apply.<br>" +
			" - If a site activation attempt fails, it may be due to another publish operation already in progress for the same site_id.<br>" +
			" - Run ```terraform refresh``` to sync the state of this resource explicitly.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "For internal use only, for testing. Equals site_id:publish_id.",
				Computed:    true,
			},
			"site_id": schema.StringAttribute{
				Description: "SiteId of the site the user wishes to publish or unpublish.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"revision_id": schema.StringAttribute{
				Description: "Unique identifier of the configuration version that was published or unpublished.",
				Required:    true,
			},
			"certificate_id": schema.Int64Attribute{
				Description: "The ID of the certificate you want to link to this site. Cannot co-exist with certificate_template_id.",
				Optional:    true,
				Validators: []validator.Int64{
					validators.NewMutualExclusiveValidator(path.Root("certificate_template_id")),
				},
			},
			"certificate_template_id": schema.Int64Attribute{
				Description: "The ID of the certificate template you want to link to this site. Cannot co-exist with certificate_id.",
				Optional:    true,
				Validators: []validator.Int64{
					validators.NewMutualExclusiveValidator(path.Root("certificate_id")),
				},
			},
			"publish_id": schema.StringAttribute{
				Description: "The ID of the publishing operation for which you want to retrieve metadata.",
				Computed:    true,
			},
			"creation_time_milli": schema.Int64Attribute{
				Description: "The time when the publish operation was created, in epoch time.",
				Computed:    true,
			},
			"owner_org_id": schema.StringAttribute{
				Description: "The organization that owns the site.",
				Computed:    true,
			},
			"last_update_time_milli": schema.Int64Attribute{
				Description: "When the publishing operation was last updated, in epoch time.",
				Computed:    true,
			},
			"target": schema.StringAttribute{
				Description: "The value will always be 'ga'.",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username that initiated the publishing operation.",
				Computed:    true,
			},
			"publish_state": schema.StringAttribute{
				Description: "For internal use. Use the 'publishStatus' values instead.",
				Computed:    true,
			},
			"publish_status": schema.StringAttribute{
				Description: "The publishing operation status. The 'publishStatus' values aggregate the 'publishState' values into broader categories. \n\n" +
					"  - Success - The operation succeeded.\n" +
					"  - Failed - The operation failed.\n" +
					"  - Aborted - The operation was canceled.\n" +
					"  - InProgress - The operation is in progress.",
				Computed: true,
			},
			"publish_acceptance_status": schema.StringAttribute{
				Description: "The CDN validates and then accepts the publishing operation before initiating it. This attribute lets you track the acceptance process. It is not an indication of the status of the publishing operation on the CDN caches themselves.\n\n" +
					"  - Pending - Pending validation.\n" +
					"  - Invalid - Validation failed. \n" +
					"  - Dismissed - Passed validation but failed to initiate.\n" +
					"  - Aborted - The publishing operation was cancelled.\n" +
					"  - In progress - The publishing operation is validated and pending initiation.\n" +
					"  - Accepted - The publishing operation was initiated.",
				Computed: true,
			},
			"operation_type": schema.StringAttribute{
				Description: "The operation type (Publish, Unpublish) that was initiated with the request. An Unpublish operation removes a delivery service from the CDN.",
				Computed:    true,
			},
			//"status_line": schema.ListAttribute{
			//	ElementType: types.StringType,
			//	Description: "Additional information related to the publish status.",
			//	Computed:    true,
			//},
			"is_active": schema.BoolAttribute{
				Description: "Indicates if the configuration is active or inactive.",
				Computed:    true,
			},
			"validators_err_details": schema.StringAttribute{
				Description: "Details about errors generated during validation.",
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *siteActivationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cdnmodel.SiteActivation
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "siteActivationResource: create")

	// Evaluate the certificate ID
	var certificateId int64
	switch {
	// If the certificate ID is set, use it
	case !plan.CertificateId.IsNull():
		certificateId = plan.CertificateId.ValueInt64()
	// If the certificate ID is not set and certificate template ID is set, get the certificate ID
	case !plan.CertificateTemplateId.IsNull():
		certificateTemplate, err := r.client.GetCertificateTemplate(plan.CertificateTemplateId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Getting Certificate Template",
				"Could not get certificate template for Qwilt CDN Site, unexpected error: "+err.Error(),
			)
			return
		}

		// If the certificate template doesn't have a last certificate ID, the certificate template is still pending verification.
		// Inform the user and return an error.
		if certificateTemplate.LastCertificateID == nil {
			if certificateTemplate.AutoManagedCertificateTemplate {
				domainsList, err := r.client.GetChallengeDelegationDomainsListFromCertificateTemplateId(plan.CertificateTemplateId)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error Getting Challenge Delegation Domains List",
						"Could not get challenge delegation domains list for Qwilt CDN Certificate Template",
					)
					return
				}
				resp.Diagnostics.AddError(
					"Chosen Certificate Template is pending verification",
					fmt.Sprintf("Certificate Template is pending verification. Please make sure to have the CNAMEs list configured correctly:\n%v", domainsList.PrettyPrint()),
				)
				return
			}
			resp.Diagnostics.AddError(
				"Chosen Certificate Template is missing Certificate",
				fmt.Sprintf("Please upload a Certificate to the Certificate Template"),
			)
			return
		}
		certificateId = *certificateTemplate.LastCertificateID
	}

	// If we have an https site to publish, link the certificate to the site
	if certificateId != 0 {
		_, err := r.client.LinkSiteCertificate(plan.SiteId.ValueString(), strconv.Itoa(int(certificateId)))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Linking Certificate to Qwilt CDN Site",
				"Could not link certificate to Qwilt CDN Site, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Publish the site
	pubOpResp, err := r.client.Publish(plan.SiteId.ValueString(), plan.RevisionId.ValueString(), r.target)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Publishing Qwilt CDN Site",
			"Could not Publishing Qwilt CDN Site, unexpected error: "+err.Error(),
		)
		return
	}

	timeout := cdnclient.ACCEPTANCE_TIMEOUT
	pubOpResp, err = r.client.GetAndWaitForPubOpAcceptance(plan.SiteId.ValueString(), pubOpResp.PublishId, timeout) // Function that checks status of 'x' from backend
	if err != nil {
		resp.Diagnostics.AddError(
			"Timeout while Waiting for validation status in Qwilt CDN Site Publish operation",
			"Could not get Qwilt CDN Site Publish acceptance status. err: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "siteActivationResource: PUBLISH ACCEPTANCE STATUS after timeout IS: "+pubOpResp.PublishAcceptanceStatus+"\n")
	if pubOpResp.PublishAcceptanceStatus == cdnclient.ACCEPTANCE_STATUS_INVALID ||
		pubOpResp.PublishAcceptanceStatus == cdnclient.ACCEPTANCE_STATUS_DISMISSED {
		details := fmt.Sprintf("Publish failed for Qwilt CDN Site %s\n. Acceptance Status: %s\n. Err: %s\n. Status line: %s\b",
			plan.SiteId,
			pubOpResp.PublishAcceptanceStatus,
			pubOpResp.ValidatorsErrDetails,
			strings.Join(pubOpResp.StatusLine, ","))
		resp.Diagnostics.AddError(
			"Error during PUBLISH for Qwilt CDN Site", details)

		return
	}

	// Map response body to schema and populate Computed attribute values
	newPlan := cdnmodel.NewSiteActivationBuilder().
		Ctx(ctx).
		PublishId(pubOpResp.PublishId).
		RevisionId(pubOpResp.RevisionId).
		SiteId(plan.SiteId.ValueString()).
		Username(pubOpResp.Username).
		CreationTimeMilli(pubOpResp.CreationTimeMilli).
		OwnerOrgId(pubOpResp.OwnerOrgId).
		LastUpdateTimeMilli(pubOpResp.LastUpdateTimeMilli).
		CertificateId(plan.CertificateId.ValueInt64()).
		CertificateTemplateId(plan.CertificateTemplateId.ValueInt64()).
		PublishState(pubOpResp.PublishState).
		OperationType(pubOpResp.OperationType).
		Target(pubOpResp.Target).
		PublishStatus(pubOpResp.PublishStatus).
		AcceptanceStatus(pubOpResp.PublishAcceptanceStatus).
		OperationType(pubOpResp.OperationType).
		ValidateErrDetails(pubOpResp.ValidatorsErrDetails).
		IsActive(pubOpResp.IsActive).
		//StatusLine(pubOpResp.StatusLine).
		Build()

	// Set state to fully populated data
	diags = resp.State.Set(ctx, newPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *siteActivationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state cdnmodel.SiteActivation
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "siteActivationResource: read")

	// Get refreshed site value from CDN
	pubOpResp, err := r.client.GetPubOp(state.SiteId.ValueString(), state.PublishId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Qwilt CDN Site Publish",
			"Could not read Qwilt CDN Site Publish "+state.SiteId.ValueString()+": "+err.Error(),
		)
		return
	}
	certsResp, err := r.client.GetSiteCertificates(state.SiteId.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Certificates for Qwilt CDN Site",
			"Could not get certificates for Qwilt CDN Site, unexpected error: "+err.Error(),
		)
		return
	}

	var certId int64
	var certificateTemplateId int64
	if len(certsResp) > 0 {
		certId, err = strconv.ParseInt(certsResp[0].CertificateId, 10, 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Converting Certificate Id",
				"Could not convert certificate ID for Qwilt CDN Site, unexpected error: "+err.Error(),
			)
			return
		}

		certResp, err := r.client.GetCertificate(types.Int64Value(certId), false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Getting Certificate for Qwilt CDN Site",
				"Could not get certificate for Qwilt CDN Site, unexpected error: "+err.Error(),
			)
			return
		}
		if certResp.CsrId != nil {
			csrId, err := strconv.Atoi(*certResp.CsrId)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Converting Certificate Id",
					"Could not convert certificate ID for Qwilt CDN Site, unexpected error: "+err.Error(),
				)
			}
			csrResp, err := r.client.GetCertificateSigningRequest(types.Int64Value(int64(csrId)))
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Getting Certificate Signing Request for Qwilt CDN Site",
					"Could not get certificate signing request for Qwilt CDN Site, unexpected error: "+err.Error(),
				)
				return
			}

			// If this is an auto-managed CSR, we should use the certificate template ID in the Site Activation
			if csrResp.AutoManagedCSR {
				certificateTemplateIdInt, err := strconv.Atoi(csrResp.CertificateTemplateIDRef)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error Converting Certificate Template Id",
						"Could not convert certificate template ID for Qwilt CDN Site, unexpected error: "+err.Error(),
					)
					return
				}
				certificateTemplateId = int64(certificateTemplateIdInt)
			}
		}
	}

	// Overwrite items with refreshed state
	state = cdnmodel.NewSiteActivationBuilder().
		Ctx(ctx).
		PublishId(pubOpResp.PublishId).
		RevisionId(pubOpResp.RevisionId).
		SiteId(state.SiteId.ValueString()).
		Username(pubOpResp.Username).
		CreationTimeMilli(pubOpResp.CreationTimeMilli).
		OwnerOrgId(pubOpResp.OwnerOrgId).
		LastUpdateTimeMilli(pubOpResp.LastUpdateTimeMilli).
		CertificateId(certId).
		CertificateTemplateId(certificateTemplateId).
		PublishState(pubOpResp.PublishState).
		OperationType(pubOpResp.OperationType).
		Target(pubOpResp.Target).
		PublishStatus(pubOpResp.PublishStatus).
		AcceptanceStatus(pubOpResp.PublishAcceptanceStatus).
		OperationType(pubOpResp.OperationType).
		IsActive(pubOpResp.IsActive).
		ValidateErrDetails(pubOpResp.ValidatorsErrDetails).
		//StatusLine(pubOpResp.StatusLine).
		Build()

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *siteActivationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan cdnmodel.SiteActivation
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "siteActivationResource: update")

	// Retrieve values from state
	var state cdnmodel.SiteActivation
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var lastCertificateId int64
	switch {
	case !state.CertificateId.IsNull():
		lastCertificateId = state.CertificateId.ValueInt64()
	case !state.CertificateTemplateId.IsNull():
		certificateTemplate, err := r.client.GetCertificateTemplate(state.CertificateTemplateId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Getting Certificate Template",
				"Could not get certificate template for Qwilt CDN Site, unexpected error: "+err.Error(),
			)
			return
		}

		lastCertificateId = *certificateTemplate.LastCertificateID
	}

	var newCertificateId int64
	switch {
	case !plan.CertificateId.IsNull():
		newCertificateId = plan.CertificateId.ValueInt64()
	case !plan.CertificateTemplateId.IsNull():
		certificateTemplate, err := r.client.GetCertificateTemplate(plan.CertificateTemplateId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Getting Certificate Template",
				"Could not get certificate template for Qwilt CDN Site, unexpected error: "+err.Error(),
			)
			return
		}
		newCertificateId = *certificateTemplate.LastCertificateID
	}

	if lastCertificateId != newCertificateId {
		if lastCertificateId != 0 {
			//unlink previous certificate
			err := r.client.UnLinkSiteCertificate(state.SiteId.ValueString(), strconv.Itoa(int(lastCertificateId)))
			if err != nil {
				resp.Diagnostics.AddError(
					"Error UnLinking Certificate to Qwilt CDN Site",
					"Could not unlink certificate to Qwilt CDN Site, unexpected error: "+err.Error(),
				)
				return
			}
		}
		if newCertificateId != 0 {
			//there is a certificate that should be linked
			_, err := r.client.LinkSiteCertificate(plan.SiteId.ValueString(), strconv.Itoa(int(newCertificateId)))
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Linking Certificate to Qwilt CDN Site",
					"Could not link certificate to Qwilt CDN Site, unexpected error: "+err.Error(),
				)
				return
			}
		}
	}

	// Publish the site
	pubOpResp, err := r.client.Publish(plan.SiteId.ValueString(), plan.RevisionId.ValueString(), r.target)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Publishing Qwilt CDN Site",
			"Could not Publishing Qwilt CDN Site, unexpected error: "+err.Error(),
		)
		return
	}

	timeout := cdnclient.ACCEPTANCE_TIMEOUT
	pubOpResp, err = r.client.GetAndWaitForPubOpAcceptance(plan.SiteId.ValueString(), pubOpResp.PublishId, timeout) // Function that checks status of 'x' from backend
	if err != nil {
		resp.Diagnostics.AddError(
			"Timeout while Waiting for validation status in Qwilt CDN Site Publish operation",
			"Could not get Qwilt CDN Site Publish acceptance status. err: "+err.Error(),
		)
		return
	}
	if pubOpResp.PublishAcceptanceStatus == cdnclient.ACCEPTANCE_STATUS_INVALID ||
		pubOpResp.PublishAcceptanceStatus == cdnclient.ACCEPTANCE_STATUS_DISMISSED {
		details := fmt.Sprintf("Publish failed for Qwilt CDN Site %s\n. Acceptance Status: %s\n. Err: %s\n. Status line: %s\b",
			plan.SiteId,
			pubOpResp.PublishAcceptanceStatus,
			pubOpResp.ValidatorsErrDetails,
			strings.Join(pubOpResp.StatusLine, ","))
		resp.Diagnostics.AddError(
			"Error during PUBLISH for Qwilt CDN Site", details)
		return
	}
	tflog.Info(ctx, "siteActivationResource: PUBLISH ACCEPTANCE STATUS after timeout IS: "+pubOpResp.PublishAcceptanceStatus+"\n")

	// Map response body to schema and populate Computed attribute values
	newPlan := cdnmodel.NewSiteActivationBuilder().
		Ctx(ctx).
		PublishId(pubOpResp.PublishId).
		RevisionId(pubOpResp.RevisionId).
		SiteId(plan.SiteId.ValueString()).
		Username(pubOpResp.Username).
		CreationTimeMilli(pubOpResp.CreationTimeMilli).
		OwnerOrgId(pubOpResp.OwnerOrgId).
		LastUpdateTimeMilli(pubOpResp.LastUpdateTimeMilli).
		CertificateId(plan.CertificateId.ValueInt64()).
		CertificateTemplateId(plan.CertificateTemplateId.ValueInt64()).
		PublishState(pubOpResp.PublishState).
		PublishStatus(pubOpResp.PublishStatus).
		AcceptanceStatus(pubOpResp.PublishAcceptanceStatus).
		OperationType(pubOpResp.OperationType).
		Target(pubOpResp.Target).
		IsActive(pubOpResp.IsActive).
		ValidateErrDetails(pubOpResp.ValidatorsErrDetails).
		//StatusLine(pubOpResp.StatusLine).
		Build()

	// Set state to fully populated data
	diags = resp.State.Set(ctx, newPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete just removes the Terraform state on success. No real deletion for this object
func (r *siteActivationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state cdnmodel.SiteActivation
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "siteActivationResource: delete, publish status: "+state.PublishStatus.ValueString())

	//Deletion semantic is 'unpublish'
	tflog.Info(ctx, "siteActivationResource: UN-PUBLISH for publish-id: "+state.PublishId.ValueString())

	_, err := r.client.Unpublish(state.SiteId.ValueString(), state.Target.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error UnPublishing Qwilt CDN Site",
			"Could not UnPublish Qwilt CDN Site "+state.SiteId.ValueString()+": "+err.Error(),
		)
		return
	}

	if state.CertificateId.ValueInt64() != 0 {
		//unlink previous certificate
		err := r.client.UnLinkSiteCertificate(state.SiteId.ValueString(), strconv.Itoa(int(state.CertificateId.ValueInt64())))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error UnLinking Certificate to Qwilt CDN Site",
				"Could not unlink certificate to Qwilt CDN Site, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

// Configure adds the provider configured client to the resource.
func (r *siteActivationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cdnclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *cdnclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = cdnclient.NewSiteFacadeClient(api.SITES_HOSTNAME, client)
}

func (r *siteActivationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	var site_id, publish_id string
	//format:
	// option1 (explicit): site_id:publish_id
	// option2 (implicit): site_id (publish_id is defaulted to last active or last published)

	if len(idParts) == 2 && idParts[0] != "" && idParts[1] != "" {
		//option1: user wants to import a specific publish_id
		site_id = idParts[0]
		publish_id = idParts[1]
	} else if len(idParts) == 1 && idParts[0] != "" {
		//option2: import latest or active publish_id
		site_id = idParts[0]

		//try to get latest or active revision
		siteResp, err := r.client.GetSite(idParts[0], r.target, true, false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Getting latest/active revision for Qwilt CDN Site",
				"Could not get active/latest revision for Qwilt CDN Site, unexpected error: "+err.Error(),
			)
			return
		}
		if siteResp.ActiveAndLastPublishingOperation.Active != nil && siteResp.ActiveAndLastPublishingOperation.Active.OperationType != api.OPERATION_TYPE_UNPUBLISH {
			publish_id = siteResp.ActiveAndLastPublishingOperation.Active.PublishId
		} else if siteResp.ActiveAndLastPublishingOperation.Last != nil && siteResp.ActiveAndLastPublishingOperation.Last.OperationType != api.OPERATION_TYPE_UNPUBLISH {
			publish_id = siteResp.ActiveAndLastPublishingOperation.Last.PublishId
		}
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: site_id:publish_id OR site_id. Got: %q", req.ID),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Import: %s:%s", site_id, publish_id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("site_id"), site_id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("publish_id"), publish_id)...)
}
