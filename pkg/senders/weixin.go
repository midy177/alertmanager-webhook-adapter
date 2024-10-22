package senders

import (
	"fmt"

	"alertmanager-webhook-adapter/pkg/webhook-adapter/channels/weixin"
	"alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeWeixin = "weixin"
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeWeixin, createWeixinSender)
}

func createWeixinSender(request *restful.Request) (models.Sender, error) {
	token := request.QueryParameter("token")
	if token == "" {
		return nil, fmt.Errorf("not token found for weixin channel")
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = "markdown"
	}

	var sender = weixin.NewSender(token, msgType)
	return sender, nil
}
