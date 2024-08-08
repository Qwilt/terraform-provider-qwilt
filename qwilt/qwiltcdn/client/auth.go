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
