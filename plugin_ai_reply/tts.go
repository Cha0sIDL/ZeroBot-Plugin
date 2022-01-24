package aireply

import (
	"encoding/json"
	"github.com/FloatTech/AnimeAPI/aireply"
	"github.com/FloatTech/ZeroBot-Plugin/order"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io/ioutil"
	"time"
)

var (
	dbpath = "data/tts/"
	dbfile = dbpath + "tts.json"
	login  = time.Now().Unix()
)

type cfg struct {
	Appkey string   `json:"appkey"`
	Access string   `json:"access"`
	Secret string   `json:"secret"`
	Voice  []string `json:"voice"`
}

//func init() {
//	limit := rate.NewManager(time.Second*10, 1)
//
//	control.Register("mockingbird", order.PrioMockingBird, &control.Options{
//		DisableOnDefault: false,
//		Help:             "- @Bot 任意文本(任意一句话回复)",
//	}).OnMessage(zero.OnlyToMe, func(ctx *zero.Ctx) bool {
//		return limit.Load(ctx.Event.UserID).Acquire()
//	}).SetBlock(true).
//		Handle(func(ctx *zero.Ctx) {
//			msg := ctx.ExtractPlainText()
//			r := aireply.NewAIReply(getReplyMode(ctx))
//			ctx.SendChain(mockingbird.Speak(ctx.Event.UserID, func() string {
//				return r.TalkPlain(msg)
//			}))
//		})
//}

func init() {
	limit := rate.NewManager(time.Second*10, 1)

	control.Register("ai", order.PrioMockingBird, &control.Options{
		DisableOnDefault: false,
		Help:             "- @Bot 任意文本(任意一句话回复)",
	}).OnMessage(zero.OnlyToMe, func(ctx *zero.Ctx) bool {
		return limit.Load(ctx.Event.UserID).Acquire()
	}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			msg := ctx.ExtractPlainText()
			r := aireply.NewAIReply(getReplyMode(ctx))
			arg := getCfg()
			data := map[string]string{"appkey": arg.Appkey, "access": arg.Access, "secret": arg.Secret, "voice": getVoice(arg), "text": r.TalkPlain(msg)}
			reqbody, _ := json.Marshal(data)
			rsp, _ := util.SendHttp("https://www.jx3api.com/share/aliyun", reqbody)
			json := gjson.ParseBytes(rsp)
			if json.Get("code").Int() != 200 {
				ctx.SendChain(message.Text(r.TalkPlain(msg)))
			} else {
				ctx.SendChain(message.Record(json.Get("data.url").String()))
			}
		})
}

func getCfg() cfg {
	tmp, err := ioutil.ReadFile(dbfile)
	if err != nil {
		panic("读取文件失败")
	}
	c := new(cfg)
	json.Unmarshal(tmp, c)
	return *c
}

func getVoice(c cfg) string {
	today := (time.Now().Unix() - login) / 86400 % int64(len(c.Voice))
	return c.Voice[today]
}
