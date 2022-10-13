package ai

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	nls "github.com/aliyun/alibabacloud-nls-go-sdk"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/config"
	"github.com/FloatTech/ZeroBot-Plugin/nlp"
	"github.com/FloatTech/ZeroBot-Plugin/util"
)

const (
	cachePath = dbpath + "cache/"
	dbpath    = "data/ai/"
)

// func init() {
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
	en := control.Register("ai", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "- @Bot 任意文本(任意一句话回复)",
	})
	en.OnPrefix("复读").SetBlock(true).Handle(
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
	en.OnMessage(zero.OnlyToMe, zero.OnlyGroup).SetBlock(false).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			msg := ctx.ExtractPlainText()
			r := nlp.NewAIReply(getReplyMode(ctx))
			if util.Rand(1, 100) < 50 {
				ctx.SendChain(message.Text(r.TalkPlain(msg, zero.BotConfig.NickName[0])))
				return
			}
			arg := nls.DefaultSpeechSynthesisParam()
			arg.Voice = getVoice()
			VoiceFile := cachePath + strconv.FormatInt(ctx.Event.UserID, 10) + strconv.FormatInt(time.Now().Unix(), 10) + ".wav"
			err := util.TTS(VoiceFile, r.TalkPlain(msg, zero.BotConfig.NickName[0]), arg, getCfg().TTS.Appkey, getCfg().TTS.Access, getCfg().TTS.Secret)
			if err != nil {
				ctx.SendChain(message.Text(r.TalkPlain(msg, zero.BotConfig.NickName[0])))
			}
			ctx.SendChain(message.Record("file:///" + file.BOTPATH + "/" + VoiceFile))
		})
}

func getCfg() config.Config {
	return config.Cfg
}

func getVoice() string {
	// timeLayout := config.Cfg.TTS.Start
	// tmp, _ := time.Parse("2006-01-02", timeLayout)
	// login := tmp.Unix()
	// today := (time.Now().Unix() - login) / 86400 % int64(len(config.Cfg.TTS.Voice))
	rand.Seed(time.Now().Unix())
	return config.Cfg.TTS.Voice[rand.Intn(len(config.Cfg.TTS.Voice))]
}
