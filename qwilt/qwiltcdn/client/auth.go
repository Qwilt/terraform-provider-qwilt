// Package client
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// SignIn - Get a new token for user
func (c *Client) SignIn() (*AuthResponse, error) {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return nil, fmt.Errorf("Please define the username and password to authenticate")
	}

	rb, err := json.Marshal(c.Auth)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/login", c.authEndpoint), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	token, err := c.doAuthRequest(req)
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{}
	ar.Token = *token
	ar.Username = c.Auth.Username

	return &ar, nil
}
