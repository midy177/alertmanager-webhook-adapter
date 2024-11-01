package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	promModels "alertmanager-webhook-adapter/pkg/models"
	"alertmanager-webhook-adapter/pkg/senders"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/kr/pretty"
)

type Controller struct {
	signature string
	debug     bool
}

func NewController(signature string) *Controller {
	return &Controller{
		signature: signature,
	}
}

func (c *Controller) WithDebug(debug bool) *Controller {
	if debug {
		fmt.Println("debug mode enabled")
	}
	c.debug = debug
	return c
}

func (c *Controller) Install(container *restful.Container) {

	ws := new(restful.WebService)
	ws.Path("/webhook/send")

	ws.Route(
		ws.POST("/").To(c.send),
	)

	container.Add(ws)
}

func (c *Controller) logf(format string, a ...any) error {
	if c.debug {
		_, err := fmt.Printf(format, a...)
		return err
	}

	return nil
}

func (c *Controller) log(a ...any) error {
	if c.debug {
		_, err := fmt.Println(a...)
		return err
	}
	return nil
}

func (c *Controller) send(request *restful.Request, response *restful.Response) {
	c.logf("Got request : %s\n", request.Request.URL.String())

	raw, err := io.ReadAll(request.Request.Body)
	if err != nil {
		errMsg := fmt.Sprintf("Err: read request body failed, err: %s", err)
		_ = c.log(errMsg)
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, errMsg, restful.MIME_JSON)
		return
	}

	promMsg := &promModels.AlertmanagerWebhookMessage{}
	if err := json.Unmarshal(raw, promMsg); err != nil {
		errMsg := fmt.Sprintf("Err: unmarshal body failed, err: %s", err)
		_ = c.log(errMsg)
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, errMsg, restful.MIME_JSON)
		return
	}
	promMsg.SetMessageAt().SetSignature(c.signature)

	channelType := request.QueryParameter("channel_type")
	if channelType == "" {
		errMsg := "Err: no channel_type found"
		_ = c.log(errMsg)
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, errMsg, restful.MIME_JSON)
		return
	}

	senderCreator, exists := senders.ChannelsSenderCreatorMap[channelType]
	if !exists {
		errMsg := fmt.Sprintf("Err: not supported channel_type of (%s)", channelType)
		_ = c.log(errMsg)
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, errMsg, restful.MIME_JSON)
		return
	}

	sender, err := senderCreator(request)
	if err != nil {
		errMsg := fmt.Sprintf("Err: create sender failed, %v", err)
		_ = c.log(errMsg)
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, errMsg, restful.MIME_JSON)
		return
	}

	payload, err := promMsg.ToPayload(channelType, raw)
	if err != nil {
		errMsg := fmt.Sprintf("Err: create msg payload failed, %v", err)
		_ = c.log(errMsg)
		_ = response.WriteHeaderAndJson(http.StatusInternalServerError, errMsg, restful.MIME_JSON)
		return
	}
	if c.debug {
		_, _ = pretty.Println(payload)

		fmt.Println(">>> Payload Markdown")
		fmt.Print(payload.Markdown)
	}

	if err := sender.Send(payload); err != nil {
		errMsg := fmt.Sprintf("Err: sender send failed, %v", err)
		_ = c.log(errMsg)
		_ = response.WriteHeaderAndJson(http.StatusInternalServerError, errMsg, restful.MIME_JSON)
		return
	}

	_ = c.logf("Send succeed: %s\n", request.Request.URL.String())
	response.WriteHeader(http.StatusNoContent)
}
