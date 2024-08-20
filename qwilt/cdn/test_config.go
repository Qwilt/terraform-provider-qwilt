// Package qwiltcdn
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.
package cdn

var test_domain_crt = `-----BEGIN CERTIFICATE-----
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

var test_domain_key = `-----BEGIN ENCRYPTED PRIVATE KEY-----
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

const (
	// QwiltCdnProviderConfig is a shared configuration to combine with the actual
	// test configuration so the CDN client is properly configured.
	// It is also possible to use the QCDN_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	QwiltCdnProviderConfig = `
provider "qwilt" {
	env_type = "dev"
}
`

	// QwiltCdnProviderConfig is a shared configuration to combine with the actual
	// test configuration so the CDN client is properly configured.
	// It is also possible to use the QCDN_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	QwiltCdnFullProviderConfig = `
terraform {
  required_providers {
    qwilt = {
      source = "qwilt.com/qwiltinc/qwilt"
    }
  }
}

provider "qwilt" {
	env_type = "dev"
}
`
)
