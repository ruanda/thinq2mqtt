package thinq

import (
	"context"
)

type SessionService service

type SessionStartRequest struct {
	Root SessionStartRequestData `json:"lgedmRoot"`
}

type SessionStartRequestData struct {
	CountryCode  string `json:"countryCode"`
	LanguageCode string `json:"langCode"`
	LoginType    string `json:"loginType"`
	AccessToken  string `json:"token"`
}

type SessionStartResponse struct {
	Root SessionStartResponseData `json:"lgedmRoot"`
}

type SessionStartResponseData struct {
	ResponseCode    string       `json:"returnCd"`
	ResponseMessage string       `json:"returnMsg"`
	SessionID       string       `json:"jsessionId"`
	Devices         []DeviceData `json:"item"`
}

type DeviceData struct {
	Model        string `json:"modelNm"`
	Type         uint   `json:"deviceType"`
	Code         string `json:"deviceCode"`
	Alias        string `json:"alias"`
	ID           string `json:"deviceId"`
	Firmware     string `json:"fwVer"`
	ModelInfoURL string `json:"modelJsonUrl"`
	MAC          string `json:"macAddress"`
}

func (s *SessionService) Start(ctx context.Context) error {
	reqData := SessionStartRequest{
		Root: SessionStartRequestData{
			CountryCode:  s.client.config.CountryCode,
			LanguageCode: s.client.config.LanguageCode,
			LoginType:    "EMP",
			AccessToken:  s.client.AccessToken,
		},
	}

	req, err := s.client.NewJSONRequest("POST", s.client.APIRoot, "member/login", reqData)
	if err != nil {
		return err
	}

	sResp := new(SessionStartResponse)

	_, err = s.client.Do(ctx, req, sResp)
	if err != nil {
		return err
	}

	s.client.SessionID = sResp.Root.SessionID
	for _, d := range sResp.Root.Devices {
		s.client.Devices.Add(&d)
	}
	return nil
}
