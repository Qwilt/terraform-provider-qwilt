// Package cdn
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.
package cdn

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
