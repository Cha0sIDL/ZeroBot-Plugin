// Package choose 选择困难症帮手
package choose

import (
	"fmt"
	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func init() {
	engine := control.Register("choose", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		PrivateDataFolder: "choose",
		Brief:             "选择困难症帮手",
		Help: "choose\n" +
			"- 选择可口可乐还是百事可乐\n" +
			"- 选择肯德基还是麦当劳还是必胜客",
	})
	getData := fcext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		err := file.DownloadTo("https://raw.githubusercontent.com/Cha0sIDL/data/master/what2eat/eating.json",
			engine.DataFolder()+"eating.json")
		if err != nil {
			return false
		}
		err = file.DownloadTo("https://raw.githubusercontent.com/Cha0sIDL/data/master/what2eat/drinks.json",
			engine.DataFolder()+"drinks.json")
		if err != nil {
			return false
		}
		err = file.DownloadTo("https://raw.githubusercontent.com/Cha0sIDL/data/master/what2eat/crazy.json",
			engine.DataFolder()+"crazy.json")
		if err != nil {
			return false
		}
		return true
	})
	engine.OnPrefix("选择").SetBlock(true).Handle(handle)
	engine.OnKeyword("吃什么", getData).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			data, err := os.ReadFile(engine.DataFolder() + "eating.json")
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			foods := gjson.ParseBytes(data).Get("basic_food").Array()
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("建议今天吃:"+lo.Sample(foods).String()))
		})
	engine.OnKeyword("喝什么", getData).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			data, err := os.ReadFile(engine.DataFolder() + "drinks.json")
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			drinks := gjson.ParseBytes(data).Map()
			shop := lo.Sample(lo.Keys(drinks))
			drink := lo.Sample(drinks[shop].Array()).String()
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text(fmt.Sprintf(lo.Sample(
				[]string{
					"不如来杯 %s 的 %s 吧！",
					"去 %s 整杯 %s 吧！",
					"%s 的 %s 如何？",
					"%s 的 %s，好喝绝绝子！",
				}), shop, drink)))
		})
	engine.OnRegex(`疯狂星期(一|二|三|四|五|六|日|天)`, getData).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			data, err := os.ReadFile(engine.DataFolder() + "crazy.json")
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			crazy := gjson.ParseBytes(data).Get("post").Array()
			ctx.SendChain(message.Text(lo.Sample(crazy)))
		})
}

func handle(ctx *zero.Ctx) {
	rawOptions := strings.Split(ctx.State["args"].(string), "还是")
	var options = make([]string, 0)
	for count, option := range rawOptions {
		options = append(options, strconv.Itoa(count+1)+", "+option)
	}
	result := rawOptions[rand.Intn(len(rawOptions))]
	name := ctx.Event.Sender.NickName
	ctx.SendChain(message.Text("> ", name, "\n",
		"你的选项有:", "\n",
		strings.Join(options, "\n"), "\n",
		"你最终会选: ", result,
	))
}
