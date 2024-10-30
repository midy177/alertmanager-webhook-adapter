package pbxplayprompt

import (
	"alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"
)

// make cache with 5m TTL and 5 max keys
var tokenCache = cache.New(20*time.Minute, 20*time.Minute)

type Sender struct {
	Address string
	GetTokenRequest
	PlayPromptRequest
}

type GetTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PlayPromptRequest struct {
	Number         string   `json:"number"`
	Prompts        []string `json:"prompts"`
	Count          int      `json:"count"`
	DialPermission string   `json:"dial_permission"`
	AutoAnswer     string   `json:"auto_answer"`
	Volume         int      `json:"volume"`
}

type GetTokenResponse struct {
	ErrCode                int    `json:"errcode"`
	ErrMsg                 string `json:"errmsg"`
	AccessTokenExpireTime  int64  `json:"access_token_expire_time"`
	AccessToken            string `json:"access_token"`
	RefreshTokenExpireTime int64  `json:"refresh_token_expire_time"`
	RefreshToken           string `json:"refresh_token"`
}

type PlayPromptResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	CallId  string `json:"call_id"`
}

func NewSender(address, username, password, number string, prompts []string, count int, dialPermission string, autoAnswer string, volume int) models.Sender {
	return newSender(address, username, password, number, prompts, count, dialPermission, autoAnswer, volume)
}

func newSender(address, username, password, number string, prompts []string, count int, dialPermission string, autoAnswer string, volume int) *Sender {
	// default is https://zct.test.smartpbx.cn/openapi
	var sender = Sender{
		Address: address,
		GetTokenRequest: GetTokenRequest{
			Username: username,
			Password: password,
		},
		PlayPromptRequest: PlayPromptRequest{
			Number:         number,
			Prompts:        prompts,
			Count:          count,
			DialPermission: dialPermission,
			AutoAnswer:     autoAnswer,
			Volume:         volume,
		},
	}
	if sender.Address == "" {
		sender.Address = "https://pbx.ras.yeastar.com"
	}
	if sender.Count == 0 {
		sender.Count = 1
	}
	//if sender.DialPermission == "" {
	//	sender.DialPermission = ""
	//}
	if sender.AutoAnswer != "yes" && sender.AutoAnswer != "no" {
		sender.AutoAnswer = "no"
	}
	if sender.Volume == 0 {
		sender.Volume = 20
	}

	return &sender
}

func (s *Sender) Send(payload *models.Payload) error {
	token, err := s.getToken()
	if err != nil {
		return err
	}
	return s.playPrompt(token)
}

func (s *Sender) SendMsg(msgSource interface{}) error {
	token, err := s.getToken()
	if err != nil {
		return err
	}
	return s.playPrompt(token)
}

func (s *Sender) SendMsgT(msgType string, msgSource interface{}) error {
	token, err := s.getToken()
	if err != nil {
		return err
	}
	return s.playPrompt(token)
}

func (s *Sender) getToken() (string, error) {
	if v, ok := tokenCache.Get(s.Address); ok {
		if accessToken, aOk := v.(string); aOk {
			return accessToken, nil
		}
	}
	url := fmt.Sprintf("%s/openapi/v1.0/get_token", s.Address)
	requestBody, err := json.Marshal(s.GetTokenRequest)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token: %s", resp.Status)

	}

	var tokenResponse GetTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}
	if tokenResponse.AccessTokenExpireTime != 0 {
		tokenCache.Set(s.Address, tokenResponse.AccessToken, time.Duration(tokenResponse.AccessTokenExpireTime-60)*time.Second)
	} else {
		tokenCache.SetDefault(s.Address, tokenResponse.AccessToken)
	}
	return tokenResponse.AccessToken, nil

}

func (s *Sender) playPrompt(accessToken string) error {
	url := fmt.Sprintf("%s/openapi/v1.0/call/play_prompt?access_token=%s", s.Address, accessToken)
	requestBody, err := json.Marshal(s.PlayPromptRequest)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to play prompt, status code: %s", resp.Status)
	}
	return nil

}
