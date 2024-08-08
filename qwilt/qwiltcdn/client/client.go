package client

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client -
type Client struct {
	envType         string
	HTTPClient      *http.Client
	Token           string
	XApiToken       string
	Auth            AuthStruct
	authEndpoint    string
	endpointBuilder EndpointBuilder
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	Username string
	Token    string
}

// NewClient -
func NewClient(envType,
	username,
	password,
	xApiToken string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 40 * time.Second},
		envType:    "",
		Auth: AuthStruct{
			Username: username,
			Password: password,
		},
	}

	c.envType = envType
	c.XApiToken = xApiToken

	c.endpointBuilder = NewEndpointBuilder(c.envType)
	c.authEndpoint = c.endpointBuilder.Build("login")

	// Attempt to sign in using username/password if no token specified
	// and auth endpoint has been provided.
	// TODO: Can we defer this to be lazily sign in?
	if c.authEndpoint != "" && xApiToken == "" {
		ar, err := c.SignIn()
		if err != nil {
			return nil, err
		}
		c.Token = ar.Token
	}

	return &c, nil
}

// doAuthRequest - Authenticates against the login API to get a token
func (c *Client) doAuthRequest(req *http.Request) (*string, error) {
	auth := c.Auth.Username + ":" + c.Auth.Password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", basicAuth)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusFound {
		return nil, fmt.Errorf("Authentication failed - status: %d", res.StatusCode)
	}

	// Extract token from cookie string
	var token string
	cookies := res.Cookies()
	if len(cookies) > 0 {
		for _, cookie := range cookies {
			if cookie.Name == "cqloudLoginToken" {
				token = cookie.Value
				break
			}
		}
		if len(token) == 0 {
			return nil, fmt.Errorf("XApiToken cookie has an empty value.")
		}
	} else {
		return nil, fmt.Errorf("No Set-Cookie header found in the response.")
	}

	return &token, err
}

// doRequest - Performs a typical request against the API
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	if c.XApiToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("X-API-KEY %s", c.XApiToken))
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		if res.StatusCode == 401 {
			err = fmt.Errorf("401 Unauthorized. please re-athenticate - status: %d", res.StatusCode)
		} else {
			err = fmt.Errorf("API command failed - status: %d, body: %s", res.StatusCode, body)
		}
	}

	return body, err
}
