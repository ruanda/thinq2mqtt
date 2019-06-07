package thinq

import "net/url"

type AuthService service

type CallbackResult struct {
	RefreshToken string
	AccessToekn  string
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
