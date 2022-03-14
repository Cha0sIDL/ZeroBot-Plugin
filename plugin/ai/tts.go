package ai

import (
	"fmt"
	"github.com/FloatTech/AnimeAPI/aireply"
	"github.com/FloatTech/ZeroBot-Plugin/config"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/control/order"
	"github.com/FloatTech/zbputils/file"
	nls "github.com/aliyun/alibabacloud-nls-go-sdk"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	cachePath = dbpath + "cache/"
	dbpath    = "data/ai/"
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
	if file.IsNotExist(cachePath) {
		os.MkdirAll(cachePath, 0755)
	}
	limit := rate.NewManager(time.Second*10, 1)

	en := control.Register("ai", order.AcquirePrio(), &control.Options{
		DisableOnDefault: false,
		Help:             "- @Bot 任意文本(任意一句话回复)",
	})
	en.OnPrefix("复读机").SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			text := ctx.State["args"]
			VoiceFile := cachePath + strconv.FormatInt(ctx.Event.UserID, 10) + strconv.FormatInt(time.Now().Unix(), 10) + ".wav"
			err := util.TTS(VoiceFile, fmt.Sprintf("%v", text), nls.DefaultSpeechSynthesisParam(), getCfg().TTS.Appkey, getCfg().TTS.Access, getCfg().TTS.Secret)
			if err != nil {
				ctx.SendChain(message.Text("Ali NLS 调用失败"))
			} else {
				ctx.SendChain(message.Record("file:///" + file.BOTPATH + "/" + VoiceFile))
			}
		})
	en.OnMessage(zero.OnlyToMe, func(ctx *zero.Ctx) bool {
		return limit.Load(ctx.Event.UserID).Acquire()
	}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			msg := ctx.ExtractPlainText()
			r := aireply.NewAIReply(getReplyMode(ctx))
			arg := nls.DefaultSpeechSynthesisParam()
			arg.Voice = getVoice()
			VoiceFile := cachePath + strconv.FormatInt(ctx.Event.UserID, 10) + strconv.FormatInt(time.Now().Unix(), 10) + ".wav"
			err := util.TTS(VoiceFile, r.TalkPlain(msg, zero.BotConfig.NickName[0]), arg, getCfg().TTS.Appkey, getCfg().TTS.Access, getCfg().TTS.Secret)
			if err != nil {
				//data := map[string]string{"appkey": getCfg().TTS.Appkey, "access": getCfg().TTS.Access, "secret": getCfg().TTS.Secret, "voice": getVoice(), "text": r.TalkPlain(msg)}
				//reqbody, _ := json.Marshal(data)
				// JX3 api 已弃用
				//rsp, _ := util.SendHttp("https://www.jx3api.com/share/aliyun", reqbody)
				//json := gjson.ParseBytes(rsp)
				//if json.Get("code").Int() != 200 {
				//	ctx.SendChain(message.Text(r.TalkPlain(msg)))
				//} else {
				//	ctx.SendChain(message.Record(json.Get("data.url").String()))
				//}
				ctx.SendChain(message.Text(r.TalkPlain(msg, zero.BotConfig.NickName[0])))
			}
			ctx.SendChain(message.Record("file:///" + file.BOTPATH + "/" + VoiceFile))
		})
}

func getCfg() config.Config {
	return config.Cfg
}

func getVoice() string {
	//timeLayout := config.Cfg.TTS.Start
	//tmp, _ := time.Parse("2006-01-02", timeLayout)
	//login := tmp.Unix()
	//today := (time.Now().Unix() - login) / 86400 % int64(len(config.Cfg.TTS.Voice))
	rand.Seed(time.Now().Unix())
	return config.Cfg.TTS.Voice[rand.Intn(len(config.Cfg.TTS.Voice))]
}
