package jx3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	goUrl "net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/playwright-community/playwright-go"

	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/DanPlayer/timefinder"
	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/FloatTech/zbputils/math"
	"github.com/FloatTech/zbputils/web"
	"github.com/fogleman/gg"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/golang-module/carbon/v2"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	"github.com/FloatTech/ZeroBot-Plugin/util"
)

const (
	url        = "https://www.jx3api.com/app/"
	realizeUrl = "https://www.jx3api.com/realize/"
	cloudUrl   = "https://www.jx3api.com/cloud/"
)

var tuiKey = map[string]string{
	"大战":     "60f211c82d105c0014c5dd7d",
	"武林通鉴秘境": "60f211c82d105c0014c5de01",
	"武林通鉴公共": "60f211c82d105c0014c5dd97",
	"十人团队秘境": "60f211c82d105c0014c5ddcd",
	"阵营日常":   "60f211c82d105c0014c5dd9d",
}

var allServer = map[string]string{
	"斗转星移": "斗转星移",
	"姨妈":   "斗转星移",
	"蝶恋花":  "蝶恋花",
	"龙争虎斗": "龙争虎斗",
	"长安城":  "长安城",
	"幽月轮":  "幽月轮",
	"剑胆琴心": "剑胆琴心",
	"煎蛋":   "剑胆琴心",
	"乾坤一掷": "乾坤一掷",
	"华乾":   "乾坤一掷",
	"唯我独尊": "唯我独尊",
	"唯满侠":  "唯我独尊",
	"梦江南":  "梦江南",
	"双梦":   "梦江南",
	"绝代天骄": "绝代天骄",
	"绝代":   "绝代天骄",
	"破阵子":  "破阵子",
	"念破":   "破阵子",
	"天鹅坪":  "天鹅坪",
	"纵月":   "天鹅坪",
	"飞龙在天": "飞龙在天",
	"大唐万象": "大唐万象",
	"青梅煮酒": "青梅煮酒",
	"共結來緣": "共結來緣",
	"傲血戰意": "傲血戰意",
	"巔峰再起": "巔峰再起",
	"江海雲夢": "江海雲夢",
}

type cd struct {
	last     int64
	fileName string
}

type JinPrice struct {
	P5173    float64 `json:"5173"`
	Post     float64 `json:"post"`
	Official float64 `json:"official"`
	Date     string  `json:"date"`
}

var heiCd = make(map[string]cd)

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
				Id          int    `json:"id"`
				CreatedTime string `json:"createdTime"`
				UpdatedTime string `json:"updatedTime"`
				RegionName  string `json:"regionName"`
				RegionNick  string `json:"regionNick"`
				Charge      string `json:"charge"`
			} `json:"region"`
			Prices struct {
				Id          int         `json:"id"`
				Price       float64     `json:"price"`
				Region      string      `json:"region"`
				RegionAlias string      `json:"regionAlias"`
				RegionId    int         `json:"regionId"`
				Server      string      `json:"server"`
				ServerId    int         `json:"serverId"`
				SaleCode    string      `json:"saleCode"`
				TradeTime   string      `json:"tradeTime"`
				OutwardName interface{} `json:"outwardName"`
				OutwardId   int         `json:"outwardId"`
				Audit       int         `json:"audit"`
				Now         int         `json:"now"`
				Exterior    string      `json:"exterior"`
				Pricer      string      `json:"pricer"`
			} `json:"prices"`
		} `json:"other"`
		Region struct {
			Id          int    `json:"id"`
			CreatedTime string `json:"createdTime"`
			UpdatedTime string `json:"updatedTime"`
			RegionName  string `json:"regionName"`
			RegionNick  string `json:"regionNick"`
			Charge      string `json:"charge"`
		} `json:"region"`
		Prices []struct {
			Id          int         `json:"id"`
			Price       float64     `json:"price"`
			Region      string      `json:"region"`
			RegionAlias string      `json:"regionAlias"`
			RegionId    int         `json:"regionId"`
			Server      string      `json:"server"`
			ServerId    int         `json:"serverId"`
			SaleCode    string      `json:"saleCode"`
			TradeTime   string      `json:"tradeTime"`
			OutwardName interface{} `json:"outwardName"`
			OutwardId   int         `json:"outwardId"`
			Audit       int         `json:"audit"`
			Now         int         `json:"now"`
			Exterior    string      `json:"exterior"`
			Pricer      string      `json:"pricer"`
		} `json:"prices"`
	} `json:"data"`
	Message interface{} `json:"message"`
}

func init() {
	go startWs()
	en := control.Register("jx", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		PrivateDataFolder: "jx3",
		Help: "- 日常任务xxx(eg 日常任务绝代天骄)\n" +
			"- 开服检查xxx(eg 开服检查绝代天骄)\n" +
			"- 金价查询xxx(eg 金价查询绝代天骄)\n" +
			"- 花价|花价查询 xxx xxx xxx(eg 花价 绝代天骄 绣球花 广陵邑)\n" +
			"- 小药\n" +
			"- xxx配装(eg 分山劲配装)\n" +
			"- xxx奇穴(eg 分山劲奇穴)\n" +
			"- 宏xxx(eg 宏分山劲)\n" +
			"- 沙盘xxx(eg 沙盘绝代天骄)\n" +
			"- 装饰属性|装饰xxx(eg 装饰混沌此生)\n" +
			"- 奇遇条件xxx(eg 奇遇条件三山四海)\n" +
			"- 奇遇攻略xxx(eg 奇遇攻略三山四海)\n" +
			"- 维护公告\n" +
			"- 骚话（不区分大小写）\n" +
			"- 舔狗\n" +
			"-（开启|关闭）jx推送\n" +
			"- /roll随机roll点\n" +
			"- 物价xxx\n" +
			"- 团队相关见 https://docs.qq.com/doc/DUGJRQXd1bE5YckhB",
	})
	go func() {
		initialize()
	}()
	datapath := file.BOTPATH + "/" + en.DataFolder()
	en.OnFullMatchGroup([]string{"日常", "日常任务"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			decorator(daily)(ctx)
		})
	en.OnRegex(`^开服检查(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			data := map[string]string{"server": strings.Replace(str, " ", "", -1)}
			reqbody, err := json.Marshal(data)
			rsp, err := util.SendHttp(url+"check", reqbody)
			if err != nil {
				log.Errorln("jx3daily:", err)
			}
			json := gjson.ParseBytes(rsp)
			var text []interface{}
			for _, value := range json.Get("data").Array() {
				if value.Get("status").Int() == 1 {
					text = append(text, value.Get("server").Str+"：开服\n")
				} else {
					text = append(text, value.Get("server").Str+"：停服\n")
				}
			}
			ctx.SendChain(message.Text(
				fmt.Sprint(text),
			))
		})
	en.OnPrefixGroup([]string{"金价", "金价查询"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			jinjia(ctx, datapath)
		})
	en.OnRegex(`^(花价|花价查询).*?\s(.*).*?\s(.*).*?\s(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			server := ctx.State["regex_matched"].([]string)[2]
			flower := ctx.State["regex_matched"].([]string)[3]
			homeMap := ctx.State["regex_matched"].([]string)[4]
			if len(server) == 0 {
				ctx.SendChain(message.Text(
					"请输入区服",
				))
			} else {
				data := map[string]string{"server": strings.Replace(server, " ", "", -1), "flower": flower, "map": homeMap}
				reqbody, err := json.Marshal(data)
				rsp, err := util.SendHttp(url+"flower", reqbody)
				if err != nil {
					log.Errorln("jx3daily:", err)
				}
				json := gjson.ParseBytes(rsp)
				text := ""
				for _, value := range json.Get("data").Array() {
					value.ForEach(func(key, v gjson.Result) bool {
						switch key.String() {
						case "name", "color", "price":
							text = text + key.String() + ":" + v.String() + "\n"
						}
						return true
					})
				}
				ctx.SendChain(message.Text(text))
			}
		})
	en.OnPrefixGroup([]string{"沙盘"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) != 1 {
				ctx.SendChain(message.Text("参数输入有误！\n" + "沙盘 绝代天骄"))
				return
			}
			server := commandPart[0]
			if _, ok := allServer[server]; ok {
				client := web.NewDefaultClient()
				request, err := http.NewRequest("GET", fmt.Sprintf("https://www.j3sp.com/api/sand/?serverName=%s&shadow=0&is_history=1", server), nil)
				if err == nil {
					// 增加header选项
					var response *http.Response
					request.Header.Add("Cookie", "spc_token=1e245818-e241-437e-a292-dfa5544a2c9f")
					response, err = client.Do(request)
					if err == nil {
						if response.StatusCode != http.StatusOK {
							ctx.SendChain(message.Text("请求出错了稍后再试试吧~"))
							return
						}
						data, _ := io.ReadAll(response.Body)
						response.Body.Close()
						strData := binary.BytesToString(data)
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
	en.OnRegex(`^(装饰属性|装饰)(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := ctx.State["regex_matched"].([]string)[2]
			data := map[string]string{"name": strings.Replace(name, " ", "", -1)}
			reqbody, err := json.Marshal(data)
			rsp, err := util.SendHttp(url+"furniture", reqbody)
			if err != nil {
				log.Errorln("jx3daily:", err)
			}
			json := gjson.ParseBytes(rsp)
			ctx.SendChain(message.Text(
				"名称：", json.Get("data.name"), "\n",
				"品质：", json.Get("data.quality"), "\n",
				"产出地图：", json.Get("data.source"), "\n",
				"等级限制：", json.Get("data.level_limit"), "\n",
				"品质等级：", json.Get("data.quality_level"), "\n",
				"观赏分数：", json.Get("data.view_score"), "\n",
				"实用分数：", json.Get("data.practical_score"), "\n",
				"tips：", json.Get("data.tip"), "\n",
			), message.Image(
				json.Get("data.image_path").String(),
			))
		})
	en.OnRegex(`^前置(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := ctx.State["regex_matched"].([]string)[1]
			data := map[string]string{"name": strings.Replace(name, " ", "", -1)}
			reqbody, err := json.Marshal(data)
			rsp, err := util.SendHttp(url+"require", reqbody)
			if err != nil {
				log.Errorln("jx3daily:", err)
			}
			json := gjson.ParseBytes(rsp)
			ctx.SendChain(
				message.Text(
					"名称：", json.Get("data.name"), "\n",
					"方法：", json.Get("data.means"), "\n",
					"前置：", json.Get("data.require"), "\n",
					"奖励：", json.Get("data.reward"), "\n",
				),
				message.Image(
					json.Get("data.upload").String()),
			)
		})
	en.OnSuffix("小药").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Image(fileUrl + "medicine.png"))
		})
	en.OnSuffix("配装").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := ctx.State["args"].(string)
			if len(name) == 0 {
				ctx.SendChain(message.Text("请输入职业！！！！"))
			} else {
				data := map[string]string{"name": getMental(strings.Replace(name, " ", "", -1))}
				reqbody, err := json.Marshal(data)
				rsp, err := util.SendHttp(url+"equip", reqbody)
				if err != nil {
					log.Errorln("jx3daily:", err)
				}
				json := gjson.ParseBytes(rsp)
				ctx.SendChain(
					message.Text("PVE：\n"),
					message.Image(
						json.Get("data.pve").String()),
					message.Text("\nPVP：\n"),
					message.Image(
						json.Get("data.pvp").String()),
				)
			}
		})
	en.OnSuffix("奇穴").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := ctx.State["args"].(string)
			if len(name) == 0 {
				ctx.SendChain(message.Text("请输入职业！！！！"))
			} else {
				data := map[string]string{"name": getMental(strings.Replace(name, " ", "", -1))}
				reqbody, err := json.Marshal(data)
				rsp, err := util.SendHttp(url+"qixue", reqbody)
				if err != nil {
					log.Errorln("jx3daily:", err)
				}
				json := gjson.ParseBytes(rsp)
				ctx.SendChain(
					message.Text("通用：\n"),
					message.Image(
						json.Get("data.all").String()),
					message.Text("\n吃鸡：\n"),
					message.Image(
						json.Get("data.longmen").String()),
				)
			}
		})
	en.OnPrefix("宏").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := ctx.State["args"].(string)
			mental := getMentalData(strings.Replace(name, " ", "", -1))
			mentalUrl := fmt.Sprintf("https://cms.jx3box.com/api/cms/posts?type=macro&per=10&page=1&order=update&client=std&search=%s", goUrl.QueryEscape(mental.Name))
			data, err := web.RequestDataWith(web.NewDefaultClient(), mentalUrl, "GET", "application/x-www-form-urlencoded", web.RandUA())
			DataList := gjson.Get(binary.BytesToString(data), "data.list").Array()
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
				rsp += "数据来源于JXBOX，dps请自行测试"
				ctx.SendChain(message.Text(rsp))
				time.Sleep(time.Second * 4)
			}
		})
	en.OnSuffix("阵眼").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := ctx.State["args"].(string)
			if utf8.RuneCountInString(name) > 5 {
				log.Println("name len max")
			} else {
				data := map[string]string{"name": getMental(strings.Replace(name, " ", "", -1))}
				reqbody, err := json.Marshal(data)
				rsp, err := util.SendHttp(url+"matrix", reqbody)
				if err != nil {
					log.Errorln("jx3daily:", err)
				}
				//	json := gjson.ParseBytes(rsp)
				log.Errorln(string(rsp))
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
	//				dwList := gjson.Get(binary.BytesToString(dwData), "list").Array()
	//				if len(dwList) == 0 {
	//					ctx.SendChain(message.Text(fmt.Sprintf("没有找到%s呢，你是不是乱输的哦~", name)))
	//					return
	//				}
	//				dwId := dwList[0].Get("dwID").String()
	//				json, _ := web.GetData("https://icon.jx3box.com/pvx/serendipity/output/serendipity.json")
	//				articleId := gjson.Get(binary.BytesToString(json), dwId).String()
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
	en.OnPrefix("攻略").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) != 1 {
				ctx.SendChain(message.Text("输入参数有误！！！攻略 炼狱厨神"))
				return
			} else {
				name := commandPart[0]
				dbData := getAdventure(name)
				if len(dbData.Pic) == 0 || carbon.Now().DiffAbsInSeconds(carbon.CreateFromTimestamp(dbData.Time)) > 3600*10 {
					dwData, _ := web.GetData(fmt.Sprintf("https://node.jx3box.com/serendipities?name=%s", goUrl.QueryEscape(name)))
					dwList := gjson.Get(binary.BytesToString(dwData), "list").Array()
					if len(dwList) == 0 {
						ctx.SendChain(message.Text(fmt.Sprintf("没有找到%s呢，你是不是乱输的哦~", name)))
						return
					}
					dwId := dwList[0].Get("dwID").String()
					json, _ := web.GetData("https://icon.jx3box.com/pvx/serendipity/output/serendipity.json")
					articleId := gjson.Get(binary.BytesToString(json), dwId).String()
					articleUrl := fmt.Sprintf("https://www.jx3box.com/cj/#/view/%s", articleId)
					pw, err := playwright.Run()
					if err != nil {
						playwright.Install()
						playwright.Run()
					}
					defer pw.Stop()
					browser, err := pw.Chromium.Launch()
					if err != nil {
						playwright.Install()
					}
					page, err := browser.NewPage(playwright.BrowserNewContextOptions{
						IsMobile: playwright.Bool(true),
					})
					if err != nil {
						return
					}
					_, err = page.Goto(articleUrl, playwright.PageGotoOptions{
						WaitUntil: playwright.WaitUntilStateNetworkidle,
						Timeout:   playwright.Float(10000),
					})
					if err != nil {
						return
					}
					page.Click("//*[@id=\"app\"]/aside/span")
					result, _ := page.QuerySelector("div[class=\"c-article-chunk on\"]")
					html, _ := result.InnerHTML()
					htmlPage, err := browser.NewPage()
					if err != nil {
						return
					}
					err = htmlPage.SetContent(html, playwright.PageSetContentOptions{
						WaitUntil: playwright.WaitUntilStateNetworkidle,
					})
					htmlPage.Keyboard().Down("PageDown")
					time.Sleep(time.Second * 10)
					htmlPage.Keyboard().Up("PageDown")
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
					updateAdventure(db)
					ctx.SendChain(message.ImageBytes(b))
				} else {
					ctx.SendChain(message.ImageBytes(dbData.Pic))
				}
			}
		})
	en.OnRegex(`^(维护公告|更新公告)(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			data := map[string]string{"limit": "3"}
			reqbody, err := json.Marshal(data)
			rsp, err := util.SendHttp(url+"announce", reqbody)
			if err != nil {
				log.Errorln("jx3daily:", err)
			}
			text := ""
			gjson.Get(helper.BytesToString(rsp), "data").ForEach(func(_, value gjson.Result) bool {
				text = text + value.Get("title").String() + "\n" + value.Get("date").String() + "\n" + value.Get("url").String() + "\n"
				return true
			})
			ctx.SendChain(message.Text(text))
		})
	en.OnRegex(`^(?i)骚话(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			rsp, err := util.SendHttp(url+"random", nil)
			if err != nil {
				log.Errorln("jx3daily:", err)
			}
			json := gjson.ParseBytes(rsp)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(json.Get("data.text")))
		})
	en.OnRegex(`^舔狗(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			rsp, err := util.SendHttp(realizeUrl+"random", nil)
			if err != nil {
				log.Errorln("jx3daily:", err)
			}
			json := gjson.ParseBytes(rsp)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(json.Get("data.text")))
		})
	en.OnFullMatch("开启jx推送", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			area := enable(ctx.Event.GroupID)
			if len(area) == 0 {
				var server []string
				rsp, _ := util.SendHttp(url+"check", nil)
				json := gjson.ParseBytes(rsp)
				for _, value := range json.Get("data").Array() {
					server = append(server, value.Get("server").String())
				}
				ctx.Send(message.Text("开启成功，检测到当前未绑定区服，请输入\n绑定区服xxx\n进行绑定，可选服务器有：\n" + fmt.Sprint(server)))
			} else {
				ctx.Send(message.Text("开启成功，当前绑定区服为：" + area))
			}
		})
	en.OnFullMatch("更新内容").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			data := map[string]string{"robot": zero.BotConfig.NickName[0]}
			reqbody, err := json.Marshal(data)
			if err != nil {
				log.Errorln("jx3daily:", err)
			}
			rsp, _ := util.SendHttp(cloudUrl+"content", reqbody)
			json := gjson.ParseBytes(rsp)
			ctx.Send(message.Image(json.Get("data.url").String()))
		})
	en.OnFullMatch("/roll").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(fmt.Sprintf("%s 投出了%d点。", ctx.Event.Sender.NickName, util.Rand(1, 100))), message.Reply(ctx.Event.MessageID))
		})
	en.OnFullMatch("关闭jx推送", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			disable(ctx.Event.GroupID)
			ctx.Send(message.Text("关闭成功"))
		})
	en.OnPrefix("绑定区服", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			area := strings.Replace(ctx.State["args"].(string), " ", "", -1)
			if _, ok := allServer[area]; ok {
				bindArea(ctx.Event.GroupID, area)
				ctx.Send(message.Text("绑定成功"))
			} else {
				ctx.Send(message.Text("区服输入有误"))
			}
		})
	// 开团 时间 副本名 备注
	en.OnPrefixGroup([]string{"开团", "新建团队", "创建团队"}, func(ctx *zero.Ctx) bool {
		return isOk(ctx.Event.UserID)
	}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) != 3 {
				ctx.SendChain(message.Text("开团参数输入有误！"))
				return
			}
			startTime := parseDate(commandPart[0])
			dungeon := commandPart[1]
			comment := commandPart[2]
			leaderId := ctx.Event.UserID
			teamId, err := createNewTeam(startTime, dungeon, comment, leaderId, ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Text("Error :", err))
				return
			}
			ctx.SendChain(message.Text("开团成功，团队id为：", teamId))
		})
	// 报团 团队ID 心法 角色名 [是否双休] 按照报名时间先后默认排序 https://docs.qq.com/doc/DUGJRQXd1bE5YckhB
	en.OnPrefixGroup([]string{"报名", "报团", "报名团队"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			double := 0
			if len(commandPart) == 3 {
				double = 0
			} else if len(commandPart) == 4 {
				double, _ = strconv.Atoi(commandPart[3])
			} else {
				ctx.SendChain(message.Text("报团参数有误。"))
				return
			}
			teamId, _ := strconv.Atoi(commandPart[0])
			mental := getMentalData(commandPart[1])
			nickName := commandPart[2]
			if mental.ID == 0 {
				ctx.SendChain(message.Text("心法输入有误"))
				return
			}
			Team := getTeamInfo(teamId)
			if carbon.Now().Timestamp() >= Team.StartTime || Team.GroupId != ctx.Event.GroupID {
				ctx.SendChain(message.Text("当前团队已过期或团队不存在。"))
				return
			}
			if isInTeam(teamId, ctx.Event.UserID) {
				ctx.SendChain(message.Text("你已经在团队中了。"))
				return
			}
			var member = Member{
				TeamId:         teamId,
				MemberQQ:       ctx.Event.UserID,
				MemberNickName: nickName,
				MentalId:       mental.ID,
				Double:         double,
				SignUp:         carbon.Now().Timestamp(),
			}
			addMember(&member)
			ctx.SendChain(message.Text("报团成功"), message.Reply(ctx.Event.MessageID))
			ctx.SendChain(message.Text("当前团队:\n"), message.Image("base64://"+helper.BytesToString(util.Image2Base64(drawTeam(teamId)))))
		})
	en.OnPrefixGroup([]string{"撤销报团"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			teamId, _ := strconv.Atoi(commandPart[0])
			if !isBelongGroup(teamId, ctx.Event.GroupID) {
				ctx.SendChain(message.Text("参数输入有误。"))
				return
			}
			deleteMember(teamId, ctx.Event.UserID)
			ctx.SendChain(message.Text("撤销成功"), message.Reply(ctx.Event.MessageID))
			ctx.SendChain(message.Text("当前团队:\n"), message.Image("base64://"+helper.BytesToString(util.Image2Base64(drawTeam(teamId)))))
		})
	en.OnFullMatchGroup([]string{"我报的团", "我的报名"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			SignUp := util.RemoveRepByMap(getSignUp(ctx.Event.UserID))
			var InfoTeam []Team
			for _, d := range SignUp {
				Team := getEfficientTeamInfo(
					fmt.Sprintf("WHERE teamID = '%d' AND startTime > '%d' AND groupId = '%d'", d, carbon.Now().Timestamp(), ctx.Event.GroupID))
				if len(Team) > 0 {
					InfoTeam = append(InfoTeam, Team[0])
				}
			}
			out := ""
			for _, data := range InfoTeam {
				out = out + fmt.Sprintf("团队id：%d,团长 ：%d,副本：%s，开始时间：%s，备注：%s\n",
					data.TeamId, data.LeaderId, data.Dungeon, carbon.CreateFromTimestamp(data.StartTime).ToDateTimeString(), data.Comment)
			}
			ctx.SendChain(message.Text(out))
		})
	en.OnFullMatchGroup([]string{"我的开团"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			InfoSlice := getEfficientTeamInfo(
				fmt.Sprintf("WHERE leaderId = '%d' AND startTime > '%d' AND groupId = '%d'", ctx.Event.UserID, carbon.Now().Timestamp(), ctx.Event.GroupID))
			out := ""
			for _, data := range InfoSlice {
				out = out + fmt.Sprintf("团队id：%d,团长 ：%d,副本：%s，开始时间：%s，备注：%s\n",
					data.TeamId, data.LeaderId, data.Dungeon, carbon.CreateFromTimestamp(data.StartTime).ToDateTimeString(), data.Comment)
			}
			ctx.SendChain(message.Text(out))
		})
	en.OnFullMatchGroup([]string{"全团显示"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			InfoSlice := getEfficientTeamInfo(
				fmt.Sprintf("WHERE startTime > '%d' AND groupId = '%d'", carbon.Now().Timestamp(), ctx.Event.GroupID))
			out := ""
			for _, data := range InfoSlice {
				out = out + fmt.Sprintf("团队id：%d,团长 ：%d,副本：%s，开始时间：%s，备注：%s\n",
					data.TeamId, data.LeaderId, data.Dungeon, carbon.CreateFromTimestamp(data.StartTime).ToDateTimeString(), data.Comment)
			}
			ctx.SendChain(message.Text(out))
		})
	en.OnPrefixGroup([]string{"准备进本"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			teamId, err := strconv.Atoi(commandPart[0])
			if err != nil {
				return
			}
			if !isBelongGroup(teamId, ctx.Event.GroupID) {
				ctx.SendChain(message.Text("参数输入有误。"))
				return
			}
			var at message.Message
			mSlice := getMemberInfo(teamId)
			for _, m := range mSlice {
				at = append(at, message.At(m.MemberQQ))
			}
			at = append(at, message.Text("\n准备进本啦！！"))
			ctx.Send(at)
		})
	// 查看团队 teamid
	en.OnPrefixGroup([]string{"查看团队", "查询团队", "查团"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			teamId, _ := strconv.Atoi(commandPart[0])
			if !isBelongGroup(teamId, ctx.Event.GroupID) {
				ctx.SendChain(message.Text("参数输入有误。"))
				return
			}
			ctx.SendChain(message.Image("base64://" + helper.BytesToString(util.Image2Base64(drawTeam(teamId)))))
		})
	// 申请团长 团牌
	en.OnPrefixGroup([]string{"申请团长"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			permission := 0
			var teamName string
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			teamName = ""
			if len(commandPart) == 1 {
				teamName = commandPart[0]
			}
			if ctx.Event.Sender.Role != "member" {
				permission = 1
			}
			err := newLeader(ctx.Event.UserID, ctx.Event.Sender.NickName, permission, teamName)
			if err == 0 {
				ctx.SendChain(message.Text("申请团长成功。"))
			}
			if err == -1 {
				ctx.SendChain(message.Text("贵人多忘事，你已经申请过了"))
			}
		})
	// 取消开团 团队id
	en.OnPrefixGroup([]string{"取消开团", "删除团队", "撤销团队", "撤销开团"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) < 1 {
				ctx.SendChain(message.Text("撤销开团参数有误"))
			}
			teamId, err := strconv.Atoi(commandPart[0])
			if err != nil || !isBelongGroup(teamId, ctx.Event.GroupID) {
				ctx.SendChain(message.Text("团队id输入有误"))
				return
			}
			status := delTeam(teamId, ctx.Event.UserID)
			switch status {
			case -1:
				ctx.SendChain(message.Text("这个团不是你的。无法删除"))
			case 0:
				ctx.SendChain(message.Text("删除成功"))
			}
		})
	// 同意审批@qq
	en.OnRegex(`^同意审批.*?(\d+)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			qq := math.Str2Int64(ctx.State["regex_matched"].([]string)[1])
			teamName := acceptLeader(qq)
			ctx.SendChain(message.At(qq), message.Text("已成为团长,团队名称为："),
				message.Text(teamName))
		})
	en.OnRegex(`^删除团长.*?(\d+)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			qq := math.Str2Int64(ctx.State["regex_matched"].([]string)[1])
			deleteLeader(qq)
			ctx.SendChain(message.At(qq), message.Text("删除成功"))
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
			qiyuUrl := fmt.Sprintf("https://pull.j3cx.com/api/serendipity?server=%s&role=%s&pageIndex=1&pageSize=30", server, name)
			rspData, err := web.RequestDataWith(web.NewDefaultClient(), qiyuUrl, "GET", "", web.RandUA())
			if err != nil || gjson.Get(binary.BytesToString(rspData), "code").Int() != 0 {
				ctx.SendChain(message.Text("出错了联系管理员看看吧"))
				return
			}
			jData := gjson.Get(binary.BytesToString(rspData), "data.data")
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
			////rspData, err := web.RequestDataWith(web.NewDefaultClient(), qiyuUrl, "GET", "", web.RandUA())
			////log.Errorln(qiyuUrl, string(rspData), "err", err)
			// if err != nil || gjson.Get(binary.BytesToString(rspData), "code").Int() != 200 {
			//	ctx.SendChain(message.Text("出错了联系管理员看看吧"))
			//	return
			//}
			// jData := gjson.Get(binary.BytesToString(rspData), "result")
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
}

func daily(ctx *zero.Ctx, server string) {
	var msg string
	msg += "今天是：" + carbon.Now().ToDateString() + " " + util.GetWeek() + "\n"
	riUrl := fmt.Sprintf("https://team.api.jx3box.com/xoyo/daily/task?date=%d", carbon.Now().Timestamp())
	daily, err := web.RequestDataWith(web.NewDefaultClient(), riUrl, "GET", "", web.RandUA())
	if err != nil || gjson.Get(binary.BytesToString(daily), "code").Int() != 0 {
		ctx.SendChain(message.Text("出错了联系管理员看看吧~"))
		return
	}
	for _, d := range gjson.Get(binary.BytesToString(daily), "data").Array() {
		msg += d.Get("taskType").String() + "：" + d.Get("activityName").String() + "\n"
	}
	for k := range tuiKey {
		tuilanData := tuilan(k)
		quest_name := gjson.Get(tuilanData, "data.quest_name").String()
		if len(tuilanData) == 0 || k == "大战" || len(quest_name) == 0 || k == "阵营日常" { // 大战美人图获取jxbox
			continue
		}
		msg += k + "：" + quest_name + "\n"
	}
	meiUrl := fmt.Sprintf("https://spider.jx3box.com/meirentu?server=%s", goUrl.QueryEscape(server))
	meiData, err := web.RequestDataWith(web.NewDefaultClient(), meiUrl, "GET", "", web.RandUA())
	if err != nil || gjson.Get(binary.BytesToString(meiData), "code").Int() != 0 {
		msg += "美人图：今天没有美人图呢~\n"
	} else {
		msg += "美人图：" + gjson.Get(binary.BytesToString(meiData), "data.name").String() + "\n"
	}
	msg += fmt.Sprintf("今日活动：%s\n", util.PrettyPrint(date[carbon.Now().Week()]))
	msg += "--------------------------------\n"
	msg += "数据来源JXBOX和推栏"
	ctx.SendChain(message.Text(msg))
}

func jinjia(ctx *zero.Ctx, datapath string) {
	var lineStruct []JinPrice
	commandPart := util.SplitSpace(ctx.State["args"].(string))
	var rsp string
	if len(commandPart) != 1 {
		ctx.SendChain(message.Text("参数输入有误！\n" + "金价 绝代天骄"))
		return
	}
	server := commandPart[0]
	if val, ok := allServer[server]; ok {
		data, err := web.RequestDataWith(web.NewDefaultClient(), "https://spider.jx3box.com/jx3price", "GET", "application/x-www-form-urlencoded", web.RandUA())
		strData := binary.BytesToString(data)
		if err != nil || gjson.Get(strData, "code").Int() != 0 {
			ctx.SendChain(message.Text("出错了，请稍后再试吧"))
			return
		}
		jin := gjson.Get(strData, fmt.Sprintf("data.%s", val))
		rsp += fmt.Sprintf("今日%s平均金价为：\n", val)
		rsp += "5173：" + average(jin.Get("today.5173")) + "￥\n"
		rsp += "万宝楼：" + average(jin.Get("today.official")) + "￥\n"
		rsp += "贴吧：" + average(jin.Get("today.post")) + "￥\n"
		rsp += "------------------------------------------\n"
		rsp += "数据来源万宝楼\n"
		json.Unmarshal([]byte(jin.Get("trend").String()), &lineStruct)
		html := jibPrice2line(lineStruct, datapath)
		finName, err := util.Html2pic(datapath, server+util.TodayFileName(), "price.html", html)
		ctx.SendChain(message.Text(rsp), message.Image("file:///"+finName))
	} else {
		ctx.SendChain(message.Text("没有找到这个服呢，你是不是乱输的哦~"))
		return
	}
}

func jibPrice2line(lineStruct []JinPrice, datapath string) string {
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
	var xdata, officialdata, postdata, p5173 []string
	for _, data := range lineStruct {
		xdata = append(xdata, data.Date)
		officialdata = append(officialdata, fmt.Sprintf("%.2f", data.Official))
		postdata = append(postdata, fmt.Sprintf("%.2f", data.Post))
		p5173 = append(p5173, fmt.Sprintf("%.2f", data.P5173))
	}
	page := components.NewPage()
	page.AddCharts(
		drawJinLine("日期", "金价", xdata, map[string][]string{"official": officialdata,
			"贴吧":   postdata,
			"5173": p5173}),
	)
	f, err := os.Create(datapath + "line.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	page.Render(io.MultiWriter(f))
	html, _ := ioutil.ReadFile(datapath + "line.html")
	return binary.BytesToString(html)
}

func drawJinLine(XName, YName string, xdata []string, data map[string][]string) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithLegendOpts(opts.Legend{Show: true, Bottom: "1px"}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: YName,
			SplitLine: &opts.SplitLine{
				Show: false,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: XName,
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
	if hei, ok := heiCd[name]; ok && (carbon.Now().Timestamp()-hei.last) < 18000 {
		ctx.SendChain(message.Image(hei.fileName))
	} else {
		goodUrl := fmt.Sprintf("https://www.j3price.top:8088/black-api/api/outward?name=%s", goUrl.QueryEscape(name))
		rspData, err := web.RequestDataWith(web.NewDefaultClient(), goodUrl, "GET", "", web.RandUA())
		if err != nil || gjson.Get(binary.BytesToString(rspData), "state").Int() != 0 {
			ctx.SendChain(message.Text("出错了联系管理员看看吧"))
			return
		}
		if len(gjson.Get(binary.BytesToString(rspData), "data").Array()) == 0 { // 如果输入无数据则请求
			searchUrl := fmt.Sprintf("https://www.j3price.top:8088/black-api/api/outward/search?step=0&page=1&size=20&name=%s", goUrl.QueryEscape(name))
			searchData, err := web.PostData(searchUrl, "application/x-www-form-urlencoded", nil)
			searchList := gjson.Get(binary.BytesToString(searchData), "data.list").Array()
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
		goodid := gjson.Get(binary.BytesToString(rspData), "data.0.id").Int() // 获得商品id
		infoUrl := fmt.Sprintf("https://www.j3price.top:8088/black-api/api/common/search/index/prices?regionId=1&outwardId=%d", goodid)
		wuJiaData, err := web.PostData(infoUrl, "application/x-www-form-urlencoded", nil)
		json.Unmarshal(wuJiaData, &data)
		if err != nil || data.State != 0 {
			ctx.SendChain(message.Text("出错了联系管理员看看吧"))
			return
		}
		wujiaPicUrl := fmt.Sprintf("https://www.j3price.top:8088/black-api/api/common/search/index/outward?regionId=1&imageLimit=1&outwardId=%d", goodid)
		wujiaPic, err := util.RequestDataWith(wujiaPicUrl)
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
		d := map[string]interface{}{
			"image": gjson.Get(binary.BytesToString(wujiaPic), "data.images.0.image"),
			"name":  name,
			"data":  price,
		}
		lineHtml := priceData2line(price, datapath)
		html := util.Template2html("price.html", d)
		finName, err := util.Html2pic(datapath, name+util.TodayFileName(), "price.html", html+lineHtml)
		heiCd[name] = cd{
			last:     carbon.Now().Timestamp(),
			fileName: "file:///" + finName,
		}
		ctx.SendChain(message.Image("file:///" + finName))
	}
}

func priceData2line(price map[string][]map[string]interface{}, datapath string) string {
	var x []string
	var y []string
	var tmp []map[string]interface{}
	for _, val := range price {
		tmp = append(tmp, val...)
	}
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
	page.Render(io.MultiWriter(f))
	html, _ := ioutil.ReadFile(datapath + "line.html")
	return binary.BytesToString(html)
}

func drawLine(XName, YName string, x, data []string) *charts.Line {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithYAxisOpts(opts.YAxis{
			Name: YName, // 纵坐标
			SplitLine: &opts.SplitLine{
				Show: false,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: XName, // 横坐标
		}),
		charts.WithLegendOpts(opts.Legend{Show: true, Bottom: "1px"}),
	)
	line.SetXAxis(x).
		AddSeries("price", generateLineData(data),
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
		server := bind(ctx.Event.GroupID)
		if len(server) != 0 {
			f(ctx, server)
			return
		}
		ctx.SendChain(message.Text("本群还没绑定区服呢"))
	}
}

func sendNotice(payload gjson.Result) {
	var rsp []message.MessageSegment
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		for _, g := range ctx.GetGroupList().Array() {
			grp := g.Get("group_id").Int()
			isEnable, bindArea := isEnable(grp)
			switch payload.Get("type").Int() {
			case 2001:
				var status string
				if bindArea == payload.Get("data.server").String() {
					switch payload.Get("data.status").Int() {
					case 1:
						status = "开服啦！！！，快上游戏了\n"
					case 0:
						status = "停服了！！！，该干活干活了，该睡觉睡觉了\n"
					}
					rsp = []message.MessageSegment{
						message.Text(payload.Get("data.server").String() + status),
					}
				} else {
					rsp = []message.MessageSegment{}
				}
			case 2002:
				rsp =
					[]message.MessageSegment{
						message.Text("有新的资讯请查收:\n"),
						message.Text(payload.Get("data.type").String() + "\n" + payload.Get("data.title").String() + "\n" +
							payload.Get("data.url").String() + "\n" + payload.Get("data.date").String()),
					}
			}
			if isEnable && len(rsp) != 0 {
				ctx.SendGroupMessage(grp, rsp)
			}
		}
		return true
	})
}

func tuilan(tuiType string) string {
	url := "https://m.pvp.xoyo.com/activitygw/activity/calendar/detail"
	m := make(map[string]interface{})
	if id, ok := tuiKey[tuiType]; ok {
		m = map[string]interface{}{"id": id}
		b, _ := json.Marshal(m)
		tuilanData, _ := web.PostData(url, "application/json", bytes.NewReader(b))
		return binary.BytesToString(tuilanData)
	} else {
		return ""
	}
	return ""
}

func parseDate(msg string) int64 {
	extract := timefinder.New().TimeExtract(msg)
	return carbon.Time2Carbon(extract[0]).Timestamp()
}

func drawTeam(teamId int) image.Image {
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
	team := getTeamInfo(teamId)
	title := strconv.Itoa(team.TeamId) + " " + team.Dungeon
	_, th := dc.MeasureString("哈")
	t := 1200/2 - (float64(len([]rune(title))) / 2)
	dc.DrawStringAnchored(title, t, th, 0.5, 0.5)
	dc.DrawStringAnchored(team.Comment, 1200/2-float64(len([]rune(team.Comment)))/2, 3*th, 0.5, 0.5)
	// 团队
	mSlice := getMemberInfo(teamId)
	dc.LoadFontFace(text.FontFile, 30)
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
		back, _ := gg.LoadImage(iconfile + strconv.Itoa(int(m.MentalId)) + ".png")
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
