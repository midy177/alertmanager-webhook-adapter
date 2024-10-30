package senders

import (
	"alertmanager-webhook-adapter/pkg/webhook-adapter/channels/pbxcall"
	"alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"fmt"
	"github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypePbxCall = "pbx_call"
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypePbxCall, createPbxCallSender)
}

func createPbxCallSender(request *restful.Request) (models.Sender, error) {
	number := request.QueryParameter("number")
	if number == "" {
		return nil, fmt.Errorf("not number found for pbx_play_prompt channel")
	}

	prompts := request.QueryParameters("prompts")
	if len(prompts) == 0 {
		return nil, fmt.Errorf("not prompts list found for pbx_play_prompt channel")
	}

	address := request.QueryParameter("address")

	dialPermission := request.QueryParameter("dial_permission")

	scheme := request.QueryParameter("scheme")

	token := request.QueryParameter("token")

	var sender = pbxcall.NewSender(address, scheme, token, number, dialPermission, prompts)
	return sender, nil
}
