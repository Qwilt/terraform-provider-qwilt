package qwiltcdn

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strings"
	"testing"
)

var domain_crt = `-----BEGIN CERTIFICATE-----
MIIDwzCCAqugAwIBAgIUCyRlvJ8W/4OLc39I04FH6KU0tAUwDQYJKoZIhvcNAQEL
BQAwgYkxCzAJBgNVBAYTAmlsMQ8wDQYDVQQIDAZtZXJrYXoxETAPBgNVBAcMCHRl
bC1hdml2MRAwDgYDVQQKDAdwcml2YXRlMQwwCgYDVQQLDANybmQxGDAWBgNVBAMM
D3d3dy5leGFtcGxlLmNvbTEcMBoGCSqGSIb3DQEJARYNc3NzQGdtYWlsLmNvbTAe
Fw0yNDA2MTkwNjE4MTBaFw0yNTA2MTkwNjE4MTBaMIGJMQswCQYDVQQGEwJpbDEP
MA0GA1UECAwGbWVya2F6MREwDwYDVQQHDAh0ZWwtYXZpdjEQMA4GA1UECgwHcHJp
dmF0ZTEMMAoGA1UECwwDcm5kMRgwFgYDVQQDDA93d3cuZXhhbXBsZS5jb20xHDAa
BgkqhkiG9w0BCQEWDXNzc0BnbWFpbC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQDIhN5oZ/0ESV2ZOauU0RAVD74siKWWEadgc6pdy/Fv+um3qaiP
Wz7ju7cziSrw4Ynt2MYrHnfNAgeFp5NKZy8NzzJD1gt2KperH+xgljW1//4fkSzQ
0qQEoMeZyF1Cfy1GNc/YV6LGiddEjjtI57ZLYHz0KYNIWNyWf5QaPmV5kYKsjg4S
0VizHHG97ZVmHzBGxmshogemkfdUUjiebdDKbsCkiw9jFA+s0XRpO3IPD/LhBMTB
z0rquaMAIBA2mhg9XY23zba6FwRMQFfsr9du8GIACrxyyE8e9IKtTPLjSuFx/xiC
0ArHFewsGpioyhgJpla+fUxNV5WdGUUDK/JFAgMBAAGjITAfMB0GA1UdDgQWBBSn
RkzAstcJG/aqEVPsOvmBFlURLTANBgkqhkiG9w0BAQsFAAOCAQEABC5tLVSSBxTH
iiW8MlinQb5tz+QneSFkwKyajysc3AziXgnMWfNg3quxU0L9WOZ3Ij2MyGLSmvm4
Z072dl8lc7cSLq5FAqBdiKTx8lp0LqqFUgZfjRRz1RGxvrJYbYw0sfnnQc6Yx4Rn
vOr6IVXabcLAt9uZNkd9kaSzUngrToeN85JOjE/AXaC/Tq2vizXKAmcsAxeg/95U
+UmdiFqmfCPrGKVzkWODoSv0uU31KCLBjiQf9VDKTP4IQtwBgakMSYeFRGCxR6kG
FyFHV4d/VOCFecm3XFBqNhlP5WxqK3P6twUH/L1395ny1iDL2laPKdl5qXhYBioq
xkwL9SJBxg==
-----END CERTIFICATE-----`

var domain_key = `-----BEGIN ENCRYPTED PRIVATE KEY-----
MIIFJDBWBgkqhkiG9w0BBQ0wSTAxBgkqhkiG9w0BBQwwJAQQEaFzlPK16IzrWscJ
LJEa5wICCAAwDAYIKoZIhvcNAgkFADAUBggqhkiG9w0DBwQIZ4sNc988oJ0EggTI
jrNQF89NwkYwIfYQv4LfZpxwJ+O+KD35d+Zh6RD9fdVxFRmE0WcYaJPr2nGpds2m
RToW+4mNlgcVF+E39Np/vPhwlTFvTpdaTGtTCo+NiU81MfLCgBCcbE0Nc3O5QG98
QtVTxAHV4Y8clbkiHyeTA4PiOp1AiNxIBbbJqGmJqw8oIv61Lg8YTd59lWnYOP9x
rAS03vg6VBoDXAfEBXU6r/J1xuy/gFD6hLLw9k2E797FPmGNL7gFDqpS53X3Cp7X
+S7zUBwwLLX7MQfzjNViX+O4VAdUwh7F48J5EnebY9YFxBSOXK0HC1T5TCRyRMlX
C3Et8ssthx6qU/tIOkqYzz+Olp6A1FdqIXg+NktShrul5L98Zte8BGYASsPIEG/n
XKBI47I3OjlJAY9Nkr0Pv6Itur6Zisj6IleVhRJYwTqW4bkhzHoeP0rt0x9XcYoD
ycRvfjI8AielFS5QWfFRPYpxwhzkDrg9xBO2hkmzI1yLvCWcQWnzKYE3qirqGeXV
jjbxUq1NjF/7r0mNFM0GCs0cwIsfIs0Wq3wII6723oUyY2Ozv8wOWKX0mfCSLVCP
H1Posc0pbCWVoGPT8uMLUdBlkSJQHlJHlpxb4SDq0kTbwuUihtSY8Ef5A/hOLh64
mYkoN6rbKAr70hqwMEHw3OpMNb6hw7e171aR7l6gX+kyGULgHm6Z+3Ag2MEfgGbO
6UNxASuBZoNgj/uk99U8dUb4rjLPMm5G2jq8wT53MNn7VF8585dr6era343Ix544
RfCsVmppx2bs+Nj6cRl8y/OaU3F3b1Ipcte89zIas5X1m36ao27NIaXj8E5rhlua
1vlSU3UKW8RxQ+jX5sCR1QAlIZt0vBahRn7CyBaUbpGzzSahfOLgNbOQRVj2xPN7
puXkM3pFq/sroX2JbaF912Bog0xoY9hEjVxlGF/WfIyogmoARCVXHQmpVvuTnI3/
kr0SRpe50bkmV4BzRS/YtR+rHm3/G1OLMI3PuYYPe532xszcEceftRRz/FxFKCMQ
HWhwYl1WrTUu+l2Qhunito4twUJQmm4HMFFV5hY6KbouxMeBnXMRLM/t6COpbF+O
KAC24vR7HG+Ia6dbmS82Ofdp2k3EDo/BCW68a+vy6On0LV/IOeTI7n6aNqcPA43D
YWNqQXV3vwovUhqdGvQGa3pnKQNSBvS9BGQt8N0hL2oeEzpOREQWjFAWUE6L1ne0
jDGX7RlGb4hKvraHG01tWy9lqBPzKyyAG8xnXrVH8SICSv1y9/ZUanNsDa9X4q6H
KA2pu8F6cBZuxgDMoHP7oPKbjWYtlaP6I2BPGrW1u3zL15t2s2isgUYSZ1Lc9f/F
3NK8wiXhf6gU4/Ei0Zk10PEBOtMvOzsoz1drzKtGFbAN1DVyB2wrZ0uOMDeo8dEz
OjU15s+t+eGBKgFZHyeRnTu79p7Yiv65ab3x6dnX7bkAezD3oyxL/xnk22UNpfPz
duG6PXZEa8ercuTSuI+FZVMZXXnQ+x/4bUVQqtAGhiv9uMJAmPUWrMWeoS4ULK4g
N3SEZmdjk6Flikt1vd/mN6qUro4wK3tlPcME69a+eMCWUlPXD72/4l1z6GNPo4W5
ELIDymzaDzbECPBOXQwRDK56js31sQYL
-----END ENCRYPTED PRIVATE KEY-----`

func TestCertificateResource(t *testing.T) {

	t.Logf("Starting TestSiteResource test")

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

	//tf.SetStdout(os.Stdout)
	//tf.SetStderr(os.Stderr)

	terraformBuilder := NewTerraformConfigBuilder().CertResource("test", domain_key, domain_crt, "ccc")
	terraformConfig := terraformBuilder.Build()

	//t.Logf("config: %s", terraformConfig)
	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err := tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(state.Values.RootModule.Resources))

	certState := findStateResource(state, "qwiltcdn_certificate", "test")
	assert.NotNil(t, certState)

	certId := certState.AttributeValues["cert_id"]
	t.Logf("cert_id: %s", certId)
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(domain_crt)), strings.TrimSpace(certState.AttributeValues["certificate"].(string)))
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(domain_crt)), strings.TrimSpace(certState.AttributeValues["certificate_chain"].(string)))
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(domain_key)), strings.TrimSpace(certState.AttributeValues["private_key"].(string)))
	assert.Equal(t, "ccc", certState.AttributeValues["description"])
	assert.Equal(t, "www.example.com", certState.AttributeValues["domain"])
	assert.Equal(t, "devorg", certState.AttributeValues["tenant"])

	//check that plan gives no diff - this actually checks the refresh and that all attributes in the state are the same as in the configuration
	plan, err := tf.Plan(context.Background())
	assert.Equal(t, nil, err)
	assert.False(t, plan) //no diff

	//update just the description
	terraformBuilder.CertResource("test", domain_key, domain_crt, "description-after-change")
	terraformConfig = terraformBuilder.Build()

	err = os.WriteFile(tfFilePath, []byte(terraformConfig), 0644)
	assert.Equal(t, nil, err)

	err = tf.Apply(context.Background())
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(state.Values.RootModule.Resources))

	certState = findStateResource(state, "qwiltcdn_certificate", "test")
	assert.NotNil(t, certState)

	assert.Equal(t, certId, certState.AttributeValues["cert_id"]) //make sure therewas no replace, just update
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(domain_crt)), strings.TrimSpace(certState.AttributeValues["certificate"].(string)))
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(domain_crt)), strings.TrimSpace(certState.AttributeValues["certificate_chain"].(string)))
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(domain_key)), strings.TrimSpace(certState.AttributeValues["private_key"].(string)))
	assert.Equal(t, "description-after-change", certState.AttributeValues["description"])
	assert.Equal(t, "www.example.com", certState.AttributeValues["domain"])
	assert.Equal(t, "devorg", certState.AttributeValues["tenant"])

	//check that plan gives no diff - this actually checks the refresh and that all attributes in the state are the same as in the configuration
	plan, err = tf.Plan(context.Background())
	assert.Equal(t, nil, err)
	assert.False(t, plan) //no diff

	//prepare for import
	err = tf.StateRm(context.Background(), "qwiltcdn_certificate.test")
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Nil(t, state.Values)

	err = tf.Import(context.Background(), "qwiltcdn_certificate.test", fmt.Sprintf("%s", certId))
	assert.Equal(t, nil, err)

	state, err = tf.Show(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(state.Values.RootModule.Resources))

	certState = findStateResource(state, "qwiltcdn_certificate", "test")
	assert.NotNil(t, certState)

	assert.Equal(t, certId, certState.AttributeValues["cert_id"])
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(domain_crt)), strings.TrimSpace(certState.AttributeValues["certificate"].(string)))
	assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(domain_crt)), strings.TrimSpace(certState.AttributeValues["certificate_chain"].(string)))
	//assert.Equal(t, b64.URLEncoding.EncodeToString([]byte(domain_key)), strings.TrimSpace(certState.AttributeValues["private_key"].(string)))
	assert.Equal(t, "description-after-change", certState.AttributeValues["description"])
	assert.Equal(t, "www.example.com", certState.AttributeValues["domain"])
	assert.Equal(t, "devorg", certState.AttributeValues["tenant"])

	//check that plan gives no diff - this actually checks the refresh and that all attributes in the state are the same as in the configuration
	plan, err = tf.Plan(context.Background())
	assert.Equal(t, nil, err)
	assert.False(t, plan) //no diff

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
