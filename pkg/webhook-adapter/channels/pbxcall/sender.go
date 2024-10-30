package pbxcall

import (
	"alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Sender struct {
	Address        string   `json:"-"`
	Scheme         string   `json:"-"`
	Token          string   `json:"-"`
	Number         string   `json:"number"`
	DialPermission string   `json:"dial_permission"`
	Prompts        []string `json:"prompts"`
}

func NewSender(address, scheme, token, number, dialPermission string, prompts []string) models.Sender {
	if address == "" {
		address = "pbxcall.yeastardigital.com:8181"
	}
	if scheme != "http" && scheme != "https" {
		scheme = "http"
	}
	return &Sender{
		Address:        address,
		Scheme:         scheme,
		Token:          token,
		Number:         number,
		DialPermission: dialPermission,
		Prompts:        prompts,
	}
}

func (s *Sender) Send(payload *models.Payload) error {
	return s.sendPbxCallRequest()
}

func (s *Sender) SendMsg(msgSource interface{}) error {
	return nil
}

func (s *Sender) SendMsgT(msgType string, msgSource interface{}) error {
	return nil
}

func (s *Sender) sendPbxCallRequest() error {

	postUrl := fmt.Sprintf("%s://%s/play_prompt", s.Scheme, s.Address)

	// 将数据编码为 JSON
	jsonData, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("Error marshalling JSON: %v", err)
	}
	// 创建新的请求
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.Token)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// 处理响应
	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		return fmt.Errorf("Request failed with status: %s", resp.Status)
	}
}
