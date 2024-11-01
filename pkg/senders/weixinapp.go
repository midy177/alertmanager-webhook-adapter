package senders

import (
	"fmt"
	"strconv"

	"alertmanager-webhook-adapter/pkg/webhook-adapter/channels/weixinapp"
	"alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeWeixinApp = "weixinapp"
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeWeixinApp, createWeixinappSender)
}

func createWeixinappSender(request *restful.Request) (models.Sender, error) {
	corpID := request.QueryParameter("corp_id")
	if corpID == "" {
		return nil, fmt.Errorf("not core_id found for weixin channel")
	}

	agentID := request.QueryParameter("agent_id")
	if agentID == "" {
		return nil, fmt.Errorf("not agent_id found for weixin channel")
	}

	aID, err := strconv.Atoi(agentID)
	if err != nil {
		return nil, fmt.Errorf("agent_id must be integer")
	}

	agentSecret := request.QueryParameter("agent_secret")
	if agentSecret == "" {
		return nil, fmt.Errorf("not agent_secret found for weixin channel")
	}

	toUser := request.QueryParameter("to_user")
	toParty := request.QueryParameter("to_party")
	toTag := request.QueryParameter("to_tag")

	if toUser == "" && toParty == "" && toTag == "" {
		return nil, fmt.Errorf("must specify one of to_user,to_party,to_tag")
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = "markdown"
	}

	var sender = weixinapp.NewSender(corpID, aID, agentSecret, msgType, toUser, toParty, toTag)
	return sender, nil
}
