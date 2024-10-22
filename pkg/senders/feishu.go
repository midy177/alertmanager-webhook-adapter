package senders

import (
	"fmt"

	"alertmanager-webhook-adapter/pkg/webhook-adapter/channels/feishu"
	"alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeFeishu = "feishu"
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeFeishu, createFeiShuSender)
}

func createFeiShuSender(request *restful.Request) (models.Sender, error) {
	token := request.QueryParameter("token")
	if token == "" {
		return nil, fmt.Errorf("not token found for feishu channel")
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = "markdown"
	}

	var sender = feishu.NewSender(token, msgType)
	return sender, nil
}
