package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// certificate maps Certificate schema data.
type Certificate struct {
	Id               types.Int64  `tfsdk:"id"`
	CertId           types.Int64  `tfsdk:"cert_id"`
	Certificate      types.String `tfsdk:"certificate"`
	CertificateChain types.String `tfsdk:"certificate_chain"`
	PrivateKey       types.String `tfsdk:"private_key"`
	//Email            types.String `tfsdk:"email"`
	Description types.String `tfsdk:"description"`
	PkHash      types.String `tfsdk:"pk_hash"`
	Tenant      types.String `tfsdk:"tenant"`
	Domain      types.String `tfsdk:"domain"`
	Status      types.String `tfsdk:"status"`
	Type        types.String `tfsdk:"type"`
}

type CertificateBuilder struct {
	cert Certificate
}

func (b CertificateBuilder) CertificateId(value int64) CertificateBuilder {
	b.cert.CertId = types.Int64Value(value)
	b.cert.Id = b.cert.CertId
	return b
}
func (b CertificateBuilder) Certificate(value string) CertificateBuilder {
	b.cert.Certificate = types.StringValue(value)
	return b
}
func (b CertificateBuilder) CertificateChain(value string) CertificateBuilder {
	b.cert.CertificateChain = types.StringValue(value)
	return b
}

func (b CertificateBuilder) PrivateKey(value string) CertificateBuilder {
	b.cert.PrivateKey = types.StringValue(value)
	return b
}
func (b CertificateBuilder) Description(description string) CertificateBuilder {
	b.cert.Description = types.StringValue(description)
	return b
}
func (b CertificateBuilder) PkHash(value string) CertificateBuilder {
	b.cert.PkHash = types.StringValue(value)
	return b
}
func (b CertificateBuilder) Tenant(description string) CertificateBuilder {
	b.cert.Tenant = types.StringValue(description)
	return b
}
func (b CertificateBuilder) Domain(description string) CertificateBuilder {
	b.cert.Domain = types.StringValue(description)
	return b
}
func (b CertificateBuilder) Status(description string) CertificateBuilder {
	b.cert.Status = types.StringValue(description)
	return b
}
func (b CertificateBuilder) Type(description string) CertificateBuilder {
	b.cert.Type = types.StringValue(description)
	return b
}

func (b CertificateBuilder) Build() Certificate {
	return b.cert
}
