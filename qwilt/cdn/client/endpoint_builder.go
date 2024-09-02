// Package client
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.

package client

import "fmt"

type EndpointBuilder struct {
	domain string
	prefix string
}

func NewEndpointBuilder(envType string) EndpointBuilder {
	b := EndpointBuilder{}

	if envType == "stage" {
		b.domain = "stage.cqloud.com"
		b.prefix = "stage-"
	} else if envType == "prestg" {
		b.domain = "prestg.cqloud.com"
		b.prefix = "prestg-"
	} else if envType == "dev" {
		b.domain = "rnd.cqloud.com"
		b.prefix = "kan11-" //change this and re-build to run locally on different env on development time
	} else if envType == "prod" {
		b.domain = "cqloud.com"
		b.prefix = ""
	}
	return b
}

func (e EndpointBuilder) Build(name string) string {
	return fmt.Sprintf("https://%s%s.%s", e.prefix, name, e.domain)
}
