package thinq

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type AuthService service

type CallbackResult struct {
	RefreshToken string
	AccessToekn  string
}

type AccessTokenResponse struct {
	Status      uint   `json:"status"`
	ErrorCode   string `json:"lgoauth_error_code,omitempty"`
	Message     string `json:"message,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   string `json:"expires_in,omitempty"`
}

func oAuthSignature(message, secret string) string {
	messageBytes := []byte(message)
	secretBytes := []byte(secret)
	h := hmac.New(sha1.New, secretBytes)
	h.Write(messageBytes)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func oAuthTimestamp() string {
	return time.Now().UTC().Format("Mon, 2 Jan 2006 15:04:05 +0000")
}

func (s *AuthService) GetOAuthURL() (*url.URL, error) {
	u, err := s.client.AuthBase.Parse("login/sign_in")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("country", s.client.config.CountryCode)
	q.Add("language", s.client.config.LanguageCode)
	q.Add("svcCode", s.client.config.ServiceCode)
	q.Add("client_id", s.client.config.ClientID)
	q.Add("authSvr", "oauth2")
	q.Add("division", "ha")
	q.Add("grant_type", "password")

	u.RawQuery = q.Encode()

	return u, nil
}

func (s *AuthService) ParseOAuthCallback(callbackURL string) (*CallbackResult, error) {
	u, err := url.Parse(callbackURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	res := &CallbackResult{
		RefreshToken: q.Get("refresh_token"),
		AccessToekn:  q.Get("access_token"),
	}
	return res, nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context) error {
	if s.client.RefreshToken == "" {
		return errors.New("auth: client has no refresh token set")
	}

	reqData := url.Values{}
	reqData.Set("grant_type", "refresh_token")
	reqData.Set("refresh_token", s.client.RefreshToken)

	signURL := fmt.Sprintf("/oauth2/token?grant_type=refresh_token&refresh_token=%s", s.client.RefreshToken)
	timestamp := oAuthTimestamp()
	signature := oAuthSignature(fmt.Sprintf("%s\n%s", signURL, timestamp), s.client.config.ClientSecret)

	req, err := s.client.NewRequest("POST", s.client.OAuthRoot, "oauth2/token", strings.NewReader(reqData.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.Header.Set(headerClientID, s.client.config.ClientID)
	req.Header.Set(headerSignature, signature)
	req.Header.Set(headerDate, timestamp)

	resp := new(AccessTokenResponse)

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return err
	}
	if resp.Status != 1 {
		return errors.New(resp.Message)
	}

	expiresIn, err := strconv.Atoi(resp.ExpiresIn)
	if err != nil {
		return err
	}

	s.client.AccessToken = resp.AccessToken
	s.client.AccessTokenExpirationDate = time.Now().Add(time.Second * time.Duration(expiresIn))
	return nil
}
