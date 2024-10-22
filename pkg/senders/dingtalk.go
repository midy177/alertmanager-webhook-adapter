package senders

import (
	"fmt"

	"alertmanager-webhook-adapter/pkg/webhook-adapter/channels/dingtalk"
	"alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeDingTalk = "dingtalk"
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeDingTalk, createDingTalkSender)
}

func createDingTalkSender(request *restful.Request) (models.Sender, error) {
	token := request.QueryParameter("token")
	if token == "" {
		return nil, fmt.Errorf("not token found for dingtalk channel")
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = "markdown"
	}

	var sender models.Sender = dingtalk.NewSender(token, msgType)
	return sender, nil
}
