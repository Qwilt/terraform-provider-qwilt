package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CertificateTemplate maps CertificateTemplate schema data.
type CertificateTemplate struct {
	Id                             types.Int64    `tfsdk:"id"`
	CertificateTemplateId          types.Int64    `tfsdk:"certificate_template_id"`
	Country                        types.String   `tfsdk:"country"`
	Tenant                         types.String   `tfsdk:"tenant"`
	State                          types.String   `tfsdk:"state"`
	Locality                       types.String   `tfsdk:"locality"`
	OrganizationName               types.String   `tfsdk:"organization_name"`
	CommonName                     types.String   `tfsdk:"common_name"`
	SANs                           []types.String `tfsdk:"sans"`
	AutoManagedCertificateTemplate types.Bool     `tfsdk:"auto_managed_certificate_template"`
	LastCertificateID              types.Int64    `tfsdk:"last_certificate_id"`
	CsrIds                         []types.Int64  `tfsdk:"csr_ids"`
}

type CertificateTemplateBuilder struct {
	cert CertificateTemplate
}

func NewCertificateTemplateBuilder() *CertificateTemplateBuilder {
	b := CertificateTemplateBuilder{}
	return &b
}

func (b *CertificateTemplateBuilder) CertificateTemplateId(value int64) *CertificateTemplateBuilder {
	b.cert.CertificateTemplateId = types.Int64Value(value)
	b.cert.Id = b.cert.CertificateTemplateId
	return b
}
func (b *CertificateTemplateBuilder) Tenant(description string) *CertificateTemplateBuilder {
	b.cert.Tenant = types.StringValue(description)
	return b
}

func (b *CertificateTemplateBuilder) Country(value string) *CertificateTemplateBuilder {
	b.cert.Country = types.StringValue(value)
	return b
}

func (b *CertificateTemplateBuilder) State(value string) *CertificateTemplateBuilder {
	b.cert.State = types.StringValue(value)
	return b
}

func (b *CertificateTemplateBuilder) Locality(value string) *CertificateTemplateBuilder {
	b.cert.Locality = types.StringValue(value)
	return b
}

func (b *CertificateTemplateBuilder) OrganizationName(value string) *CertificateTemplateBuilder {
	b.cert.OrganizationName = types.StringValue(value)
	return b
}

func (b *CertificateTemplateBuilder) CommonName(value string) *CertificateTemplateBuilder {
	b.cert.CommonName = types.StringValue(value)
	return b
}

func (b *CertificateTemplateBuilder) AutoManagedCertificateTemplate(value bool) *CertificateTemplateBuilder {
	b.cert.AutoManagedCertificateTemplate = types.BoolValue(value)
	return b
}

func (b *CertificateTemplateBuilder) LastCertificateID(value int64) *CertificateTemplateBuilder {
	b.cert.LastCertificateID = types.Int64Value(value)
	return b
}

func (b *CertificateTemplateBuilder) AddSANs(sans ...string) *CertificateTemplateBuilder {
	for _, san := range sans {
		b.cert.SANs = append(b.cert.SANs, types.StringValue(san))
	}
	return b
}

func (b *CertificateTemplateBuilder) AddCsrIds(csrIds ...int64) *CertificateTemplateBuilder {
	for _, id := range csrIds {
		b.cert.CsrIds = append(b.cert.CsrIds, types.Int64Value(id))
	}
	return b
}

func (b *CertificateTemplateBuilder) Build() CertificateTemplate {
	return b.cert
}
