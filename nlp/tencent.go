// Package nlp 腾讯nlp插件
package nlp

import (
	"fmt"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	nlp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/nlp/v20190408"
	"github.com/tidwall/gjson"

	"github.com/FloatTech/ZeroBot-Plugin/config"
)

// Tencent 腾讯nlp结构体
type Tencent struct {
	n string
	b []string
}

const (
	// BotName 腾讯机器人的名字
	BotName = "小龙女"
)

// NewTencent 返回腾讯结构体
func NewTencent(name string, banwords ...string) *Tencent {
	return &Tencent{n: name, b: banwords}
}

func (*Tencent) String() string {
	return "腾讯"
}

// Talk 取得带 CQ 码的回复消息
func (t *Tencent) Talk(_ int64, msg, nickname string) string {
	return t.TalkPlain(0, msg, nickname)
}

// TalkPlain 取得回复消息
func (t *Tencent) TalkPlain(_ int64, msg, nickname string) string {
	credential := common.NewCredential(
		config.Cfg.SecretID,
		config.Cfg.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "nlp.tencentcloudapi.com"
	client, _ := nlp.NewClient(credential, "ap-guangzhou", cpf)

	request := nlp.NewChatBotRequest()

	request.Query = common.StringPtr(msg)

	response, err := client.ChatBot(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return ""
	}
	if err != nil {
		panic(err)
	}
	replystr := gjson.Get(response.ToJsonString(), "Response.Reply").String()
	replystr = strings.ReplaceAll(replystr, BotName, nickname)
	return replystr
}
