package aireply

import (
	"encoding/json"
	"github.com/FloatTech/AnimeAPI/aireply"
	"github.com/FloatTech/ZeroBot-Plugin/config"
	"github.com/FloatTech/ZeroBot-Plugin/order"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
	"time"
)

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
			data := map[string]string{"appkey": arg.TTS.Appkey, "access": arg.TTS.Access, "secret": arg.TTS.Secret, "voice": getVoice(), "text": r.TalkPlain(msg)}
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

func getCfg() config.Config {
	return config.Cfg
}

func getVoice() string {
	timeLayout := config.Cfg.TTS.Start
	tmp, _ := time.Parse("2006-01-02", timeLayout)
	login := tmp.Unix()
	today := (time.Now().Unix() - login) / 86400 % int64(len(config.Cfg.TTS.Voice))
	return config.Cfg.TTS.Voice[today]
}
