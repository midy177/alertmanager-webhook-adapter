package senders

import (
	"alertmanager-webhook-adapter/pkg/webhook-adapter/channels/pbxplayprompt"
	"alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"strconv"
)

const (
	ChannelTypePbxPlayPrompt = "pbx_play_prompt"
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypePbxPlayPrompt, createPbxPlayPromptSender)
}

func createPbxPlayPromptSender(request *restful.Request) (models.Sender, error) {
	username := request.QueryParameter("username")
	if username == "" {
		return nil, fmt.Errorf("not username found for pbx_play_prompt channel")
	}

	password := request.QueryParameter("password")
	if password == "" {
		return nil, fmt.Errorf("not password found for pbx_play_prompt channel")
	}

	number := request.QueryParameter("number")
	if number == "" {
		return nil, fmt.Errorf("not number found for pbx_play_prompt channel")
	}

	prompts := request.QueryParameters("prompts")
	if len(prompts) == 0 {
		return nil, fmt.Errorf("not prompts list found for pbx_play_prompt channel")
	}

	countStr := request.QueryParameter("count")
	count, _ := strconv.Atoi(countStr)

	address := request.QueryParameter("address")

	dialPermission := request.QueryParameter("dial_permission")

	autoAnswer := request.QueryParameter("auto_answer")

	volumeStr := request.QueryParameter("volume")
	volume, _ := strconv.Atoi(volumeStr)

	var sender = pbxplayprompt.NewSender(address, username, password, number, prompts, count, dialPermission, autoAnswer, volume)
	return sender, nil
}
