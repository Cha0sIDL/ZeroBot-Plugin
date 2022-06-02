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

type Tencent struct{}

const (
	BotName = "小龙女"
)

func (*Tencent) String() string {
	return "腾讯"
}

// Talk 取得带 CQ 码的回复消息
func (t *Tencent) Talk(msg, nickname string) string {
	replystr := t.TalkPlain(msg, nickname)
	replystr = strings.ReplaceAll(replystr, "{face:", "[CQ:face,id=")
	replystr = strings.ReplaceAll(replystr, "{br}", "\n")
	replystr = strings.ReplaceAll(replystr, "}", "]")
	replystr = strings.ReplaceAll(replystr, BotName, nickname)

	return replystr
}

// TalkPlain 取得回复消息
func (t *Tencent) TalkPlain(msg, nickname string) string {
	credential := common.NewCredential(
		config.Cfg.SecretId,
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
	replystr := fmt.Sprintf("%s", gjson.Get(response.ToJsonString(), "Response.Reply"))
	replystr = strings.ReplaceAll(replystr, BotName, nickname)
	return replystr
}
