package thinq

import (
	"context"
	"net/url"
)

type GatewayService service

type GatewayListRequest struct {
	Root GatewayListRequestData `json:"lgedmRoot"`
}

type GatewayListRequestData struct {
	CountryCode  string `json:"countryCode"`
	LanguageCode string `json:"langCode"`
}

type GatewayListResponse struct {
	Root GatewayListResponseData `json:"lgedmRoot"`
}

type GatewayListResponseData struct {
	ThinqURI string `json:"thinqUri"`
	EmpURI   string `json:"empUri"`
	OAuthURI string `json:"oauthUri"`
}

func (s *GatewayService) Discover(ctx context.Context) error {
	reqData := GatewayListRequest{
		Root: GatewayListRequestData{
			CountryCode:  s.client.config.CountryCode,
			LanguageCode: s.client.config.LanguageCode,
		},
	}

	req, err := s.client.NewRequest("POST", s.client.GatewayURL, "", reqData)
	if err != nil {
		return err
	}

	gResp := new(GatewayListResponse)

	_, err = s.client.Do(ctx, req, gResp)
	if err != nil {
		return err
	}

	authBase, err := url.Parse(gResp.Root.EmpURI)
	if err != nil {
		return err
	}

	apiRoot, err := url.Parse(gResp.Root.ThinqURI)
	if err != nil {
		return err
	}

	oAuthRoot, err := url.Parse(gResp.Root.OAuthURI)
	if err != nil {
		return err
	}

	s.client.AuthBase = authBase
	s.client.APIRoot = apiRoot
	s.client.OAuthRoot = oAuthRoot

	return nil
}
