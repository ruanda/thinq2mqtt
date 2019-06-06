package thinq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultGatewayURL = "https://kic.lgthinq.com:46030/api/common/gatewayUriList"

	headerApplicationKey = "x-thinq-application-key"
	headerSecurityKey    = "x-thinq-security-key"

	defaultApplicationKey = "thinq"
	defaultSecurityKey    = "thinq"
)

type Config struct {
	CountryCode  string
	LanguageCode string
	ServiceCode  string
	ClientID     string
}

type Client struct {
	config Config
	client *http.Client

	ApplicationKey string
	SecurityKey    string

	GatewayURL *url.URL
	AuthBase   *url.URL
	APIRoot    *url.URL
	OAuthRoot  *url.URL

	common service

	Gateway *GatewayService
	Auth    *AuthService
}

type service struct {
	client *Client
}

func NewClient(config Config, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}
	gatewayURL, _ := url.Parse(defaultGatewayURL)

	c := &Client{
		config:         config,
		client:         httpClient,
		ApplicationKey: defaultApplicationKey,
		SecurityKey:    defaultSecurityKey,
		GatewayURL:     gatewayURL,
	}
	c.common.client = c
	c.Gateway = (*GatewayService)(&c.common)
	c.Auth = (*AuthService)(&c.common)
	return c, nil
}

func (c *Client) NewRequest(method string, baseURL *url.URL, url string, body interface{}) (*http.Request, error) {

	u, err := baseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set(headerApplicationKey, c.ApplicationKey)
	req.Header.Set(headerSecurityKey, c.SecurityKey)
	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return resp, err
}

type ErrorResponse struct {
	Response *http.Response
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d",
		r.Response.Request.Method,
		r.Response.Request.URL,
		r.Response.StatusCode)
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	return errorResponse
}
