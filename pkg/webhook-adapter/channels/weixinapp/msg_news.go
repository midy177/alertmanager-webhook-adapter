package weixinapp

import (
	"alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"alertmanager-webhook-adapter/pkg/webhook-adapter/utils"
)

func init() {
	Payload2MsgFnMap[MsgTypeNews] = NewMsgNewsFromPayload
}

type News struct {
	Articles []*Article `json:"articles"` // 图文消息，一个图文消息支持1到8条图文
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`    // 点击后跳转的链接
	PicURL      string `json:"picurl"` // 图文消息的图片链接，支持JPG、PNG格式，较好的效果为大图 1068*455，小图150*150
}

func NewArticle(title string, url string) *Article {
	return &Article{
		Title: utils.TruncateToValidUTF8(title, maxTitleBytes, truncatedMark),
		URL:   url,
	}
}

func (a *Article) SetDescription(descr string) *Article {
	a.Description = utils.TruncateToValidUTF8(descr, maxDescriptionBytes, truncatedMark)
	return a
}

func (a *Article) SetPicURL(picURL string) *Article {
	a.PicURL = picURL
	return a
}

func NewMsgNews(articles []*Article) *Msg {
	var a []*Article

	if len(articles) > 8 {
		a = articles[:8]
	} else {
		a = articles
	}

	return &Msg{
		MsgType: MsgTypeNews,
		News: &News{
			Articles: a,
		},
	}

}

func NewMsgNewsFromPayload(payload *models.Payload) *Msg {
	// Todo, construct articles from payload
	articles := []*Article{}
	return NewMsgNews(articles)
}
