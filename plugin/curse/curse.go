// Package curse 骂人插件(求骂,自卫)
package curse

import (
	"time"

	"github.com/FloatTech/floatbox/math"
	"github.com/FloatTech/floatbox/process"
	control "github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	fcext "github.com/FloatTech/floatbox/ctxext"
	ctrl "github.com/FloatTech/zbpctrl"
)

const (
	minLevel = "min"
	maxLevel = "max"
)

func init() {
	engine := control.AutoRegister(&ctrl.Options[*zero.Ctx]{
		DisableOnDefault: true,
		Brief:            "骂人反击",
		Help:             "- 骂他@qq(求骂,自卫)",
		PublicDataFolder: "Curse",
	})

	getdb := fcext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		db.DBPath = engine.DataFolder() + "curse.db"
		_, err := engine.GetLazyData("curse.db", true)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		err = db.Open(time.Hour)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		err = db.Create("curse", &curse{})
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		c, err := db.Count("curse")
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		logrus.Infoln("[curse]加载", c, "条骂人语录")
		return true
	})

	engine.OnRegex(`^骂(他|它|她).*?(\d+)`, zero.OnlyGroup, getdb).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		process.SleepAbout1sTo2s()
		qq := math.Str2Int64(ctx.State["regex_matched"].([]string)[2]) // 被骂的人的qq
		for _, su := range zero.BotConfig.SuperUsers {
			if su == qq {
				return
			}
		}
		text := getRandomCurseByLevel(minLevel).Text
		ctx.SendChain(message.At(qq), message.Text(text))
	})
	engine.OnRegex(`^大力骂他.*?(\d+)`, zero.OnlyGroup, getdb).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		process.SleepAbout1sTo2s()
		qq := math.Str2Int64(ctx.State["regex_matched"].([]string)[1]) // 被骂的人的qq
		for _, su := range zero.BotConfig.SuperUsers {
			if su == qq {
				return
			}
		}
		text := getRandomCurseByLevel(maxLevel).Text
		ctx.SendChain(message.At(qq), message.Text(text))
	})

	engine.OnKeywordGroup([]string{"他妈", "公交车", "你妈", "操", "屎", "去死", "快死", "我日", "逼", "尼玛", "艾滋", "癌症", "有病", "烦你", "你爹", "屮", "cnm"}, zero.OnlyToMe, getdb).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			text := getRandomCurseByLevel(minLevel).Text
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(text))
		})
}
