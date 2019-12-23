package bone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var (
	_ Client         = (*DefaultClient)(nil)
	_ ServiceManager = (*DefaultClient)(nil)
)

// Client is API client
type Client interface {
	SetClient(c *http.Client)
	ServiceManager
}

// ServiceManager manages API services
type ServiceManager interface {
	GetService(name string) (Service, error)
	SetService(name string, s Service)
}

// Service is interface for API service
type Service interface {
	SetClient(c Client)
}

// DefaultClient is the default Client and is used API request
type DefaultClient struct {
	// HTTP client used to communicate with the DO API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// Services manages API services
	Services map[string]Service

	// header name for access token
	TokenKey string

	// API Access Token
	Token string

	CheckResponse ResponseChecker
}

// ClientOpt are options for New.
type ClientOpt func(Client) error

// NewClient returns a new cunstom API client.
func NewClient(v Client, httpClient *http.Client, opts ...ClientOpt) (err error) {
	v.SetClient(httpClient)
	for _, opt := range opts {
		if err := opt(v); err != nil {
			return err
		}
	}
	return nil
}

// SetClient sets http client
func (c *DefaultClient) SetClient(httpClient *http.Client) {
	if httpClient != nil {
		c.client = httpClient
	} else if c.client == nil {
		c.client = http.DefaultClient
	}
	c.SetResponseChecker(CheckResponse)
}

// GetService gets a API service
func (c *DefaultClient) GetService(name string) (s Service, err error) {
	s, ok := c.Services[name]
	if !ok {
		return nil, fmt.Errorf("%s does not exist", name)
	}
	return s, nil
}

// SetService sets a API service
func (c *DefaultClient) SetService(name string, s Service) {
	if c.Services == nil {
		c.Services = map[string]Service{}
	}
	c.Services[name] = s
}

// SetUserAgent sets the user agent.
func (c *DefaultClient) SetUserAgent(ua string) {
	if ua != "" {
		c.UserAgent = ua
	}
}

// SetBaseURL sets the base URL.
func (c *DefaultClient) SetBaseURL(bu string) (err error) {
	u, err := url.Parse(bu)
	if err != nil {
		return err
	}
	c.BaseURL = u
	return nil
}

// SetResponseChecker sets CheckResponse
func (c *DefaultClient) SetResponseChecker(f ResponseChecker) {
	c.CheckResponse = CheckResponse
}

func (c *DefaultClient) NewRequest(method, urlStr string, headers map[string]string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)
	if c.TokenKey != "" && c.Token != "" {
		req.Header.Set(c.TokenKey, c.Token)
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	return req, nil
}

// Do requests API
func (c *DefaultClient) Do(ctx context.Context, req *http.Request, decode ResponseDecode, v interface{}) (response *http.Response, err error) {
	return Do(ctx, c.client, req, decode, c.CheckResponse, v)
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(client *http.Client, req *http.Request) (*http.Response, error) {
	return client.Do(req)
}

// Do requests API
func Do(ctx context.Context, c *http.Client, req *http.Request, decode ResponseDecode, rc ResponseChecker, v interface{}) (response *http.Response, err error) {
	var resp *http.Response
	if resp, err = DoRequestWithClient(c, req); err != nil {
		return nil, err
	}
	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if err = rc(resp); err != nil {
		return resp, err
	}

	if v != nil {
		if err = decode(resp.Body, v); err != nil {
			return nil, err
		}
	}
	return resp, nil
}
