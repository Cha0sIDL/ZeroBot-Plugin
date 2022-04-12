// Package curse 骂人插件(求骂,自卫)
package curse

import (
	"github.com/FloatTech/zbputils/math"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"

	control "github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/process"
)

const (
	minLevel = "min"
	maxLevel = "max"
)

func init() {
	engine := control.Register("curse", &control.Options{
		DisableOnDefault: true,
		Help:             "骂他@qq(求骂,自卫)",
		PublicDataFolder: "Curse",
	})

	go func() {
		dbpath := engine.DataFolder()
		db.DBPath = dbpath + "curse.db"
		_, err := file.GetLazyData(db.DBPath, false, true)
		if err != nil {
			panic(err)
		}
		err = db.Create("curse", &curse{})
		if err != nil {
			panic(err)
		}
		c, _ := db.Count("curse")
		logrus.Infoln("[curse]加载", c, "条骂人语录")
	}()

	engine.OnRegex(`^骂(他|它|她).*?(\d+)`, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		process.SleepAbout1sTo2s()
		qq := math.Str2Int64(ctx.State["regex_matched"].([]string)[1]) // 被骂的人的qq
		for _, su := range zero.BotConfig.SuperUsers {
			if su == strconv.FormatInt(qq, 10) {
				return
			}
		}
		text := getRandomCurseByLevel(minLevel).Text
		ctx.SendChain(message.At(qq), message.Text(text))
	})
	engine.OnRegex(`^大力骂他.*?(\d+)`, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		process.SleepAbout1sTo2s()
		qq := math.Str2Int64(ctx.State["regex_matched"].([]string)[1]) // 被骂的人的qq
		for _, su := range zero.BotConfig.SuperUsers {
			if su == strconv.FormatInt(qq, 10) {
				return
			}
		}
		text := getRandomCurseByLevel(maxLevel).Text
		ctx.SendChain(message.At(qq), message.Text(text))
	})

	//engine.OnFullMatch("大力骂我").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
	//	process.SleepAbout1sTo2s()
	//	text := getRandomCurseByLevel(maxLevel).Text
	//	ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(text))
	//})

	engine.OnKeywordGroup([]string{"他妈", "公交车", "你妈", "操", "屎", "去死", "快死", "我日", "逼", "尼玛", "艾滋", "癌症", "有病", "烦你", "你爹", "屮", "cnm"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			text := getRandomCurseByLevel(maxLevel).Text
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(text))
		})
}
