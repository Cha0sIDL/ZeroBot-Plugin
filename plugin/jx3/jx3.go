// Package jx3 剑网相关插件
package jx3

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	goUrl "net/url"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/go-resty/resty/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/playwright-community/playwright-go"
	"github.com/samber/lo"
	"github.com/tidwall/sjson"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	"github.com/FloatTech/floatbox/process"

	"github.com/FloatTech/ZeroBot-Plugin/config"

	"github.com/antchfx/htmlquery"

	"github.com/fumiama/cron"

	ctrl "github.com/FloatTech/zbpctrl"

	binutils "github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/zbputils/control"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/util"

	"github.com/fogleman/gg"
)

var tuiKey = map[string]string{
	"大战":     "60f211c82d105c0014c5dd7d",
	"武林通鉴秘境": "60f211c82d105c0014c5de01",
	"武林通鉴公共": "60f211c82d105c0014c5dd97",
	"十人团队秘境": "60f211c82d105c0014c5ddcd",
	"阵营日常":   "60f211c82d105c0014c5dd9d",
}

var chatServer = map[string]string{
	"绝代天骄": "45",
	"唯我独尊": "19",
}

var allServer = map[string][2]string{
	"斗转星移": {"斗转星移", "电信五区"},
	"姨妈":   {"斗转星移", "电信五区"},
	"蝶恋花":  {"蝶恋花", "电信一区"},
	"龙争虎斗": {"龙争虎斗", "电信一区"},
	"长安城":  {"长安城", "电信一区"},
	"幽月轮":  {"幽月轮", "电信五区"},
	"剑胆琴心": {"剑胆琴心", "电信五区"},
	"煎蛋":   {"剑胆琴心", "电信五区"},
	"乾坤一掷": {"乾坤一掷", "电信五区"},
	"华乾":   {"乾坤一掷", "电信五区"},
	"唯我独尊": {"唯我独尊", "电信五区"},
	"唯满侠":  {"唯我独尊", "电信五区"},
	"梦江南":  {"梦江南", "电信五区"},
	"双梦":   {"梦江南", "电信五区"},
	"绝代天骄": {"绝代天骄", "电信八区"},
	"绝代":   {"绝代天骄", "电信八区"},
	"破阵子":  {"破阵子", "双线一区"},
	"念破":   {"破阵子", "双线一区"},
	"天鹅坪":  {"天鹅坪", "双线一区"},
	"纵月":   {"天鹅坪", "双线一区"},
	"飞龙在天": {"飞龙在天", "双线二区"},
	"大唐万象": {"大唐万象", "电信五区"},
	"青梅煮酒": {"青梅煮酒", "双线四区"},
	"共結來緣": {"共結來緣"},
	"傲血戰意": {"傲血戰意"},
	"巔峰再起": {"巔峰再起"},
	"江海雲夢": {"江海雲夢"},
}

var serverIP = map[string]string{
	"斗转星移": "125.88.195.133:3724",
	"蝶恋花":  "125.88.195.112:3724",
	"龙争虎斗": "125.88.195.69:3724",
	"长安城":  "125.88.195.52:3724",
	"幽月轮":  "125.88.195.117:3724",
	"剑胆琴心": "125.88.195.42:3724",
	"乾坤一掷": "125.88.195.154:3724",
	"唯我独尊": "125.88.195.89:3724",
	"梦江南":  "125.88.195.59:3724",
	"绝代天骄": "125.88.195.178:3724",
	"破阵子":  "103.228.229.128:3724",
	"天鹅坪":  "103.228.229.129:3724",
	"飞龙在天": "103.228.229.130:3724",
	"青梅煮酒": "103.228.229.127:3724",
}

type cd struct {
	last     int64
	fileName string
}

type jinPrice struct {
	P5173    float64 `json:"5173"`
	Post     float64 `json:"post"`
	Official float64 `json:"official"`
	Date     string  `json:"date"`
}

type sandBox struct {
	sandToken       string
	sandTokenExpire int64
}

var sand sandBox

var controlCd = make(map[string]cd)

var xiaoheiIndx = map[string]string{
	"电信点卡": "server1",
	"双线一区": "server2",
	"电信一区": "server3",
	"双线二区": "server4",
}

type xiaohei struct {
	State int `json:"state"`
	Data  struct {
		Other []struct {
			Region struct {
				ID          int    `json:"id"`
				CreatedTime string `json:"createdTime"`
				UpdatedTime string `json:"updatedTime"`
				RegionName  string `json:"regionName"`
				RegionNick  string `json:"regionNick"`
				Charge      string `json:"charge"`
			} `json:"region"`
			Prices struct {
				ID          int         `json:"id"`
				Price       float64     `json:"price"`
				Region      string      `json:"region"`
				RegionAlias string      `json:"regionAlias"`
				RegionID    int         `json:"regionId"`
				Server      string      `json:"server"`
				ServerID    int         `json:"serverId"`
				SaleCode    string      `json:"saleCode"`
				TradeTime   string      `json:"tradeTime"`
				OutwardName interface{} `json:"outwardName"`
				OutwardID   int         `json:"outwardId"`
				Audit       int         `json:"audit"`
				Now         int         `json:"now"`
				Exterior    string      `json:"exterior"`
				Pricer      string      `json:"pricer"`
			} `json:"prices"`
		} `json:"other"`
		Region struct {
			ID          int    `json:"id"`
			CreatedTime string `json:"createdTime"`
			UpdatedTime string `json:"updatedTime"`
			RegionName  string `json:"regionName"`
			RegionNick  string `json:"regionNick"`
			Charge      string `json:"charge"`
		} `json:"region"`
		Prices []struct {
			ID          int         `json:"id"`
			Price       float64     `json:"price"`
			Region      string      `json:"region"`
			RegionAlias string      `json:"regionAlias"`
			RegionID    int         `json:"regionId"`
			Server      string      `json:"server"`
			ServerID    int         `json:"serverId"`
			SaleCode    string      `json:"saleCode"`
			TradeTime   string      `json:"tradeTime"`
			OutwardName interface{} `json:"outwardName"`
			OutwardID   int         `json:"outwardId"`
			Audit       int         `json:"audit"`
			Now         int         `json:"now"`
			Exterior    string      `json:"exterior"`
			Pricer      string      `json:"pricer"`
		} `json:"prices"`
	} `json:"data"`
	Message interface{} `json:"message"`
}

// GroupList 获取群配置的结构
type GroupList struct {
	grp    int64
	server string
}

func init() {
	en := control.Register("jx3", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		PrivateDataFolder: "jx3",
		Brief:             "剑网相关插件",
		Help: "- 日常任务|日常(eg 日常任务绝代天骄)\n" +
			"- 开服\n" +
			"- 更新公告\n" +
			"- 技改\n" +
			"- 小药\n" +
			"- 金价|金价查询 xxx (eg 金价 绝代天骄)\n" +
			"- 宏xxx(eg 宏分山劲)\n" +
			"- 沙盘xxx(eg 沙盘绝代天骄)\n" +
			"- 奇遇|奇遇查询 xxx xxx(eg 奇遇 唯我独尊 柳连柳奶)\n" +
			"- 攻略xxx(eg 攻略三山四海)\n" +
			"- 骚话\n" +
			"- 舔狗\n" +
			"-（开启|关闭）jx推送\n" +
			"- /roll 随机roll点\n" +
			"- 物价xxx\n" +
			"- 绑定区服xxx\n" +
			"- 团队相关见 https://docs.qq.com/doc/DUGJRQXd1bE5YckhB",
	}).ApplySingle(ctxext.DefaultSingle)
	c := cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger), cron.SkipIfStillRunning(cron.DefaultLogger)))
	_, err := c.AddFunc("0 5 * * *", func() {
		err := updateTalk()
		if err != nil {
			return
		}
	})
	c.AddFunc("@every 30s", func() { //nolint:errcheck
		zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
			var grpList []GroupList
			controls := jdb.isEnable()
			for _, g := range ctx.GetGroupList().Array() {
				grp := g.Get("group_id").Int()
				if val, ok := controls[grp]; ok {
					grpList = append(grpList, GroupList{
						grp:    grp,
						server: val,
					})
				}
			}
			news(ctx, grpList)
			return true
		})
	})
	c.AddFunc("@every 2m", func() { //nolint:errcheck
		zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
			var grpList []GroupList
			controls := jdb.isEnable()
			for _, g := range ctx.GetGroupList().Array() {
				grp := g.Get("group_id").Int()
				if val, ok := controls[grp]; ok {
					grpList = append(grpList, GroupList{
						grp:    grp,
						server: val,
					})
				}
			}
			checkServer(ctx, grpList)
			return true
		})
	})
	if err == nil && runtime.GOOS == "linux" {
		c.Start()
	}
	jdb = initialize()
	if config.Cfg.JxChat != nil {
		for _, chat := range *config.Cfg.JxChat {
			go startChatWs(chat)
		}
	}
	datapath := file.BOTPATH + "/" + en.DataFolder()
	en.OnFullMatchGroup([]string{"日常", "日常任务"}, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			decorator(daily)(ctx)
		})
	en.OnFullMatch("开服").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			decorator(server)(ctx)
		})
	en.OnFullMatch("更新公告").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			pic := util.ScreenShot("https://jx3.xoyo.com/launcher/update/latest.html")
			ctx.SendChain(message.ImageBytes(pic))
		})
	en.OnFullMatch("技改").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			pic := util.ScreenShot("https://jx3.xoyo.com/launcher/update/latest_exp.html")
			ctx.SendChain(message.ImageBytes(pic))
		})
	en.OnPrefixGroup([]string{"金价", "金价查询"}).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			jinjia(ctx, datapath)
		})
	en.OnPrefixGroup([]string{"沙盘"}).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) != 1 {
				ctx.SendChain(message.Text("参数输入有误！\n" + "沙盘 绝代天骄"))
				return
			}
			server := commandPart[0]
			if fullName, ok := allServer[server]; ok {
				if len(sand.sandToken) == 0 || carbon.Now().Timestamp() > sand.sandTokenExpire {
					login, err := web.GetData("https://www.j3sp.com/api/user/login?account=ChrisSandBox%40outlook.com&password=123456")
					if err != nil || gjson.ParseBytes(login).Get("code").Int() != 1 {
						log.Errorln("jx3daily:", err)
						return
					}
					sand = sandBox{
						sandToken:       gjson.ParseBytes(login).Get("data.userinfo.token").String(),
						sandTokenExpire: carbon.Now().Timestamp() + 43200,
					}
				}
				client := NewTimeOutDefaultClient()
				request, err := http.NewRequest("GET", fmt.Sprintf("https://www.j3sp.com/api/sand/?serverName=%s&shadow=0&is_history=1", fullName[0]), nil)
				if err == nil {
					// 增加header选项
					var response *http.Response
					request.Header.Add("Cookie", fmt.Sprintf("spc_token=%s", sand.sandToken))
					response, err = client.Do(request)
					if err == nil {
						if response.StatusCode != http.StatusOK {
							ctx.SendChain(message.Text("请求出错了稍后再试试吧~"))
							return
						}
						data, _ := io.ReadAll(response.Body)
						response.Body.Close()
						strData := binutils.BytesToString(data)
						if gjson.Get(strData, "msg").String() != "success" {
							ctx.SendChain(message.Text("请求出错了稍后再试试吧~"))
							return
						}
						ctx.SendChain(message.Image(gjson.Get(strData, "data.sand_data.sandImage").String()))
					}
				}
			} else {
				ctx.Send(message.Text("区服输入有误"))
			}
		})
	en.OnSuffix("小药").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Image(fileURL + "medicine.png"))
		})
	en.OnPrefix("宏").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			name := ctx.State["args"].(string)
			mental := jdb.getMentalData(strings.ReplaceAll(name, " ", ""))
			mentalURL := fmt.Sprintf("https://cms.jx3box.com/api/cms/posts?type=macro&per=10&page=1&order=update&client=std&search=%s", goUrl.QueryEscape(mental.Name))
			data, err := web.RequestDataWith(NewTimeOutDefaultClient(), mentalURL, "GET", "application/x-www-form-urlencoded", web.RandUA(), nil)
			DataList := gjson.Get(binutils.BytesToString(data), "data.list").Array()
			if err != nil || len(DataList) == 0 {
				ctx.SendChain(message.Text("出错了请检查参数或稍后试试吧~"))
				return
			}
			for idx, m := range DataList {
				rsp := ""
				if idx >= 2 {
					break
				}
				rsp += "作者:" + m.Get("author").String() + "\n" + m.Get("post_title").String() + "\n"
				for _, meta := range m.Get("post_meta.data").Array() {
					rsp += "\n" + meta.Get("desc").String() + "\n" + meta.Get("name").String() + ":\n" + meta.Get("macro").String() + "\n" + "----------------------------------------------\n"
				}
				rsp += fmt.Sprintf("详情请点击: https://www.jx3box.com/macro/%s", m.Get("ID").String()) + "\n"
				rsp += "数据来源于JXBOX，dps请自行测试"
				ctx.SendChain(message.Text(rsp))
				time.Sleep(time.Second * 4)
			}
		})
	// en.OnRegex(`^攻略(.*)`).SetBlock(true).
	//	Handle(func(ctx *zero.Ctx) {
	//		name := ctx.State["regex_matched"].([]string)[1]
	//		if len(name) == 0 {
	//			ctx.SendChain(message.Text("输入参数有误！！！"))
	//		} else {
	//			dbData := getAdventure(name)
	//			if len(dbData.Pic) == 0 || carbon.Now().DiffAbsInSeconds(carbon.CreateFromTimestamp(dbData.Time)) > 3600*10 {
	//				dwData, _ := web.GetData(fmt.Sprintf("https://node.jx3box.com/serendipities?name=%s", goUrl.QueryEscape(name)))
	//				dwList := gjson.Get(binutils.BytesToString(dwData), "list").Array()
	//				if len(dwList) == 0 {
	//					ctx.SendChain(message.Text(fmt.Sprintf("没有找到%s呢，你是不是乱输的哦~", name)))
	//					return
	//				}
	//				dwId := dwList[0].Get("dwID").String()
	//				json, _ := web.GetData("https://icon.jx3box.com/pvx/serendipity/output/serendipity.json")
	//				articleId := gjson.Get(binutils.BytesToString(json), dwId).String()
	//				articleUrl := fmt.Sprintf("https://www.jx3box.com/cj/#/view/%s", articleId)
	//				pw, err := playwright.Run()
	//				if err != nil {
	//					playwright.Install()
	//					playwright.Run()
	//				}
	//				defer pw.Stop()
	//				browser, err := pw.Chromium.Launch()
	//				if err != nil {
	//					playwright.Install()
	//				}
	//				page, err := browser.NewPage(playwright.BrowserNewContextOptions{
	//					IsMobile: playwright.Bool(true),
	//				})
	//				if err != nil {
	//					return
	//				}
	//				_, err = page.Goto(articleUrl, playwright.PageGotoOptions{
	//					WaitUntil: playwright.WaitUntilStateNetworkidle,
	//					Timeout:   playwright.Float(10000),
	//				})
	//				if err != nil {
	//					return
	//				}
	//				page.Click("//*[@id=\"app\"]/aside/span")
	//				result, _ := page.QuerySelector("div[class=\"c-article-chunk on\"]")
	//				result.WaitForSelector("image")
	//				result.ScrollIntoViewIfNeeded()
	//				b, err := result.Screenshot()
	//				if err != nil {
	//					ctx.SendChain(message.Text("出错了，稍后再试试吧~"))
	//				}
	//				db := &Adventure{
	//					Name: name,
	//					Pic:  b,
	//					Time: carbon.Now().Timestamp(),
	//				}
	//				updateAdventure(db)
	//				ctx.SendChain(message.ImageBytes(b))
	//			} else {
	//				ctx.SendChain(message.ImageBytes(dbData.Pic))
	//			}
	//		}
	//	})
	en.OnPrefix("攻略").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) != 1 {
				ctx.SendChain(message.Text("输入参数有误！！！攻略 炼狱厨神"))
				return
			}
			name := commandPart[0]
			dbData := jdb.getAdventure(name)
			if len(dbData.Pic) == 0 || carbon.Now().DiffAbsInSeconds(carbon.CreateFromTimestamp(dbData.Time)) > 3600*10 {
				dwData, _ := web.GetData(fmt.Sprintf("https://node.jx3box.com/serendipities?name=%s", goUrl.QueryEscape(name)))
				dwList := gjson.Get(binutils.BytesToString(dwData), "list").Array()
				if len(dwList) == 0 {
					ctx.SendChain(message.Text(fmt.Sprintf("没有找到%s呢，你是不是乱输的哦~", name)))
					return
				}
				dwID := dwList[0].Get("dwID").String()
				json, _ := web.GetData("https://icon.jx3box.com/pvx/serendipity/output/serendipity.json")
				articleID := gjson.Get(binutils.BytesToString(json), dwID).String()
				articleURL := fmt.Sprintf("https://www.jx3box.com/cj/#/view/%s", articleID)
				pw, err := playwright.Run()
				if err != nil {
					playwright.Install() //nolint:errcheck
					playwright.Run()     //nolint:errcheck
				}
				defer pw.Stop() //nolint:errcheck
				browser, err := pw.Chromium.Launch()
				if err != nil {
					playwright.Install() //nolint:errcheck
				}
				page, err := browser.NewPage(playwright.BrowserNewContextOptions{
					IsMobile: playwright.Bool(true),
				})
				if err != nil {
					return
				}
				_, err = page.Goto(articleURL, playwright.PageGotoOptions{
					WaitUntil: playwright.WaitUntilStateNetworkidle,
					Timeout:   playwright.Float(30000),
				})
				if err != nil {
					ctx.SendChain(message.Text("Err:", err))
					return
				}
				err = page.Click("//*[@id=\"app\"]/aside/span")
				if err != nil {
					ctx.SendChain(message.Text("Err:", err))
					return
				}
				result, _ := page.QuerySelector("div[class=\"c-article-chunk on\"]")
				html, _ := result.InnerHTML()
				htmlPage, err := browser.NewPage()
				if err != nil {
					return
				}
				htmlPage.SetContent(html, playwright.PageSetContentOptions{ //nolint:errcheck
					WaitUntil: playwright.WaitUntilStateNetworkidle,
				})
				htmlPage.Keyboard().Down("PageDown") //nolint:errcheck
				time.Sleep(time.Second * 10)
				htmlPage.Keyboard().Up("PageDown") //nolint:errcheck
				b, err := htmlPage.Screenshot(
					playwright.PageScreenshotOptions{
						Type:     playwright.ScreenshotTypeJpeg,
						Quality:  playwright.Int(100),
						FullPage: playwright.Bool(true),
					})
				if err != nil {
					ctx.SendChain(message.Text("出错了，稍后再试试吧~"))
				}
				db := &Adventure{
					Name: name,
					Pic:  b,
					Time: carbon.Now().Timestamp(),
				}
				jdb.updateAdventure(db)
				ctx.SendChain(message.ImageBytes(b))
			} else {
				ctx.SendChain(message.ImageBytes(dbData.Pic))
			}
		})
	en.OnPrefix("交易行").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			var server string
			var itemName string
			switch {
			case len(commandPart) == 1:
				server = jdb.bind(ctx.Event.GroupID)
				itemName = commandPart[0]
				if len(server) == 0 {
					ctx.SendChain(message.Text("本群尚未绑定区服"))
					return
				}
			case len(commandPart) == 2:
				server = commandPart[0]
				itemName = commandPart[1]
			default:
				ctx.SendChain(message.Text("参数输入有误！\n" + "战绩 绝代天骄 xxx"))
				return
			}
			if normServer, ok := allServer[server]; ok {
				itemURL := fmt.Sprintf("https://helper.jx3box.com/api/item/search?page=1&limit=15&client=std&keyword=%s", goUrl.QueryEscape(itemName))
				data, err := web.RequestDataWith(NewTimeOutDefaultClient(), itemURL, "GET", "https://www.jx3box.com/", web.RandUA(), nil)
				jsonItem := gjson.ParseBytes(data)
				if err != nil || jsonItem.Get("code").Int() != 200 {
					ctx.SendChain(util.HTTPError()...)
					return
				}
				jsonItemArr := jsonItem.Get("data.data").Array()
				if len(jsonItemArr) == 0 {
					ctx.SendChain(message.Text("没有找到" + itemName + "你是不是乱输的哦~"))
					return
				}
				var rsp = "找到以下道具,请输入序号选择要查询的物品：\n"
				for idx, itemInfo := range jsonItemArr {
					rsp += fmt.Sprintf("%d. %s %s\n", idx, itemInfo.Get("Name").String(), util.GetChinese(itemInfo.Get("Desc").String()))
				}
				next := zero.NewFutureEvent("message", 2, false, zero.RegexRule(`^\d+$`), ctx.CheckSession())
				recv, cancel := next.Repeat()
				defer cancel()
				ctx.SendChain(message.Text(rsp))
				for {
					select {
					case <-time.After(time.Second * 120):
						ctx.SendChain(message.Text("交易行查询指令过期"))
						return
					case c := <-recv:
						msg := c.Event.Message.ExtractPlainText()
						num, err := strconv.Atoi(msg)
						if err != nil {
							ctx.SendChain(message.Text("请输入数字!"))
							continue
						}
						if num < 0 || num > len(jsonItemArr) {
							ctx.SendChain(message.Text("序号非法!"))
							continue
						}
						ItemID := jsonItem.Get(fmt.Sprintf("data.data.%d.id", num)).String()
						// ItemIcon := jsonItem.Get(fmt.Sprintf("data.%d.IconID", num)).Int()
						tradingPrice, err := web.GetData(fmt.Sprintf("https://next2.jx3box.com/api/item-price/%s/detail?server=%s", ItemID, goUrl.QueryEscape(normServer[0])))
						priceItem := gjson.ParseBytes(tradingPrice)
						priceSlice := priceItem.Get("data.prices").Array()
						if err != nil || priceItem.Get("code").Int() != 0 {
							ctx.SendChain(util.HTTPError()...)
							return
						}
						rsp = ""
						sort.SliceStable(priceSlice, func(i, j int) bool {
							return priceSlice[i].Get("unit_price").Int() < priceSlice[j].Get("unit_price").Int()
						})
						for _, s := range priceSlice {
							rsp += fmt.Sprintf("单价：%s , 数量：%d\n", price2hRead(s.Get("unit_price").Int()), s.Get("n_count").Int())
						}
						rsp += "-------------数据来自JXBOX----------------\n"
						ctx.SendChain(message.Text(rsp))
						return
					}
				}
			} else {
				ctx.SendChain(message.Text("输入区服有误，请检查qaq~"))
			}
		})
	en.OnRegex(`^(?i)骚话(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var t Jokes
			jdb.Pick(&t)
			ctx.SendChain(message.Text(t.Talk))
		})
	en.OnFullMatch("/roll").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(fmt.Sprintf("%s 投出了%d点。", ctx.Event.Sender.NickName, util.Rand(1, 100))), message.Reply(ctx.Event.MessageID))
		})
	en.OnFullMatch("开启jx推送", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			area := enable(ctx.Event.GroupID)
			if len(area) == 0 {
				var server []string
				for key := range serverIP {
					server = append(server, key)
				}
				ctx.Send(message.Text("开启成功，检测到当前未绑定区服，请输入\n绑定区服xxx\n进行绑定，可选服务器有：\n" + util.PrettyPrint(server)))
			} else {
				ctx.Send(message.Text("开启成功，当前绑定区服为：" + area))
			}
		})
	en.OnFullMatch("关闭jx推送", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			jdb.disable(ctx.Event.GroupID)
			ctx.Send(message.Text("关闭成功"))
		})
	en.OnPrefix("绑定区服", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			area := strings.ReplaceAll(ctx.State["args"].(string), " ", "")
			if val, ok := allServer[area]; ok {
				jdb.bindArea(ctx.Event.GroupID, val[0])
				ctx.Send(message.Text("绑定成功"))
			} else {
				ctx.Send(message.Text("区服输入有误"))
			}
		})
	// 开团 副本名 备注
	en.OnPrefixGroup([]string{"开团", "新建团队", "创建团队"}, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) != 2 {
				ctx.SendChain(message.Text("开团参数输入有误！开团 副本名 备注"))
				return
			}
			dungeon := commandPart[0]
			comment := commandPart[1]
			teamID, err := jdb.createNewTeam(&Team{
				LeaderID: ctx.Event.UserID,
				Dungeon:  dungeon,
				Comment:  comment,
				GroupID:  ctx.Event.GroupID,
			})
			if err != nil {
				ctx.SendChain(message.Text("Error :", err))
				return
			}
			ctx.SendChain(message.Text("开团成功，团队id为：", teamID))
		})
	// 报团 团队ID 心法 角色名 [是否双休] 按照报名时间先后默认排序 https://docs.qq.com/doc/DUGJRQXd1bE5YckhB
	en.OnPrefixGroup([]string{"报名", "报团", "报名团队", "代报名"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			double := 0
			switch {
			case len(commandPart) == 3:
				double = 0
			case len(commandPart) == 4:
				double, _ = strconv.Atoi(commandPart[3])
			default:
				ctx.SendChain(message.Text("报团参数有误。"))
				return
			}
			teamID, err := strconv.Atoi(commandPart[0])
			if err != nil {
				ctx.SendChain(message.Text("团队编号输入有误"))
				return
			}
			mental := jdb.getMentalData(commandPart[1])
			nickName := commandPart[2]
			if mental.ID == 0 {
				ctx.SendChain(message.Text("心法输入有误"))
				return
			}
			Team := jdb.getTeamInfo(teamID)
			if Team.GroupID != ctx.Event.GroupID {
				ctx.SendChain(message.Text("当前团队不存在。"))
				return
			}
			if []rune(ctx.MessageString())[0] == '代' {
				if jdb.isInTeam("team_id = ? and member_nick_name = ?", teamID, nickName) {
					ctx.SendChain(message.Text(nickName, "已经在团队中了。"))
					return
				}
			} else {
				if jdb.isInTeam("team_id = ? and member_qq = ?", teamID, ctx.Event.UserID) {
					ctx.SendChain(message.Text("你已经在团队中了。"))
					return
				}
			}
			var member = Member{
				TeamID:         uint(teamID),
				MemberQQ:       ctx.Event.UserID,
				MemberNickName: nickName,
				MentalID:       mental.ID,
				Double:         double,
				SignUp:         carbon.Now().Timestamp(),
			}
			err = jdb.addMember(&member)
			if err != nil {
				ctx.SendChain(message.Text("数据库写入失败,Err:", err))
				return
			}
			ctx.SendChain(message.Text("报团成功"), message.Reply(ctx.Event.MessageID))
			ctx.SendChain(message.Text("当前团队:\n"), message.Image("base64://"+helper.BytesToString(util.Image2Base64(drawTeam(teamID)))))
		})
	en.OnPrefixGroup([]string{"撤销报团"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			teamID, _ := strconv.Atoi(commandPart[0])
			team := jdb.getTeamInfo(teamID)
			if team.GroupID != ctx.Event.GroupID {
				ctx.SendChain(message.Text("参数输入有误。"))
				return
			}
			err := jdb.deleteMember(teamID, ctx.Event.UserID)
			if err != nil {
				ctx.SendChain(message.Text("Err:", err))
				return
			}
			ctx.SendChain(message.Text("撤销成功"), message.Reply(ctx.Event.MessageID))
			ctx.SendChain(message.Text("当前团队:\n"), message.Image("base64://"+helper.BytesToString(util.Image2Base64(drawTeam(teamID)))))
		})
	en.OnFullMatchGroup([]string{"我报的团", "我的报名"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var sTeam []Team
			err := jdb.Find("groupId = ?", &sTeam, ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Text("Err:", err))
			}
			s := lo.Map(sTeam, func(item Team, _ int) uint {
				return item.TeamID
			})
			SignUp := lo.Uniq(jdb.getSignUp(ctx.Event.UserID))

			ctx.SendChain(message.Text("本群你报名过的团队id：\n", util.IntersectArray(s, SignUp)))
		})
	en.OnFullMatchGroup([]string{"我的开团"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var sTeam []Team
			err := jdb.Find("leaderId = ? and groupId = ?", &sTeam, ctx.Event.UserID, ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Text("Err:", err))
				return
			}
			out := ""
			for _, data := range sTeam {
				out += fmt.Sprintf("团队id：%d,团长 ：%d,副本：%s，备注：%s\n",
					data.TeamID, data.LeaderID, data.Dungeon, data.Comment)
			}
			ctx.SendChain(message.Text(out))
		})
	en.OnFullMatchGroup([]string{"查看全部团队"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var sTeam []Team
			err := jdb.Find("groupId = ?", &sTeam, ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Text("Err:", err))
				return
			}
			if len(sTeam) == 0 {
				ctx.SendChain(message.Text("本群没有有效团队哦"))
				return
			}
			out := ""
			for _, data := range sTeam {
				out += fmt.Sprintf("团队id：%d,团长 ：%d,副本：%s，备注：%s\n",
					data.TeamID, data.LeaderID, data.Dungeon, data.Comment)
			}
			ctx.SendChain(message.Text(out))
		})
	// 查看团队 teamid
	en.OnPrefixGroup([]string{"查看团队", "查询团队", "查团"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			teamID, _ := strconv.Atoi(commandPart[0])
			team := jdb.getTeamInfo(teamID)
			if team.GroupID != ctx.Event.GroupID {
				ctx.SendChain(message.Text("团队id输入有误。"))
				return
			}
			ctx.SendChain(message.Image("base64://" + helper.BytesToString(util.Image2Base64(drawTeam(teamID)))))
		})
	// 取消开团 团队id
	en.OnPrefixGroup([]string{"取消开团", "删除团队", "结束团队"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) < 1 {
				ctx.SendChain(message.Text("参数有误"))
			}
			teamID, err := strconv.Atoi(commandPart[0])
			team := jdb.getTeamInfo(teamID)
			if err != nil || team.GroupID != ctx.Event.GroupID || team.LeaderID != ctx.Event.UserID {
				ctx.SendChain(message.Text("团队id输入有误"))
				return
			}
			err = jdb.delTeam(teamID, ctx.Event.UserID)
			if err != nil {
				ctx.SendChain(message.Text(err))
			}
			ctx.SendChain(message.Text("取消成功"))
		})
	en.OnPrefixGroup([]string{"奇遇", "奇遇查询"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var msg string
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) != 2 {
				ctx.SendChain(message.Text("参数输入有误！\n" + "奇遇 唯我独尊 柳连柳奶"))
				return
			}
			server := commandPart[0]
			name := commandPart[1]
			qiyuURL := fmt.Sprintf("https://pull.j3cx.com/api/serendipity?server=%s&role=%s&pageIndex=1&pageSize=30", server, name)
			rspData, err := web.RequestDataWith(NewTimeOutDefaultClient(), qiyuURL, "GET", "", web.RandUA(), nil)
			if err != nil || gjson.Get(binutils.BytesToString(rspData), "code").Int() != 0 {
				ctx.SendChain(message.Text("出错了联系管理员看看吧"))
				return
			}
			jData := gjson.Get(binutils.BytesToString(rspData), "data.data")
			if len(jData.String()) == 0 {
				ctx.SendChain(message.Text("没有查到本账号的奇遇呢"))
				return
			}
			for idx, data := range jData.Array() {
				if idx == 0 {
					msg += server + "\n"
				}
				msg = msg + name + "  " + data.Get("serendipity").String() + "  " + data.Get("date_str").String() + "\n"
			}
			ctx.SendChain(message.Text(msg))
			// var msg string
			// commandPart := util.SplitSpace(ctx.State["args"].(string))
			// if len(commandPart) != 2 {
			//	ctx.SendChain(message.Text("参数输入有误！\n" + "奇遇 唯我独尊 柳连柳奶"))
			//	return
			//}
			// server := commandPart[0]
			// name := commandPart[1]
			// qiyuUrl := fmt.Sprintf("https://www.jx3mm.com/home/qyinfo?S=%s&n=%s&u=不限&t=&token=%s", server, name, config.Cfg.MMToken)
			// rspData, err := util.SendHttp(qiyuUrl, []byte(""))
			////rspData, err := web.RequestDataWith(NewTimeOutDefaultClient(), qiyuUrl, "GET", "", web.RandUA())
			////log.Errorln(qiyuUrl, string(rspData), "err", err)
			// if err != nil || gjson.Get(binutils.BytesToString(rspData), "code").Int() != 200 {
			//	ctx.SendChain(message.Text("出错了联系管理员看看吧"))
			//	return
			//}
			// jData := gjson.Get(binutils.BytesToString(rspData), "result")
			// if len(jData.Array()) == 0 {
			//	ctx.SendChain(message.Text("没有查到本账号的奇遇呢"))
			//	return
			//}
			// for idx, data := range jData.Array() {
			//	if idx == 0 {
			//		msg += server + "\n"
			//	}
			//	msg = msg + name + "  " + data.Get("serendipity").String() + "  " + carbon.CreateFromTimestamp(data.Get("time").Int()).ToDateTimeString() + "\n"
			//}
			// ctx.SendChain(message.Text(msg))
		})
	en.OnPrefixGroup([]string{"物价"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			wujia(ctx, datapath, 0)
		})
	en.OnFullMatch("更新骚话", zero.SuperUserPermission).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			err := updateTalk()
			if err != nil {
				ctx.SendChain(message.Text("更新失败", err))
				return
			}
			ctx.SendChain(message.Text("更新成功"))
		})
	en.OnPrefixGroup([]string{"属性"}, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(
		func(ctx *zero.Ctx) {
			attributes(ctx, datapath)
		},
	)
	en.OnPrefixGroup([]string{"战绩"}, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(
		func(ctx *zero.Ctx) {
			indicator(ctx, datapath)
		},
	)
}

func server(ctx *zero.Ctx, server string) {
	if len(serverStatus) != 0 {
		if _, ok := serverIP[server]; ok {
			ctx.SendChain(message.Text("正在尝试Ping ", server, "  ٩(๑´0`๑)۶"))
			process.SleepAbout1sTo2s()
			if !serverStatus[server] {
				ctx.SendChain(message.Text(server, " 垃圾服务器又在维护中  w(ﾟДﾟ)w~"))
				return
			}
			ctx.SendChain(message.Text(server, " 已经开服啦ヽ(✿ﾟ▽ﾟ)ノ~"))
		}
	}
}

func daily(ctx *zero.Ctx, server string) {
	daily4Db := jdb.findDaily(server)
	if carbon.Now().Timestamp()-daily4Db.Time >= 86400 {
		var msg string
		msg += "今天是：" + carbon.Now().ToDateString() + " " + getWeek() + "\n"
		riURL := fmt.Sprintf("https://team.api.jx3box.com/xoyo/daily/task?date=%d", carbon.Now().Timestamp())
		daily, err := web.RequestDataWith(NewTimeOutDefaultClient(), riURL, "GET", "", web.RandUA(), nil)
		if err != nil || gjson.Get(binutils.BytesToString(daily), "code").Int() != 0 {
			ctx.SendChain(util.HTTPError()...)
			return
		}
		for _, d := range gjson.Get(binutils.BytesToString(daily), "data").Array() {
			msg += d.Get("taskType").String() + "：" + d.Get("activityName").String() + "\n"
		}
		for k := range tuiKey {
			tuilanData := tuilan(k)
			questName := gjson.Get(tuilanData, "data.quest_name").String()
			if len(tuilanData) == 0 || k == "大战" || len(questName) == 0 || k == "阵营日常" { // 大战美人图获取jxbox
				continue
			}
			msg += k + "：" + questName + "\n"
		}
		meiURL := fmt.Sprintf("https://spider.jx3box.com/meirentu?server=%s", goUrl.QueryEscape(server))
		meiData, err := web.RequestDataWith(NewTimeOutDefaultClient(), meiURL, "GET", "", web.RandUA(), nil)
		if err != nil || gjson.Get(binutils.BytesToString(meiData), "code").Int() != 0 {
			msg += "美人图：今天没有美人图呢~\n"
		} else {
			msg += "美人图：" + gjson.Get(binutils.BytesToString(meiData), "data.name").String() + "\n"
		}
		msg += fmt.Sprintf("今日活动：%s\n", util.PrettyPrint(date[carbon.Now().Week()]))
		msg += "--------------------------------\n"
		msg += "数据来源JXBOX和推栏"
		ctx.SendChain(message.Text(msg))
		jdb.Insert(&Daily{ //nolint:errcheck
			Server:    server,
			DailyTask: msg,
			Time:      carbon.CreateFromTime(7, 0, 0).Timestamp(),
		})
	} else {
		ctx.SendChain(message.Text(daily4Db.DailyTask))
	}
}

func getWeek() string {
	s := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	intWeek := carbon.Now().Week()
	return s[intWeek]
}

func jinjia(ctx *zero.Ctx, datapath string) {
	var lineStruct []jinPrice
	commandPart := util.SplitSpace(ctx.State["args"].(string))
	var rsp string
	if len(commandPart) != 1 {
		ctx.SendChain(message.Text("参数输入有误！\n" + "金价 绝代天骄"))
		return
	}
	server := commandPart[0]
	if val, ok := allServer[server]; ok {
		data, err := web.RequestDataWith(NewTimeOutDefaultClient(), "https://spider.jx3box.com/jx3price", "GET", "application/x-www-form-urlencoded", web.RandUA(), nil)
		strData := binutils.BytesToString(data)
		if err != nil || gjson.Get(strData, "code").Int() != 0 {
			ctx.SendChain(util.HTTPError()...)
			return
		}
		jin := gjson.Get(strData, fmt.Sprintf("data.%s", val[0]))
		rsp += fmt.Sprintf("今日%s平均金价为：\n", val[0])
		rsp += "5173：" + average(jin.Get("today.5173")) + "￥\n"
		rsp += "万宝楼：" + average(jin.Get("today.official")) + "￥\n"
		rsp += "贴吧：" + average(jin.Get("today.post")) + "￥\n"
		rsp += "------------------------------------------\n"
		json.Unmarshal([]byte(jin.Get("trend").String()), &lineStruct) //nolint:errcheck
		sort.Slice(lineStruct, func(i, j int) bool {
			dateA := strings.Split(lineStruct[i].Date, "-")
			dateB := strings.Split(lineStruct[j].Date, "-")
			for k := 0; k < len(dateA); k++ {
				switch strings.Compare(dateA[k], dateB[k]) {
				case 1:
					return false
				case -1:
					return true
				default:
					continue
				}
			}
			return true
		})
		priceHTML := util.Template2html("goldprice.html", map[string]interface{}{
			"server": val[0],
			"data":   lo.Reverse(lineStruct),
		})
		//html := jibPrice2line(lineStruct, datapath)
		finName, err := util.HTML2pic(datapath, server+util.TodayFileName(), priceHTML)
		if err != nil {
			ctx.SendChain(message.Text("Err:", err))
			return
		}
		ctx.SendChain(message.Text(rsp), message.Image("file:///"+finName))
	} else {
		ctx.SendChain(message.Text("没有找到这个服呢，你是不是乱输的哦~"))
		return
	}
}

func jibPrice2line(lineStruct []jinPrice, datapath string) string {
	var xdata, officialdata, postdata, p5173 []string //nolint:prealloc
	for _, data := range lineStruct {
		xdata = append(xdata, data.Date)
		officialdata = append(officialdata, fmt.Sprintf("%.2f", data.Official))
		postdata = append(postdata, fmt.Sprintf("%.2f", data.Post))
		p5173 = append(p5173, fmt.Sprintf("%.2f", data.P5173))
	}
	page := components.NewPage()
	page.AddCharts(
		drawJinLine("日期", "金价", xdata, map[string][]string{"万宝楼": officialdata,
			"贴吧":   postdata,
			"5173": p5173}),
	)
	f, err := os.Create(datapath + "line.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	page.Render(io.MultiWriter(f)) //nolint:errcheck
	html, _ := os.ReadFile(datapath + "line.html")
	return binutils.BytesToString(html)
}

func drawJinLine(xName, yName string, xdata []string, data map[string][]string) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithLegendOpts(opts.Legend{Show: true, Bottom: "1px"}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: yName,
			SplitLine: &opts.SplitLine{
				Show: false,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: xName,
			//AxisLabel: &opts.AxisLabel{
			//	Interval: "0",
			// },
		}),
	)
	xLine := line.SetXAxis(xdata)
	for name, val := range data {
		xLine = xLine.AddSeries(name, generateLineData(val))
	}
	xLine.SetSeriesOptions(
		charts.WithMarkLineNameTypeItemOpts(opts.MarkLineNameTypeItem{
			Name: "Average",
			Type: "average",
		}),
		charts.WithLineChartOpts(opts.LineChart{
			Smooth: true,
		}),
		charts.WithMarkPointStyleOpts(opts.MarkPointStyle{
			Label: &opts.Label{
				Show:      true,
				Formatter: "{a}: {b}",
			},
		}),
	)
	return line
}

func wujia(ctx *zero.Ctx, datapath string, control int8) {
	if control >= 2 {
		return
	}
	var m sync.Mutex
	m.Lock()
	defer m.Unlock()
	var price = make(map[string][]map[string]interface{})
	var data xiaohei
	commandPart := util.SplitSpace(ctx.State["args"].(string))
	if len(commandPart) != 1 {
		ctx.SendChain(message.Text("参数输入有误！\n" + "物价 牛金"))
		return
	}
	name := commandPart[0]
	if hei, ok := controlCd[name]; ok && (carbon.Now().Timestamp()-hei.last) < 18000 {
		ctx.SendChain(message.Image(hei.fileName))
	} else {
		goodURL := fmt.Sprintf("https://www.j3price.top:8088/black-api/api/outward?name=%s", goUrl.QueryEscape(name))
		rspData, err := web.RequestDataWith(NewTimeOutDefaultClient(), goodURL, "GET", "", web.RandUA(), nil)
		if err != nil || gjson.Get(binutils.BytesToString(rspData), "state").Int() != 0 {
			ctx.SendChain(message.Text("出错了联系管理员看看吧", err, gjson.Get(binutils.BytesToString(rspData), "state").Int()))
			return
		}
		if len(gjson.Get(binutils.BytesToString(rspData), "data").Array()) == 0 { // 如果输入无数据则请求
			searchURL := fmt.Sprintf("https://www.j3price.top:8088/black-api/api/outward/search?step=0&page=1&size=20&name=%s", goUrl.QueryEscape(name))
			searchData, err := web.PostData(searchURL, "application/x-www-form-urlencoded", nil)
			searchList := gjson.Get(binutils.BytesToString(searchData), "data.list").Array()
			if err != nil || len(searchList) == 0 {
				ctx.SendChain(message.Text(fmt.Sprintf("没有找到%s，你是不是乱输的哦~", name)))
				return
			}
			msg := "你可能找的是以下结果：\n"
			for _, s := range searchList {
				msg += s.Get("outwardName").String() + "\n"
			}
			msg += "自动帮你查询：" + searchList[0].Get("outwardAlias").String()
			ctx.SendChain(message.Text(msg))
			ctx.State["args"] = searchList[0].Get("outwardName").String()
			wujia(ctx, datapath, control+1)
			return
		}
		goodid := gjson.Get(binutils.BytesToString(rspData), "data.0.id").Int() // 获得商品id
		infoURL := fmt.Sprintf("https://www.j3price.top:8088/black-api/api/common/search/index/prices?regionId=1&outwardId=%d", goodid)
		wuJiaData, err := web.PostData(infoURL, "application/x-www-form-urlencoded", nil)
		json.Unmarshal(wuJiaData, &data) //nolint:errcheck
		if err != nil || data.State != 0 {
			ctx.SendChain(message.Text("出错了联系管理员看看吧"))
			return
		}
		wujiaPicURL := fmt.Sprintf("https://www.j3price.top:8088/black-api/api/common/search/index/outward?regionId=1&imageLimit=1&outwardId=%d", goodid)
		wujiaPic, err := RequestDataWith(wujiaPicURL)
		if err != nil {
			ctx.SendChain(message.Text("Err:", err))
			return
		}
		for _, rprice := range data.Data.Other {
			if server, ok := xiaoheiIndx[rprice.Prices.Region]; ok {
				price[server] = append(price[server], map[string]interface{}{
					"date":   rprice.Prices.TradeTime,
					"server": rprice.Prices.Server,
					"value":  fmt.Sprintf("%.2f", rprice.Prices.Price),
					"sale":   rprice.Prices.SaleCode,
				})
			}
		}
		for _, rprice := range data.Data.Prices {
			if server, ok := xiaoheiIndx[rprice.Region]; ok {
				price[server] = append(price[server], map[string]interface{}{
					"date":   rprice.TradeTime,
					"server": rprice.Server,
					"value":  fmt.Sprintf("%.2f", rprice.Price),
					"sale":   rprice.SaleCode,
				})
			}
		}
		lineHTML := priceData2line(price, datapath)
		html := util.Template2html("price.html", map[string]interface{}{
			"image": gjson.Get(binutils.BytesToString(wujiaPic), "data.images.0.image"),
			"name":  name,
			"data":  price,
		})
		finName, err := util.HTML2pic(datapath, name+util.TodayFileName(), html+lineHTML)
		if err != nil {
			ctx.SendChain(message.Text("Err:", err))
			return
		}
		controlCd[name] = cd{
			last:     carbon.Now().Timestamp(),
			fileName: "file:///" + finName,
		}
		ctx.SendChain(message.Image("file:///" + finName))
	}
}

func updateTalk() error {
	url := "https://cms.jx3box.com/api/cms/post/jokes?per=%d&page=%d"
	var page int64 = 1
	per := 30
	var Mutex sync.Mutex
	Mutex.Lock()
	defer Mutex.Unlock()
	for {
		data, err := web.GetData(fmt.Sprintf(url, per, page))
		jsonData := binutils.BytesToString(data)
		if err != nil {
			return err
		}
		for _, talkData := range gjson.Get(jsonData, "data.list").Array() {
			jdb.Insert(&Jokes{ //nolint:errcheck
				ID:   talkData.Get("id").Int(),
				Talk: talkData.Get("content").String(),
			})
		}
		if page >= gjson.Get(jsonData, "data.pages").Int() {
			return nil
		}
		page++
		time.Sleep(time.Millisecond * 500)
	}
}

func indicator(ctx *zero.Ctx, datapath string) {
	commandPart := util.SplitSpace(ctx.State["args"].(string))
	var server string
	var name string
	switch {
	case len(commandPart) == 1:
		server = jdb.bind(ctx.Event.GroupID)
		name = commandPart[0]
		if len(server) == 0 {
			ctx.SendChain(message.Text("本群尚未绑定区服"))
			return
		}
	case len(commandPart) == 2:
		server = commandPart[0]
		name = commandPart[1]
	default:
		ctx.SendChain(message.Text("参数输入有误！\n" + "战绩 绝代天骄 xxx"))
		return
	}
	if normServer, ok := allServer[server]; ok {
		zone := normServer[1]
		server = normServer[0]
		var user User
		err := jdb.Find("id = ?", &user, name+"_"+chatServer[server])
		gameRoleID := gjson.Parse(user.Data).Get("body.msg.0.sRoleId").String()
		if err != nil {
			ctx.SendChain(message.Text("没有查找到这个角色呢,试着在世界频道说句话试试吧~"))
			return
		}
		var data = make(map[string]interface{})
		indicator, err := getIndicator(struct {
			RoleID string `json:"role_id"`
			Server string `json:"server"`
			Zone   string `json:"zone"`
			TS     string `json:"ts"`
		}{
			RoleID: gameRoleID,
			Server: server,
			Zone:   zone,
			TS:     ts(),
		})
		if err != nil {
			ctx.SendChain(message.Text("请求剑网推栏失败,请稍后重试~"))
			return
		}
		strIndicator := binutils.BytesToString(indicator)
		templateData := map[string]interface{}{
			"name":   gjson.Get(strIndicator, "data.role_info.name").String(),
			"server": gjson.Get(strIndicator, "data.role_info.zone").String() + "_" + gjson.Get(strIndicator, "data.role_info.server").String(),
			"data":   data,
		}
		performanceData := make(map[string]interface{})
		for _, indicatorData := range gjson.Get(strIndicator, "data.indicator").Array() {
			t := indicatorData.Get("type").String()
			var key string
			performance := indicatorData.Get("performance").IsObject()
			if !performance {
				continue
			}
			switch t {
			case "2c":
				key = "pvp2"
			case "3c":
				key = "pvp3"
			case "5c":
				key = "pvp5"
			}
			performanceData[key] = map[string]string{
				"totalCount": indicatorData.Get("performance.total_count").String(),
				"mvpCount":   indicatorData.Get("performance.mvp_count").String(),
				"winCount":   indicatorData.Get("performance.win_count").String(),
				"mmr":        indicatorData.Get("performance.mmr").String(),
				"ranking":    indicatorData.Get("performance.ranking").String(),
				"winRate":    fmt.Sprintf("%.2f", indicatorData.Get("performance.win_count").Float()/indicatorData.Get("performance.total_count").Float()*100),
				"grade":      indicatorData.Get("performance.grade").String(),
			}
		}
		data["performance"] = performanceData
		history, err := getPersonHistory(struct {
			TS       string `json:"ts"`
			PersonID string `json:"person_id"`
			Cursor   int    `json:"cursor"`
			Size     int    `json:"size"`
		}{
			TS:       ts(),
			PersonID: gjson.Parse(user.Data).Get("body.msg.0.sPersonId").String(),
			Size:     10,
			Cursor:   0,
		})
		if err != nil {
			ctx.SendChain(message.Text("请求剑网推栏失败,请稍后重试~"))
			return
		}
		historyStr := binutils.BytesToString(history)
		for idx, historyData := range gjson.Parse(historyStr).Get("data").Array() {
			startTime := historyData.Get("start_time").Int()
			endTime := historyData.Get("end_time").Int()
			historyStr, _ = sjson.Set(historyStr, "data."+fmt.Sprintf("%d", idx)+".time",
				util.DiffTime(startTime, endTime))
			historyStr, _ = sjson.Set(historyStr, "data."+fmt.Sprintf("%d", idx)+".ago", carbon.CreateFromTimestamp(endTime).ToDateTimeString())
		}
		data["history"] = util.JSONToMap(historyStr)
		templateData["data"] = data
		html := util.Template2html("match.html", templateData)
		finName, err := util.HTML2pic(datapath, name+"_match", html)
		if err != nil {
			ctx.SendChain(message.Text("Err:", err))
			return
		}
		ctx.SendChain(message.Image("file:///" + finName))
	} else {
		ctx.SendChain(message.Text("输入区服有误，请检查qaq~"))
	}
}

func getIndicator(body interface{}) ([]byte, error) {
	xSk := sign(body)
	client := resty.New()
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Host", "m.pvp.xoyo.com").
		SetHeader("Connection", "keep-alive").
		SetHeader("Accept", "application/json").
		SetHeader("fromsys", "APP").
		SetHeader("gamename", "jx3").
		SetHeader("X-Sk", xSk).
		SetHeader("Accept-Language", "zh-CN,zh-Hans;q=0.9").
		SetHeader("apiversion", "3").
		SetHeader("platform", "ios").
		SetHeader("token", (*config.Cfg.JxChat)[0].Token).
		SetHeader("deviceid", "jzrjvE6MDwUbMQTIFIiDQg==").
		SetHeader("Cache-Control", "no-cache").
		SetHeader("clientkey", "1").
		SetHeader("User-Agent", "SeasunGame/193 CFNetwork/1385 Darwin/22.0.0").
		SetHeader("sign", "true").
		SetBody(body).
		Post("https://m.pvp.xoyo.com/role/indicator")
	return res.Body(), err
}

func getPersonHistory(body interface{}) ([]byte, error) {
	xSk := sign(body)
	client := resty.New()
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Host", "m.pvp.xoyo.com").
		SetHeader("Connection", "keep-alive").
		SetHeader("Accept", "application/json").
		SetHeader("fromsys", "APP").
		SetHeader("gamename", "jx3").
		SetHeader("X-Sk", xSk).
		SetHeader("Accept-Language", "zh-CN,zh-Hans;q=0.9").
		SetHeader("apiversion", "3").
		SetHeader("platform", "ios").
		SetHeader("token", (*config.Cfg.JxChat)[0].Token).
		SetHeader("deviceid", "jzrjvE6MDwUbMQTIFIiDQg==").
		SetHeader("Cache-Control", "no-cache").
		SetHeader("clientkey", "1").
		SetHeader("User-Agent", "SeasunGame/193 CFNetwork/1385 Darwin/22.0.0").
		SetHeader("sign", "true").
		SetBody(body).
		Post("https://m.pvp.xoyo.com/mine/match/person-history")
	return res.Body(), err
}

func attributes(ctx *zero.Ctx, datapath string) {
	ts := ts()
	commandPart := util.SplitSpace(ctx.State["args"].(string))
	var server string
	var name string
	switch {
	case len(commandPart) == 1:
		server = jdb.bind(ctx.Event.GroupID)
		name = commandPart[0]
		if len(server) == 0 {
			ctx.SendChain(message.Text("本群尚未绑定区服"))
			return
		}
	case len(commandPart) == 2:
		server = commandPart[0]
		name = commandPart[1]
	default:
		ctx.SendChain(message.Text("参数输入有误！\n" + "属性 绝代天骄 xxx"))
		return
	}
	if normServer, ok := allServer[server]; ok {
		var user User
		zone := normServer[1]
		server = normServer[0]
		err := jdb.Find("id = ?", &user, name+"_"+chatServer[server])
		if err != nil {
			ctx.SendChain(message.Text("没有查找到这个角色呢,试着在世界频道说句话试试吧~"))
			return
		}
		gameRoleID := gjson.Parse(user.Data).Get("body.msg.0.sRoleId").String()
		body := map[string]string{
			"server":       server,
			"zone":         zone,
			"game_role_id": gameRoleID,
			"ts":           ts,
		}
		xSk := sign(body)
		client := resty.New()
		data, err := client.R().
			SetHeader("Content-Type", "application/json").
			// SetHeader("Host", "m.pvp.xoyo.com").
			SetHeader("Connection", "keep-alive").
			SetHeader("Accept", "application/json").
			SetHeader("fromsys", "APP").
			SetHeader("gamename", "jx3").
			SetHeader("X-Sk", xSk).
			SetHeader("Accept-Language", "zh-CN,zh-Hans;q=0.9").
			SetHeader("apiversion", "3").
			SetHeader("platform", "ios").
			SetHeader("token", (*config.Cfg.JxChat)[0].Token).
			SetHeader("deviceid", "jzrjvE6MDwUbMQTIFIiDQg==").
			SetHeader("Cache-Control", "no-cache").
			SetHeader("clientkey", "1").
			SetHeader("User-Agent", "SeasunGame/193 CFNetwork/1385 Darwin/22.0.0").
			SetHeader("sign", "true").
			SetHeader("proxy", "https://m.pvp.xoyo.com/mine/equip/get-role-equip").
			SetBody(body).
			Post("https://http-go-http-proxy-jvuuzynfbg.cn-hangzhou.fcapp.run")
		if err != nil {
			ctx.SendChain(message.Text("请求出错了，稍后试试吧~", err))
			return
		}
		jsonObj := gjson.ParseBytes(data.Body()).Get("data").String()
		templateData := map[string]interface{}{
			"name":   name,
			"server": zone + "_" + server,
			"data":   util.JSONToMap(jsonObj)}
		html := util.Template2html("equip.html", templateData)
		finName, err := util.HTML2pic(datapath, name, html)
		if err != nil {
			ctx.SendChain(message.Text("Err:", err))
		}
		ctx.SendChain(message.Image("file:///" + finName))
	} else {
		ctx.SendChain(message.Text("输入区服有误，请检查qaq~"))
	}
}

func priceData2line(price map[string][]map[string]interface{}, datapath string) string {
	var tmp []map[string]interface{}
	for _, val := range price {
		tmp = append(tmp, val...)
	}
	x, y := make([]string, 0, len(tmp)), make([]string, 0, len(tmp))
	sort.Slice(tmp, func(i, j int) bool {
		dateA := strings.Split(util.Interface2String(tmp[i]["date"]), "/")
		dateB := strings.Split(util.Interface2String(tmp[j]["date"]), "/")
		for k := 0; k < len(dateA); k++ {
			switch strings.Compare(dateA[k], dateB[k]) {
			case 1:
				return false
			case -1:
				return true
			default:
				continue
			}
		}
		return true
	})
	for _, d := range tmp {
		x = append(x, util.Interface2String(d["date"]))
		y = append(y, util.Interface2String(d["value"]))
	}
	page := components.NewPage()
	page.AddCharts(
		drawLine("日期", "价格", x, y),
	)
	f, err := os.Create(datapath + "line.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	page.Render(io.MultiWriter(f)) //nolint:errcheck
	html, _ := os.ReadFile(datapath + "line.html")
	return binutils.BytesToString(html)
}

func drawLine(xName, yName string, x, data []string) *charts.Line {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithYAxisOpts(opts.YAxis{
			Name: yName, // 纵坐标
			SplitLine: &opts.SplitLine{
				Show: false,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: xName, // 横坐标
		}),
		charts.WithLegendOpts(opts.Legend{Show: true, Bottom: "1px"}),
	)
	line.SetXAxis(x).
		AddSeries("价格", generateLineData(data),
			charts.WithLabelOpts(opts.Label{Show: true, Position: "top"})).
		SetSeriesOptions(
			charts.WithMarkLineNameTypeItemOpts(opts.MarkLineNameTypeItem{
				Name: "Average",
				Type: "average",
			}),
			charts.WithLineChartOpts(opts.LineChart{
				Smooth: true,
			}),
			charts.WithMarkPointStyleOpts(opts.MarkPointStyle{
				Label: &opts.Label{
					Show:      true,
					Formatter: "{a}: {b}",
				},
			}),
		)
	return line
}

func generateLineData(data []string) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(data); i++ {
		items = append(items, opts.LineData{Value: data[i]})
	}
	return items
}

func decorator(f func(ctx *zero.Ctx, server string)) func(ctx *zero.Ctx) {
	return func(ctx *zero.Ctx) {
		server := jdb.bind(ctx.Event.GroupID)
		if len(server) != 0 {
			f(ctx, server)
			return
		}
		ctx.SendChain(message.Text("本群还没绑定区服呢"))
	}
}

var serverStatus = make(map[string]bool)

func checkServer(ctx *zero.Ctx, grpList []GroupList) {
	data, err := web.GetData("https://spider2.jx3box.com/api/spider/server/server_state")
	if err != nil {
		log.Errorln("jx3servers error", err)
		return
	}
	serverData := gjson.ParseBytes(data)
	lenServer := len(serverStatus)
	type status struct {
		serverStatus bool
		dbStatus     bool
	}
	var statusList = make(map[string]*status)
	for key := range serverIP {
		statusList[key] = &status{
			serverStatus: true,
			dbStatus:     true,
		}
		if !serverData.Get(fmt.Sprintf(`#(server_name="%s").connect_state`, key)).Bool() {
			statusList[key] = &status{serverStatus: false, dbStatus: serverStatus[key]}
			serverStatus[key] = false
			continue
		}
		statusList[key].dbStatus = serverStatus[key]
		serverStatus[key] = true
	}
	if lenServer != 0 {
		for _, grpListData := range grpList {
			server := grpListData.server
			if _, ok := serverIP[server]; ok {
				if s, ok := statusList[server]; ok {
					msg := server + " 开服啦ヽ(✿ﾟ▽ﾟ)ノ~"
					if s.dbStatus != s.serverStatus {
						if !s.serverStatus {
							msg = server + " 垃圾服务器维护啦  w(ﾟДﾟ)w~"
						}
						log.Errorln("debug server", grpList, statusList[server])
						ctx.SendGroupMessage(grpListData.grp, message.Text(msg))
						process.SleepAbout1sTo2s()
					}
				}
			}
		}
	}
}

func news(ctx *zero.Ctx, grpList []GroupList) {
	msg := make([]News, 0, 5)
	count, _ := jdb.Count(&News{})
	doc, _ := htmlquery.LoadURL("https://jx3.xoyo.com/allnews/")
	li := htmlquery.Find(doc, "/html/body/div[5]/div/div/div[2]/div/div[3]/div[2]/div/div/ul/li")
	for _, node := range li {
		date := htmlquery.InnerText(htmlquery.FindOne(node, "/em"))
		attribute := htmlquery.FindOne(node, "/a")
		title := htmlquery.SelectAttr(attribute, "title")
		href := htmlquery.SelectAttr(attribute, "href")
		kind := htmlquery.InnerText(htmlquery.FindOne(attribute, "/div"))
		if !strings.Contains(href, "https://jx3.xoyo.com") {
			href = "https://jx3.xoyo.com" + href
		}
		canFind := jdb.CanFind("id = ?", &News{}, href)
		data := News{
			ID:    href,
			Date:  date,
			Title: title,
			Kind:  kind,
		}
		if canFind {
			continue
		}
		err := jdb.Insert(&data)
		if err != nil {
			continue
		}
		msg = append(msg, data)
	}
	if count == 0 {
		return
	}
	for _, grpListData := range grpList {
		for _, data := range msg {
			ctx.SendGroupMessage(grpListData.grp, fmt.Sprintf("有新的资讯请查收:\n%s\n%s\n%s\n%s", data.Kind, data.Title, data.ID, data.Date))
			process.SleepAbout1sTo2s()
		}
	}
}

func tuilan(tuiType string) string {
	id, ok := tuiKey[tuiType]
	if ok {
		body := struct {
			ID string `json:"id"`
			TS string `json:"ts"`
		}{
			ID: id,
			TS: ts(),
		}
		xSk := sign(body)
		client := resty.New()
		res, _ := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Host", "m.pvp.xoyo.com").
			SetHeader("Connection", "keep-alive").
			SetHeader("Accept", "application/json").
			SetHeader("fromsys", "APP").
			SetHeader("gamename", "jx3").
			SetHeader("X-Sk", xSk).
			SetHeader("Accept-Language", "zh-CN,zh-Hans;q=0.9").
			SetHeader("apiversion", "3").
			SetHeader("platform", "ios").
			SetHeader("token", (*config.Cfg.JxChat)[0].Token).
			SetHeader("deviceid", "jzrjvE6MDwUbMQTIFIiDQg==").
			SetHeader("Cache-Control", "no-cache").
			SetHeader("clientkey", "1").
			SetHeader("User-Agent", "SeasunGame/193 CFNetwork/1385 Darwin/22.0.0").
			SetHeader("sign", "true").
			SetBody(body).
			Post("https://m.pvp.xoyo.com/activitygw/activity/calendar/detail")
		return binutils.BytesToString(res.Body())
	}
	return ""
}

// func drawTeam(teamId int, datapath string, ctx *zero.Ctx) string {
//	var templateTeamData = make(map[string][]map[string]interface{})
//	team := jdb.getTeamInfo(teamId)
//	mSlice := jdb.getMemberInfo(teamId)
//	title := strconv.Itoa(int(team.TeamId)) + "-----" + team.Dungeon
//	for idx, member := range mSlice {
//		double := "单修"
//		if member.Double == 1 {
//			double = "双修"
//		}
//		templateTeamData["team"+strconv.Itoa(idx/5+1)] = append(templateTeamData["team"+strconv.Itoa(idx/5+1)],
//			map[string]interface{}{
//				"qq":     member.MemberQQ,
//				"name":   member.MemberNickName,
//				"double": double,
//				"img":    iconfile + strconv.Itoa(int(member.MentalId)) + ".png",
//			})
//	}
//	data := map[string]interface{}{
//		"title":   title,            //标题
//		"content": team.Comment,     //备注信息
//		"data":    templateTeamData, //团队信息
//		"group":   team.GroupID,
//		"user":    team.LeaderID,
//	}
//	html := util.Template2html("team.html", data)
//	finName, err := util.Html2pic(datapath, strconv.FormatInt(ctx.Event.GroupID, 10)+"_team", html)
//	if err != nil {
//		return ""
//	}
//	return finName
//}

func drawTeam(teamID int) image.Image {
	Fonts, err := gg.LoadFontFace(text.FontFile, 50)
	if err != nil {
		panic(err)
	}
	const W = 1200
	const H = 1200
	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	// 画直线
	for i := 0; i < 1200; {
		dc.SetRGBA(255, 255, 255, 11)
		dc.SetLineWidth(1)
		dc.DrawLine(0, float64(i), 1200, float64(i))
		dc.Stroke()
		i += 200
	}
	// 画直线
	for i := 200; i < 1200; {
		// dc.SetRGBA(255, 255, 255, 11)
		// dc.SetLineWidth(1)
		dc.DrawLine(float64(i), 200, float64(i), 1200)
		dc.Stroke()
		i += 200
	}
	dc.SetFontFace(Fonts)
	// 队伍
	for i := 1; i < 6; i++ {
		dc.DrawString(strconv.Itoa(i)+"队", 40, float64(100+200*i))
	}
	// 标题
	team := jdb.getTeamInfo(teamID)
	title := strconv.Itoa(int(team.TeamID)) + " " + team.Dungeon
	_, th := dc.MeasureString("哈")
	t := 1200/2 - (float64(len([]rune(title))) / 2)
	dc.DrawStringAnchored(title, t, th, 0.5, 0.5)
	dc.DrawStringAnchored(team.Comment, 1200/2-float64(len([]rune(team.Comment)))/2, 3*th, 0.5, 0.5)
	// 团队
	mSlice := jdb.getMemberInfo(teamID)
	dc.LoadFontFace(text.FontFile, 30) //nolint:errcheck
	_, th = dc.MeasureString("哈")
	start := 200
	for idx, m := range mSlice {
		x := float64(start + idx%5*200 + 10)
		y := float64(start+idx/5*200) + th*2
		dc.DrawString(m.MemberNickName, x, y)
		double := "单修"
		if m.Double == 1 {
			double = "双修"
		}
		dc.DrawString(double, x, y+th*2)
		back, _ := gg.LoadImage(iconfile + strconv.Itoa(int(m.MentalID)) + ".png")
		dc.DrawImage(back, int(x), int(y+th*3))
	}
	return dc.Image()
}

func average(price gjson.Result) string {
	var a float64
	price.ForEach(
		func(key, value gjson.Result) bool {
			a += value.Float()
			return true
		})
	return fmt.Sprintf("%.2f", a/price.Get("#").Float())
}

func ts() string {
	return carbon.Now().Layout("20060102150405", carbon.UTC) + util.Interface2String(carbon.Now(carbon.UTC).Millisecond())
}

func sign(data interface{}) string {
	bData, _ := json.Marshal(data)
	CombineData := util.BytesCombine(bData, []byte("@#?.#@"))
	key := []byte(config.Cfg.SignKey)
	h := hmac.New(sha256.New, key)
	h.Write(CombineData)
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

// 51.2345.67.89
func price2hRead(price int64) (readStr string) {
	strPrice := strconv.FormatInt(price, 10)
	l := len(strPrice)
	for idx, str := range strPrice {
		i := l - idx
		readStr += string(str)
		switch {
		case i == 9:
			readStr += "金砖"
		case i == 5:
			readStr += "金"
		case i == 3:
			readStr += "银"
		}
	}
	readStr += "铜"
	return
}
