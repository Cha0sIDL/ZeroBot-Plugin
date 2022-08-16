package game

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/antchfx/htmlquery"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/util"
)

const (
	servicename = "steam"
	searchUrl   = "https://steamstats.cn/api/steam/search?q=%s&page=1&format=json&lang=zh-hans"
	steamUrl    = "https://store.steampowered.com/app/%d/"
)

func init() {
	engine := control.Register(servicename, &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "查询steam游戏(使用英文名称查询更准确哦)" + "- steam xxx\n",
	})
	engine.OnPrefix("steam").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) != 1 {
				ctx.SendChain(message.Text("参数输入有误~"))
				return
			}
			gameName := commandPart[0]
			data, err := web.RequestDataWith(web.NewDefaultClient(), fmt.Sprintf(searchUrl, url.QueryEscape(gameName)), "GET", "https://steamstats.cn/", web.RandUA())
			if err != nil {
				ctx.SendChain(message.Text("出错了", err))
				return
			}
			if len(gjson.ParseBytes(data).Get("data.results").Array()) == 0 {
				ctx.SendChain(message.Text("搜索不到 ", gameName, " 检查下有没有吧~偷偷告诉你，搜英文名的效果可能会更好哟~"))
				return
			}
			results := gjson.ParseBytes(data).Get("data.results.0")
			description := getSteamGameDescription(results.Get("app_id").Int())
			ctx.SendChain(
				message.Text("搜索到以下信息：\n"),
				message.Text("游戏:", results.Get("name"), "(", results.Get("name_cn"), ")\n"),
				message.Text("游戏id：", results.Get("app_id"), "\n"),
				message.Image(results.Get("avatar").String()),
				message.Text("游戏描述：", strings.TrimSpace(description)),
				message.Text(fmt.Sprintf("\nSteamUrl : https://store.steampowered.com/app/%s/", results.Get("app_id"))),
			)
		})
}

func getSteamGameDescription(gameId int64) (description string) {
	doc, _ := htmlquery.LoadURL(fmt.Sprintf(steamUrl, gameId))
	description = htmlquery.InnerText(htmlquery.FindOne(doc, "//*[@id=\"game_highlights\"]/div[1]/div/div[2]"))
	return
}
