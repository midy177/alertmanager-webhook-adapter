package senders

import (
	"fmt"

	"alertmanager-webhook-adapter/pkg/webhook-adapter/channels/discord"
	"alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeDiscordWebhook = "discord-webhook"
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeDiscordWebhook, createDiscordWebhookSender)
}

func createDiscordWebhookSender(request *restful.Request) (models.Sender, error) {
	id := request.QueryParameter("id")
	if id == "" {
		return nil, fmt.Errorf("not id found for discord-webhook channel")
	}

	token := request.QueryParameter("token")
	if token == "" {
		return nil, fmt.Errorf("not token found for discord-webhook channel")
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = "markdown"
	}

	var sender models.Sender = discord.NewWebhookSender(id, token)
	return sender, nil
}
