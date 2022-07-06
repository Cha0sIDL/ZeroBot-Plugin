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

	"github.com/flosch/pongo2/v5"

	"github.com/antchfx/htmlquery"

	"github.com/fumiama/cron"

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

type sandBox struct {
	sandToken       string
	sandTokenExpire int64
}

var sand sandBox

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
	// for _, chat := range *config.Cfg.JxChat {
	//	go startChatWs(chat)
	//}
	pongo2.RegisterFilter("genSlices", func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
		out = pongo2.AsValue(make([]int, in.Len()))
		err = nil
		return out, err
	})
	en := control.Register("jx", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		PrivateDataFolder: "jx3",
		Help: "- 日常任务xxx(eg 日常任务绝代天骄)\n" +
			"- 开服检查xxx(eg 开服检查绝代天骄)\n" +
			"- 金价查询|金价xxx(eg 金价查询绝代天骄)\n" +
			"- 花价|花价查询 xxx xxx xxx(eg 花价 绝代天骄 绣球花 广陵邑)\n" +
			"- 小药\n" +
			"- xxx配装(eg 分山劲配装)\n" +
			"- xxx奇穴(eg 分山劲奇穴)\n" +
			"- 宏xxx(eg 宏分山劲)\n" +
			"- 沙盘xxx(eg 沙盘绝代天骄)\n" +
			"- 装饰属性|装饰xxx(eg 装饰混沌此生)\n" +
			"- 奇遇条件xxx(eg 奇遇条件三山四海)\n" +
			"- 攻略xxx(eg 攻略三山四海)\n" +
			"- 维护公告|更新公告\n" +
			"- 骚话\n" +
			"- 舔狗\n" +
			"-（开启|关闭）jx推送\n" +
			"- /roll随机roll点\n" +
			"- 物价xxx\n" +
			"- 绑定区服xxx\n" +
			"- 团队相关见 https://docs.qq.com/doc/DUGJRQXd1bE5YckhB",
	})
	c := cron.New()
	_, err := c.AddFunc("0 5 * * *", func() {
		err := updateTalk()
		if err != nil {
			return
		}
	})
	c.AddFunc("@every 30s", func() {
		news()
	})
	if err == nil {
		c.Start()
	}
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
				client := web.NewDefaultClient()
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
	// en.OnRegex(`^前置(.*)`).SetBlock(true).
	//	Handle(func(ctx *zero.Ctx) {
	//		name := ctx.State["regex_matched"].([]string)[1]
	//		data := map[string]string{"name": strings.Replace(name, " ", "", -1)}
	//		reqbody, err := json.Marshal(data)
	//		rsp, err := util.SendHttp(url+"require", reqbody)
	//		if err != nil {
	//			log.Errorln("jx3daily:", err)
	//		}
	//		json := gjson.ParseBytes(rsp)
	//		ctx.SendChain(
	//			message.Text(
	//				"名称：", json.Get("data.name"), "\n",
	//				"方法：", json.Get("data.means"), "\n",
	//				"前置：", json.Get("data.require"), "\n",
	//				"奖励：", json.Get("data.reward"), "\n",
	//			),
	//			message.Image(
	//				json.Get("data.upload").String()),
	//		)
	//	})
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
	// en.OnSuffix("奇穴").SetBlock(true).
	//	Handle(func(ctx *zero.Ctx) {
	//		name := ctx.State["args"].(string)
	//		if len(name) == 0 {
	//			ctx.SendChain(message.Text("请输入职业！！！！"))
	//		} else {
	//			data := map[string]string{"name": getMental(strings.Replace(name, " ", "", -1))}
	//			reqbody, err := json.Marshal(data)
	//			rsp, err := util.SendHttp(url+"qixue", reqbody)
	//			if err != nil {
	//				log.Errorln("jx3daily:", err)
	//			}
	//			json := gjson.ParseBytes(rsp)
	//			ctx.SendChain(
	//				message.Text("通用：\n"),
	//				message.Image(
	//					json.Get("data.all").String()),
	//				message.Text("\n吃鸡：\n"),
	//				message.Image(
	//					json.Get("data.longmen").String()),
	//			)
	//		}
	//	})
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
			var t Jokes
			db.Pick(dbTalk, &t)
			ctx.SendChain(message.Text(t.Talk))
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
	en.OnFullMatch("更新骚话", zero.SuperUserPermission).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			err := updateTalk()
			if err != nil {
				ctx.SendChain(message.Text("更新失败", err))
				return
			}
			num, _ := db.Count(dbTalk)
			ctx.SendChain(message.Text(fmt.Sprintf("更新成功,本次共更新%d条骚话", num)))
		})
	en.OnPrefixGroup([]string{"属性"}).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			attributes(ctx, datapath)
		},
	)
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
		jin := gjson.Get(strData, fmt.Sprintf("data.%s", val[0]))
		rsp += fmt.Sprintf("今日%s平均金价为：\n", val[0])
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

func updateTalk() error {
	url := "https://cms.jx3box.com/api/cms/post/jokes?per=%d&page=%d"
	var page int64 = 1
	per := 30
	var Mutex sync.Mutex
	Mutex.Lock()
	defer Mutex.Unlock()
	for {
		data, err := web.GetData(fmt.Sprintf(url, per, page))
		jsonData := binary.BytesToString(data)
		if err != nil {
			return err
		}
		for _, talkData := range gjson.Get(jsonData, "data.list").Array() {
			// isFind := db.CanFind(dbTalk, fmt.Sprintf("where id=%d", talkData.Get("id").Int()))
			// if isFind {
			//	return nil
			//}
			db.Insert(dbTalk, &Jokes{
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

func attributes(ctx *zero.Ctx, datapath string) {
	// ts := "20220705054932497"
	commandPart := util.SplitSpace(ctx.State["args"].(string))
	if len(commandPart) != 2 {
		ctx.SendChain(message.Text("参数输入有误！\n" + "物价 绝代天骄 xxx"))
		return
	}
	// server := commandPart[0]
	//name := commandPart[1]
	//if normServer, ok := allServer[server]; ok {
	//	var user User
	//	zone := normServer[1]
	//	server = normServer[0]
	//	err := db.Find(dbUser, &user, fmt.Sprintf("WHERE id = '%s'", name+"_"+chatServer[server]))
	//	if err != nil {
	//		ctx.SendChain(message.Text("数据库查询失败了,请联系管理员看看吧~", err))
	//		return
	//	}
	//	if len(user.Data) == 0 {
	//		ctx.SendChain(message.Text("没有查到这个角色呢,试着在世界频道说句话试试吧~"))
	//		return
	//	}
	//	gameRoleId := gjson.Parse(user.Data).Get("body.msg.0.sRoleId").String()
	//	body, _ := json.Marshal(map[string]string{
	//		"ts":           ts,
	//		"game_role_id": gameRoleId,
	//		"server":       server,
	//		"zone":         zone,
	//	})
	//	data, err := web.PostData("https://m.pvp.xoyo.com/mine/equip/get-role-equip", "application/json", bytes.NewReader(body))
	//	if err != nil {
	//		ctx.SendChain(message.Text("请求出错了，稍后试试吧~"))
	//	}
	//jsonObj := gjson.ParseBytes(data).String()
	jsonObj := `
 {
        "CountFiveStone": 6,
        "Equips": [{
            "Base1Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base1Max": "",
                "Base1Min": "",
                "Desc": "atInvalid"
            },
            "Base2Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base2Max": "",
                "Base2Min": "",
                "Desc": "atInvalid"
            },
            "Base3Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base3Max": "",
                "Base3Min": "",
                "Desc": "atInvalid"
            },
            "BelongForce": "天策(傲血战意)、唐门(惊羽诀)、丐帮、霸刀",
            "BelongKungfu": "aoxue,jingyu,xiaochen,beiao",
            "BelongSchool": "通用",
            "Color": "4",
            "ColorActivateLevel": 0,
            "Desc": "",
            "DetailType": "",
            "EquipBoxStrengthLevel": "66666",
            "EquipItemStrengthLevel": "0",
            "FiveStoneScore": 0,
            "ID": "33717",
            "Icon": {
                "FileName": "https://dl.pvp.xoyo.com/prod/icons/tkt_ring11.png?v=2",
                "Kind": "饰品",
                "SubKind": "戒指"
            },
            "IncreaseQuality": 533,
            "JinglianScore": 479,
            "Level": "110",
            "MaxDurability": "0",
            "MaxEquipBoxStrengthLevel": "6",
            "MaxStrengthLevel": "6",
            "ModifyType": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "体质提高1278",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atVitalityBase",
                "Increase": 96,
                "Param1Max": "1278",
                "Param1Min": "1278",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高248",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "Increase": 19,
                "Param1Max": "248",
                "Param1Min": "248",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高402",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "Increase": 30,
                "Param1Max": "402",
                "Param1Min": "402",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功破防等级提高1244",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsOvercomeBase",
                "Increase": 93,
                "Param1Max": "1244",
                "Param1Min": "1244",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "破招值提高1106",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atSurplusValueBase",
                "Increase": 83,
                "Param1Max": "1106",
                "Param1Min": "1106",
                "Param2Max": "",
                "Param2Min": ""
            }],
            "Name": "未迟戒指",
            "Quality": "7100",
            "Score": 6390,
            "SetID": "",
            "SetName": "",
            "StrengthLevel": "6",
            "TabType": "8",
            "UID": "186665",
            "UcPos": "6",
            "WDurability": "0",
            "equipBelongs": [{
                "getDesc": "",
                "getType": "副本",
                "mapName": "",
                "order": "",
                "source": "副本：25人英雄河阳之战 — 阿阁诺,周通忌"
            }]
        }, {
            "Base1Type": {
                "Attrib": {
                    "GeneratedBase": "外功防御等级提高265",
                    "GeneratedMagic": "外功防御等级提高265",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base1Max": "265",
                "Base1Min": "265",
                "Desc": "atPhysicsShieldBase"
            },
            "Base2Type": {
                "Attrib": {
                    "GeneratedBase": "内功防御等级提高212",
                    "GeneratedMagic": "内功防御等级提高212",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base2Max": "212",
                "Base2Min": "212",
                "Desc": "atMagicShield"
            },
            "Base3Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base3Max": "",
                "Base3Min": "",
                "Desc": "atInvalid"
            },
            "BelongForce": "天策(傲血战意)、唐门(惊羽诀)、丐帮、霸刀",
            "BelongKungfu": "aoxue,jingyu,xiaochen,beiao",
            "BelongSchool": "通用",
            "Color": "4",
            "ColorActivateLevel": 0,
            "Desc": "",
            "DetailType": "",
            "EquipBoxStrengthLevel": "5",
            "EquipItemStrengthLevel": "0",
            "FiveStone": [{
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "破招值提高80",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atSurplusValueBase",
                "EnchantId": "6216",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 93,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "80",
                "Param1Min": "80",
                "SlotIdx": 0
            }, {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功会心等级提高80",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsCriticalStrike",
                "EnchantId": "6216",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 93,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "80",
                "Param1Min": "80",
                "SlotIdx": 1
            }],
            "FiveStoneScore": 602,
            "ID": "59825",
            "Icon": {
                "FileName": "https://dl.pvp.xoyo.com/prod/icons/cloth_21_10_28_104.png?v=2",
                "Kind": "防具",
                "SubKind": "裤子"
            },
            "IncreaseQuality": 396,
            "JinglianScore": 713,
            "Level": "110",
            "MaxDurability": "4320",
            "MaxEquipBoxStrengthLevel": "6",
            "MaxStrengthLevel": "6",
            "ModifyType": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "体质提高2593",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atVitalityBase",
                "Increase": 143,
                "Param1Max": "2593",
                "Param1Min": "2593",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高503",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "Increase": 28,
                "Param1Max": "503",
                "Param1Min": "503",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高816",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "Increase": 45,
                "Param1Max": "816",
                "Param1Min": "816",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功会心等级提高2523",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsCriticalStrike",
                "Increase": 139,
                "Param1Max": "2523",
                "Param1Min": "2523",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "无双等级提高2242",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrainBase",
                "Increase": 123,
                "Param1Max": "2242",
                "Param1Min": "2242",
                "Param2Max": "",
                "Param2Min": ""
            }],
            "Name": "晴辰裤",
            "Quality": "7200",
            "Score": 12960,
            "Set": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "全会心等级提高695",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atAllTypeCriticalStrike",
                "Increase": 0,
                "Param1Max": "695",
                "Param1Min": "695",
                "Param2Max": "",
                "Param2Min": "",
                "SetNum": "2"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "无双等级提高695",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrainBase",
                "Increase": 0,
                "Param1Max": "695",
                "Param1Min": "695",
                "Param2Max": "",
                "Param2Min": "",
                "SetNum": "4"
            }],
            "SetID": "4766",
            "SetList": ["晴辰靴", "晴辰裤", "晴辰冠", "晴辰衣", "晴辰袖", "晴辰腰带"],
            "SetListMap": ["晴辰裤", "晴辰衣"],
            "SetName": "晴辰",
            "StrengthLevel": "5",
            "TabType": "7",
            "UID": "187389",
            "UcPos": "10",
            "WDurability": "4104",
            "WPermanentEnchant": {
                "Attributes": [{
                    "Attrib": {
                        "GeneratedBase": "",
                        "GeneratedMagic": "外功破防等级提高491",
                        "HorseBase": "",
                        "HorseMagic": "",
                        "Type": "Attribute",
                        "percent": false
                    },
                    "Attribute1Value1": "491",
                    "Attribute1Value2": "491",
                    "Compare": "",
                    "Desc": "atPhysicsOvercomeBase",
                    "DiamondCount": "",
                    "DiamondIntensity": "",
                    "DiamondType": ""
                }],
                "ID": "11439",
                "Icon": null,
                "Level": "",
                "Name": "奉天·裤·铸（外破）",
                "Type": ""
            },
            "WTemporaryEnchant": {
                "Attributes": [{
                    "Attrib": null,
                    "Attribute1Value1": "skill/装备/霸刀技能大附魔1.lua",
                    "Attribute1Value2": "",
                    "Compare": "",
                    "Desc": "atExecuteScript",
                    "DiamondCount": "",
                    "DiamondIntensity": "",
                    "DiamondType": ""
                }],
                "ID": "11399",
                "Icon": null,
                "Level": "",
                "Name": "云山经·玉简·北傲诀",
                "Type": ""
            },
            "equipBelongs": [{
                "getDesc": "",
                "getType": "商店",
                "mapName": "",
                "order": "",
                "source": "商店：叶鸦 — 绝世防具"
            }]
        }, {
            "Base1Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base1Max": "",
                "Base1Min": "",
                "Desc": "atInvalid"
            },
            "Base2Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base2Max": "",
                "Base2Min": "",
                "Desc": "atInvalid"
            },
            "Base3Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base3Max": "",
                "Base3Min": "",
                "Desc": "atInvalid"
            },
            "BelongForce": "天策(傲血战意)、唐门(惊羽诀)、丐帮、霸刀",
            "BelongKungfu": "aoxue,jingyu,xiaochen,beiao",
            "BelongSchool": "通用",
            "Color": "4",
            "ColorActivateLevel": 0,
            "Desc": "",
            "DetailType": "",
            "EquipBoxStrengthLevel": "6",
            "EquipItemStrengthLevel": "0",
            "FiveStone": [{
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功破防等级提高80",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsOvercomeBase",
                "EnchantId": "6216",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 93,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "80",
                "Param1Min": "80",
                "SlotIdx": 0
            }],
            "FiveStoneScore": 82,
            "ID": "33723",
            "Icon": {
                "FileName": "https://dl.pvp.xoyo.com/prod/icons/tkt_necklace08.png?v=2",
                "Kind": "饰品",
                "SubKind": "项链"
            },
            "IncreaseQuality": 533,
            "JinglianScore": 479,
            "Level": "110",
            "MaxDurability": "0",
            "MaxEquipBoxStrengthLevel": "6",
            "MaxStrengthLevel": "6",
            "ModifyType": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "体质提高1278",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atVitalityBase",
                "Increase": 96,
                "Param1Max": "1278",
                "Param1Min": "1278",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高248",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "Increase": 19,
                "Param1Max": "248",
                "Param1Min": "248",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高402",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "Increase": 30,
                "Param1Max": "402",
                "Param1Min": "402",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功会心等级提高1244",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsCriticalStrike",
                "Increase": 93,
                "Param1Max": "1244",
                "Param1Min": "1244",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "无双等级提高1106",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrainBase",
                "Increase": 83,
                "Param1Max": "1106",
                "Param1Min": "1106",
                "Param2Max": "",
                "Param2Min": ""
            }],
            "Name": "舟山项饰",
            "Quality": "7100",
            "Score": 6390,
            "SetID": "",
            "SetName": "",
            "StrengthLevel": "6",
            "TabType": "8",
            "UID": "186713",
            "UcPos": "5",
            "WDurability": "0",
            "equipBelongs": [{
                "getDesc": "",
                "getType": "副本",
                "mapName": "",
                "order": "",
                "source": "副本：25人英雄河阳之战 — 阿阁诺"
            }]
        }, {
            "Base1Type": {
                "Attrib": {
                    "GeneratedBase": "远程伤害提高1135",
                    "GeneratedMagic": "远程武器伤害提高1135",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base1Max": "1135",
                "Base1Min": "1135",
                "Desc": "atRangeWeaponDamageBase"
            },
            "Base2Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base2Max": "756",
                "Base2Min": "756",
                "Desc": "atRangeWeaponDamageRand"
            },
            "Base3Type": {
                "Attrib": {
                    "GeneratedBase": "速度1.5",
                    "GeneratedMagic": "远程武器速度1.5",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base3Max": "24",
                "Base3Min": "24",
                "Desc": "atRangeWeaponAttackSpeedBase"
            },
            "Base4Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "每秒伤害1008.7",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base4Max": "",
                "Base4Min": "",
                "Desc": "fWeaponSpeed"
            },
            "BelongForce": "天策(傲血战意)、唐门(惊羽诀)、丐帮、霸刀",
            "BelongKungfu": "aoxue,jingyu,xiaochen,beiao",
            "BelongSchool": "通用",
            "Color": "4",
            "ColorActivateLevel": 0,
            "Desc": "柳静海离开五毒后，唐书雁趁无人时悄悄打开了那个装着“孤鸿只影”的匣子，红锦绸缎，就这么缠在一只略有些丑陋的“鸭子”上。若还是少年，唐书雁定要笑他，“这是哪里来的野鸭”。可如今，她与他早已错过多时……若少年，若少年，多无忧。",
            "DetailType": "弓弦",
            "EquipBoxStrengthLevel": "5",
            "EquipItemStrengthLevel": "4",
            "FiveStone": [{
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高36",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "EnchantId": "6216",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 42,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "36",
                "Param1Min": "36",
                "SlotIdx": 0
            }],
            "FiveStoneScore": 472,
            "ID": "30900",
            "Icon": {
                "FileName": "https://dl.pvp.xoyo.com/prod/icons/wpn_longdis15.png?v=2",
                "Kind": "武器",
                "SubKind": "投掷囊"
            },
            "IncreaseQuality": 336,
            "JinglianScore": 362,
            "Level": "110",
            "MaxDurability": "3360",
            "MaxEquipBoxStrengthLevel": "6",
            "MaxStrengthLevel": "6",
            "ModifyType": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "体质提高1318",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atVitalityBase",
                "Increase": 72,
                "Param1Max": "1318",
                "Param1Min": "1318",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高256",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "Increase": 14,
                "Param1Max": "256",
                "Param1Min": "256",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高415",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "Increase": 23,
                "Param1Max": "415",
                "Param1Min": "415",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功破防等级提高1282",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsOvercomeBase",
                "Increase": 71,
                "Param1Max": "1282",
                "Param1Min": "1282",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "无双等级提高1140",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrainBase",
                "Increase": 63,
                "Param1Max": "1140",
                "Param1Min": "1140",
                "Param2Max": "",
                "Param2Min": ""
            }],
            "Name": "不足贵·守",
            "Quality": "6100",
            "Score": 6588,
            "SetID": "",
            "SetName": "",
            "StrengthLevel": "5",
            "TabType": "6",
            "UID": "183305",
            "UcPos": "2",
            "WDurability": "3192",
            "WPermanentEnchant": {
                "Attributes": [{
                    "Attrib": {
                        "GeneratedBase": "",
                        "GeneratedMagic": "外功破防等级提高441",
                        "HorseBase": "",
                        "HorseMagic": "",
                        "Type": "Attribute",
                        "percent": false
                    },
                    "Attribute1Value1": "441",
                    "Attribute1Value2": "441",
                    "Compare": "",
                    "Desc": "atPhysicsOvercomeBase",
                    "DiamondCount": "",
                    "DiamondIntensity": "",
                    "DiamondType": ""
                }],
                "ID": "11351",
                "Icon": null,
                "Level": "",
                "Name": "阴山鬼晶·暗器（外破）",
                "Type": ""
            },
            "equipBelongs": [{
                "getDesc": "",
                "getType": "任务",
                "mapName": "",
                "order": "",
                "source": "任务：北地盐物"
            }]
        }, {
            "Base1Type": {
                "Attrib": {
                    "GeneratedBase": "外功防御等级提高295",
                    "GeneratedMagic": "外功防御等级提高295",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base1Max": "295",
                "Base1Min": "295",
                "Desc": "atPhysicsShieldBase"
            },
            "Base2Type": {
                "Attrib": {
                    "GeneratedBase": "内功防御等级提高236",
                    "GeneratedMagic": "内功防御等级提高236",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base2Max": "236",
                "Base2Min": "236",
                "Desc": "atMagicShield"
            },
            "Base3Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base3Max": "",
                "Base3Min": "",
                "Desc": "atInvalid"
            },
            "BelongForce": "天策(傲血战意)、唐门(惊羽诀)、丐帮、霸刀",
            "BelongKungfu": "aoxue,jingyu,xiaochen,beiao",
            "BelongSchool": "通用",
            "Color": "4",
            "ColorActivateLevel": 0,
            "Desc": "",
            "DetailType": "",
            "EquipBoxStrengthLevel": "5",
            "EquipItemStrengthLevel": "0",
            "FiveStone": [{
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高18",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "EnchantId": "6216",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 21,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "18",
                "Param1Min": "18",
                "SlotIdx": 0
            }, {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功会心等级提高80",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsCriticalStrike",
                "EnchantId": "6216",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 93,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "80",
                "Param1Min": "80",
                "SlotIdx": 1
            }],
            "FiveStoneScore": 319,
            "ID": "59837",
            "Icon": {
                "FileName": "https://dl.pvp.xoyo.com/prod/icons/cloth_21_10_28_106.png?v=2",
                "Kind": "防具",
                "SubKind": "上衣"
            },
            "IncreaseQuality": 396,
            "JinglianScore": 713,
            "Level": "110",
            "MaxDurability": "4800",
            "MaxEquipBoxStrengthLevel": "6",
            "MaxStrengthLevel": "6",
            "ModifyType": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "体质提高2593",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atVitalityBase",
                "Increase": 143,
                "Param1Max": "2593",
                "Param1Min": "2593",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高503",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "Increase": 28,
                "Param1Max": "503",
                "Param1Min": "503",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高816",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "Increase": 45,
                "Param1Max": "816",
                "Param1Min": "816",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功破防等级提高2523",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsOvercomeBase",
                "Increase": 139,
                "Param1Max": "2523",
                "Param1Min": "2523",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "无双等级提高2242",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrainBase",
                "Increase": 123,
                "Param1Max": "2242",
                "Param1Min": "2242",
                "Param2Max": "",
                "Param2Min": ""
            }],
            "Name": "晴辰衣",
            "Quality": "7200",
            "Score": 12960,
            "Set": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "全会心等级提高695",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atAllTypeCriticalStrike",
                "Increase": 0,
                "Param1Max": "695",
                "Param1Min": "695",
                "Param2Max": "",
                "Param2Min": "",
                "SetNum": "2"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "无双等级提高695",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrainBase",
                "Increase": 0,
                "Param1Max": "695",
                "Param1Min": "695",
                "Param2Max": "",
                "Param2Min": "",
                "SetNum": "4"
            }],
            "SetID": "4766",
            "SetList": ["晴辰裤", "晴辰冠", "晴辰衣", "晴辰袖", "晴辰腰带", "晴辰靴"],
            "SetListMap": ["晴辰裤", "晴辰衣"],
            "SetName": "晴辰",
            "StrengthLevel": "5",
            "TabType": "7",
            "UID": "187401",
            "UcPos": "3",
            "WDurability": "4560",
            "WPermanentEnchant": {
                "Attributes": [{
                    "Attrib": {
                        "GeneratedBase": "",
                        "GeneratedMagic": "无双等级提高174",
                        "HorseBase": "",
                        "HorseMagic": "",
                        "Type": "Attribute",
                        "percent": false
                    },
                    "Attribute1Value1": "174",
                    "Attribute1Value2": "174",
                    "Compare": "",
                    "Desc": "atStrainBase",
                    "DiamondCount": "",
                    "DiamondIntensity": "",
                    "DiamondType": ""
                }],
                "ID": "11189",
                "Icon": null,
                "Level": "",
                "Name": "云扇·衣·无双",
                "Type": ""
            },
            "equipBelongs": [{
                "getDesc": "",
                "getType": "商店",
                "mapName": "",
                "order": "",
                "source": "商店：叶鸦 — 绝世防具"
            }]
        }, {
            "Base1Type": {
                "Attrib": {
                    "GeneratedBase": "外功防御等级提高205",
                    "GeneratedMagic": "外功防御等级提高205",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base1Max": "205",
                "Base1Min": "205",
                "Desc": "atPhysicsShieldBase"
            },
            "Base2Type": {
                "Attrib": {
                    "GeneratedBase": "内功防御等级提高164",
                    "GeneratedMagic": "内功防御等级提高164",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base2Max": "164",
                "Base2Min": "164",
                "Desc": "atMagicShield"
            },
            "Base3Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base3Max": "",
                "Base3Min": "",
                "Desc": "atInvalid"
            },
            "BelongForce": "霸刀",
            "BelongKungfu": "beiao",
            "BelongSchool": "霸刀",
            "Color": "4",
            "ColorActivateLevel": 0,
            "Desc": "敛刃收归埋风骨，累因虚名葬北沉。",
            "DetailType": "",
            "EquipBoxStrengthLevel": "5",
            "EquipItemStrengthLevel": "5",
            "FiveStone": [{
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高18",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "EnchantId": "6216",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 21,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "18",
                "Param1Min": "18",
                "SlotIdx": 0
            }, {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高36",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "EnchantId": "6216",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 42,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "36",
                "Param1Min": "36",
                "SlotIdx": 1
            }],
            "FiveStoneScore": 1121,
            "ID": "55989",
            "Icon": {
                "FileName": "https://dl.pvp.xoyo.com/prod/icons/cloth_21_10_28_75.png?v=2",
                "Kind": "防具",
                "SubKind": "帽子"
            },
            "IncreaseQuality": 344,
            "JinglianScore": 557,
            "Level": "110",
            "MaxDurability": "3840",
            "MaxEquipBoxStrengthLevel": "6",
            "MaxStrengthLevel": "6",
            "ModifyType": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "体质提高2026",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atVitalityBase",
                "Increase": 111,
                "Param1Max": "2026",
                "Param1Min": "2026",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高393",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "Increase": 22,
                "Param1Max": "393",
                "Param1Min": "393",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高637",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "Increase": 35,
                "Param1Max": "637",
                "Param1Min": "637",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功破防等级提高1971",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsOvercomeBase",
                "Increase": 108,
                "Param1Max": "1971",
                "Param1Min": "1971",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "破招值提高1752",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atSurplusValueBase",
                "Increase": 96,
                "Param1Max": "1752",
                "Param1Min": "1752",
                "Param2Max": "",
                "Param2Min": ""
            }],
            "Name": "择芳·敛刃冠",
            "Quality": "6250",
            "Score": 10125,
            "Set": [{
                "Analysis": "",
                "Attrib": {
                    "Desc": "“刀啸风吟”伤害提高10%,“项王击鼎”伤害提高10%",
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Equipmentrecipe",
                    "percent": false
                },
                "Desc": "atSetEquipmentRecipe",
                "Increase": 0,
                "Param1Max": "4290",
                "Param1Min": "4290",
                "Param2Max": "1",
                "Param2Min": "1",
                "SetNum": "2"
            }, {
                "Analysis": "",
                "Attrib": {
                    "Desc": "装备：施展外功伤害招式，一定几率提高自身外功会心几率4%，会心效果4%，持续6秒。",
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Skillevent",
                    "percent": false
                },
                "Desc": "atSkillEventHandler",
                "Increase": 0,
                "Param1Max": "1925",
                "Param1Min": "1925",
                "Param2Max": "0",
                "Param2Min": "0",
                "SetNum": "4"
            }],
            "SetID": "4230",
            "SetList": ["择芳·敛刃靴", "择芳·敛刃冠", "择芳·敛刃衣", "择芳·敛刃护手", "择芳·敛刃腰带"],
            "SetListMap": ["择芳·敛刃冠"],
            "SetName": "择芳·敛刃",
            "StrengthLevel": "5",
            "TabType": "7",
            "UID": "179214",
            "UcPos": "4",
            "WCommonEnchant": {
                "Desc": "若自身当前气血值大于自身最大气血值的75%，则破防等级提高551点。不在名剑大会中生效。",
                "ID": "11424"
            },
            "WDurability": "3648",
            "WPermanentEnchant": {
                "Attributes": [{
                    "Attrib": {
                        "GeneratedBase": "",
                        "GeneratedMagic": "加速等级提高289",
                        "HorseBase": "",
                        "HorseMagic": "",
                        "Type": "Attribute",
                        "percent": false
                    },
                    "Attribute1Value1": "289",
                    "Attribute1Value2": "289",
                    "Compare": "",
                    "Desc": "atHasteBase",
                    "DiamondCount": "",
                    "DiamondIntensity": "",
                    "DiamondType": ""
                }],
                "ID": "11051",
                "Icon": null,
                "Level": "",
                "Name": "奉天·头·甲（急速）",
                "Type": ""
            },
            "equipBelongs": []
        }, {
            "Base1Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base1Max": "",
                "Base1Min": "",
                "Desc": "atInvalid"
            },
            "Base2Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base2Max": "",
                "Base2Min": "",
                "Desc": "atInvalid"
            },
            "Base3Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base3Max": "",
                "Base3Min": "",
                "Desc": "atInvalid"
            },
            "BelongForce": "外功门派",
            "BelongKungfu": "physics",
            "BelongSchool": "精简",
            "Color": "4",
            "ColorActivateLevel": 0,
            "Desc": "",
            "DetailType": "",
            "EquipBoxStrengthLevel": "6",
            "EquipItemStrengthLevel": "0",
            "FiveStoneScore": 437,
            "ID": "32007",
            "Icon": {
                "FileName": "https://dl.pvp.xoyo.com/prod/icons/tkt_ring06.png?v=2",
                "Kind": "饰品",
                "SubKind": "戒指"
            },
            "IncreaseQuality": 152,
            "JinglianScore": 137,
            "Level": "110",
            "MaxDurability": "0",
            "MaxEquipBoxStrengthLevel": "6",
            "MaxStrengthLevel": "3",
            "ModifyType": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高801",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "Increase": 19,
                "Param1Max": "801",
                "Param1Min": "801",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "破招值提高617",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atSurplusValueBase",
                "Increase": 15,
                "Param1Max": "617",
                "Param1Min": "617",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功会心等级提高1172",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsCriticalStrike",
                "Increase": 28,
                "Param1Max": "1172",
                "Param1Min": "1172",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功会心效果等级提高617",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsCriticalDamagePowerBase",
                "Increase": 15,
                "Param1Max": "617",
                "Param1Min": "617",
                "Param2Max": "",
                "Param2Min": ""
            }],
            "Name": "暂息戒",
            "Quality": "6340",
            "Score": 5706,
            "SetID": "",
            "SetName": "",
            "StrengthLevel": "3",
            "TabType": "8",
            "UID": "179579",
            "UcPos": "7",
            "WDurability": "0",
            "WPermanentEnchant": {
                "Attributes": [{
                    "Attrib": {
                        "GeneratedBase": "",
                        "GeneratedMagic": "力道提高110",
                        "HorseBase": "",
                        "HorseMagic": "",
                        "Type": "Attribute",
                        "percent": false
                    },
                    "Attribute1Value1": "110",
                    "Attribute1Value2": "110",
                    "Compare": "",
                    "Desc": "atStrengthBase",
                    "DiamondCount": "",
                    "DiamondIntensity": "",
                    "DiamondType": ""
                }],
                "ID": "11485",
                "Icon": null,
                "Level": "",
                "Name": "阴山鬼晶·戒指（力道）",
                "Type": ""
            },
            "equipBelongs": [{
                "getDesc": "",
                "getType": "副本",
                "mapName": "",
                "order": "",
                "source": "副本：25人英雄雷域大泽 — 悉达罗摩,乌蒙贵；25人普通河阳之战 — 周通忌,常宿"
            }]
        }, {
            "Base1Type": {
                "Attrib": {
                    "GeneratedBase": "外功防御等级提高143",
                    "GeneratedMagic": "外功防御等级提高143",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base1Max": "143",
                "Base1Min": "143",
                "Desc": "atPhysicsShieldBase"
            },
            "Base2Type": {
                "Attrib": {
                    "GeneratedBase": "内功防御等级提高115",
                    "GeneratedMagic": "内功防御等级提高115",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base2Max": "115",
                "Base2Min": "115",
                "Desc": "atMagicShield"
            },
            "Base3Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base3Max": "",
                "Base3Min": "",
                "Desc": "atInvalid"
            },
            "BelongForce": "霸刀",
            "BelongKungfu": "beiao",
            "BelongSchool": "霸刀",
            "Color": "4",
            "ColorActivateLevel": 0,
            "Desc": "覆手狂澜惊寰宇，长风浩荡洗乾坤。",
            "DetailType": "",
            "EquipBoxStrengthLevel": "6",
            "EquipItemStrengthLevel": "0",
            "FiveStone": [{
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功会心等级提高80",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsCriticalStrike",
                "EnchantId": "6216",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 93,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "80",
                "Param1Min": "80",
                "SlotIdx": 0
            }, {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高36",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "EnchantId": "6223",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 42,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "36",
                "Param1Min": "36",
                "SlotIdx": 1
            }],
            "FiveStoneScore": 319,
            "ID": "58613",
            "Icon": {
                "FileName": "https://dl.pvp.xoyo.com/prod/icons/item_22_4_15_18.png?v=2",
                "Kind": "防具",
                "SubKind": "腰带"
            },
            "IncreaseQuality": 525,
            "JinglianScore": 662,
            "Level": "110",
            "MaxDurability": "2400",
            "MaxEquipBoxStrengthLevel": "6",
            "MaxStrengthLevel": "6",
            "ModifyType": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "体质提高1764",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atVitalityBase",
                "Increase": 132,
                "Param1Max": "1764",
                "Param1Min": "1764",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高342",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "Increase": 26,
                "Param1Max": "342",
                "Param1Min": "342",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高555",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "Increase": 42,
                "Param1Max": "555",
                "Param1Min": "555",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功会心等级提高1717",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsCriticalStrike",
                "Increase": 129,
                "Param1Max": "1717",
                "Param1Min": "1717",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "无双等级提高1526",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrainBase",
                "Increase": 114,
                "Param1Max": "1526",
                "Param1Min": "1526",
                "Param2Max": "",
                "Param2Min": ""
            }],
            "Name": "承霁·惊寰腰带",
            "Quality": "7000",
            "Score": 8820,
            "Set": [{
                "Analysis": "",
                "Attrib": {
                    "Desc": "装备：施展外功伤害招式，一定几率提高自身外功会心几率4%，会心效果4%，持续6秒。",
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Skillevent",
                    "percent": false
                },
                "Desc": "atSkillEventHandler",
                "Increase": 0,
                "Param1Max": "1925",
                "Param1Min": "1925",
                "Param2Max": "0",
                "Param2Min": "0",
                "SetNum": "2"
            }, {
                "Analysis": "",
                "Attrib": {
                    "Desc": "“刀啸风吟”伤害提高10%,“项王击鼎”伤害提高10%",
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Equipmentrecipe",
                    "percent": false
                },
                "Desc": "atSetEquipmentRecipe",
                "Increase": 0,
                "Param1Max": "4290",
                "Param1Min": "4290",
                "Param2Max": "1",
                "Param2Min": "1",
                "SetNum": "4"
            }],
            "SetID": "4638",
            "SetList": ["承霁·惊寰冠", "承霁·惊寰衣", "承霁·惊寰护手", "承霁·惊寰腰带", "承霁·惊寰靴"],
            "SetListMap": ["承霁·惊寰腰带", "承霁·惊寰护手", "承霁·惊寰靴"],
            "SetName": "承霁·惊寰",
            "StrengthLevel": "6",
            "TabType": "7",
            "UID": "184823",
            "UcPos": "8",
            "WDurability": "2280",
            "WPermanentEnchant": {
                "Attributes": [{
                    "Attrib": {
                        "GeneratedBase": "",
                        "GeneratedMagic": "无双等级提高174",
                        "HorseBase": "",
                        "HorseMagic": "",
                        "Type": "Attribute",
                        "percent": false
                    },
                    "Attribute1Value1": "174",
                    "Attribute1Value2": "174",
                    "Compare": "",
                    "Desc": "atStrainBase",
                    "DiamondCount": "",
                    "DiamondIntensity": "",
                    "DiamondType": ""
                }],
                "ID": "11191",
                "Icon": null,
                "Level": "",
                "Name": "云扇·腰·无双",
                "Type": ""
            },
            "equipBelongs": []
        }, {
            "Base1Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base1Max": "",
                "Base1Min": "",
                "Desc": "atInvalid"
            },
            "Base2Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base2Max": "",
                "Base2Min": "",
                "Desc": "atInvalid"
            },
            "Base3Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base3Max": "",
                "Base3Min": "",
                "Desc": "atInvalid"
            },
            "BelongForce": "天策(傲血战意)、唐门(惊羽诀)、丐帮、霸刀",
            "BelongKungfu": "aoxue,jingyu,xiaochen,beiao",
            "BelongSchool": "通用",
            "Color": "4",
            "ColorActivateLevel": 0,
            "Desc": "使用：大幅度提升自身外功破防等级，持续15秒。\\n",
            "DetailType": "",
            "EquipBoxStrengthLevel": "5",
            "EquipItemStrengthLevel": "0",
            "FiveStone": [{
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功破防等级提高80",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsOvercomeBase",
                "EnchantId": "6216",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 93,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "80",
                "Param1Min": "80",
                "SlotIdx": 0
            }],
            "FiveStoneScore": 82,
            "ID": "33741",
            "Icon": {
                "FileName": "https://dl.pvp.xoyo.com/prod/icons/tkt_pendant20.png?v=2",
                "Kind": "饰品",
                "SubKind": "腰坠"
            },
            "IncreaseQuality": 270,
            "JinglianScore": 243,
            "Level": "110",
            "MaxDurability": "0",
            "MaxEquipBoxStrengthLevel": "6",
            "MaxStrengthLevel": "4",
            "ModifyType": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "体质提高1278",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atVitalityBase",
                "Increase": 49,
                "Param1Max": "1278",
                "Param1Min": "1278",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高248",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "Increase": 9,
                "Param1Max": "248",
                "Param1Min": "248",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高402",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "Increase": 15,
                "Param1Max": "402",
                "Param1Min": "402",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功破防等级提高1244",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsOvercomeBase",
                "Increase": 47,
                "Param1Max": "1244",
                "Param1Min": "1244",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "无双等级提高1106",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrainBase",
                "Increase": 42,
                "Param1Max": "1106",
                "Param1Min": "1106",
                "Param2Max": "",
                "Param2Min": ""
            }],
            "Name": "思乐",
            "Quality": "7100",
            "Score": 6390,
            "SetID": "",
            "SetName": "",
            "StrengthLevel": "4",
            "TabType": "8",
            "UID": "186797",
            "UcPos": "9",
            "WDurability": "0",
            "equipBelongs": [{
                "getDesc": "",
                "getType": "副本",
                "mapName": "",
                "order": "",
                "source": "副本：25人英雄河阳之战 — 周贽"
            }]
        }, {
            "Base1Type": {
                "Attrib": {
                    "GeneratedBase": "近身伤害提高1474",
                    "GeneratedMagic": "近战武器伤害提高1474",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base1Max": "1474",
                "Base1Min": "1474",
                "Desc": "atMeleeWeaponDamageBase"
            },
            "Base2Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base2Max": "983",
                "Base2Min": "983",
                "Desc": "atMeleeWeaponDamageRand"
            },
            "Base3Type": {
                "Attrib": {
                    "GeneratedBase": "速度1.5",
                    "GeneratedMagic": "近战武器速度1.5",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base3Max": "24",
                "Base3Min": "24",
                "Desc": "atMeleeWeaponAttackSpeedBase"
            },
            "Base4Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "每秒伤害1310.3",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base4Max": "",
                "Base4Min": "",
                "Desc": "fWeaponSpeed"
            },
            "BelongForce": "霸刀",
            "BelongKungfu": "beiao",
            "BelongSchool": "霸刀",
            "Color": "4",
            "ColorActivateLevel": 3,
            "ColorStone": {
                "Attributes": [{
                    "Attrib": {
                        "GeneratedBase": "",
                        "GeneratedMagic": "力道提高55",
                        "HorseBase": "",
                        "HorseMagic": "",
                        "Type": "Attribute",
                        "percent": false
                    },
                    "Attribute1Value1": "55",
                    "Attribute1Value2": "55",
                    "Compare": "3",
                    "Desc": "atStrengthBase",
                    "DiamondCount": "13",
                    "DiamondIntensity": "45",
                    "DiamondType": "5"
                }, {
                    "Attrib": {
                        "GeneratedBase": "",
                        "GeneratedMagic": "外功破防等级提高488",
                        "HorseBase": "",
                        "HorseMagic": "",
                        "Type": "Attribute",
                        "percent": false
                    },
                    "Attribute1Value1": "488",
                    "Attribute1Value2": "488",
                    "Compare": "3",
                    "Desc": "atPhysicsOvercomeBase",
                    "DiamondCount": "14",
                    "DiamondIntensity": "75",
                    "DiamondType": "5"
                }, {
                    "Attrib": {
                        "GeneratedBase": "",
                        "GeneratedMagic": "外功会心效果等级提高975",
                        "HorseBase": "",
                        "HorseMagic": "",
                        "Type": "Attribute",
                        "percent": false
                    },
                    "Attribute1Value1": "975",
                    "Attribute1Value2": "975",
                    "Compare": "3",
                    "Desc": "atPhysicsCriticalDamagePowerBase",
                    "DiamondCount": "16",
                    "DiamondIntensity": "90",
                    "DiamondType": "5"
                }],
                "ID": "542",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_fire_9.png?v=2",
                    "Kind": "道具",
                    "SubKind": "可挑选"
                },
                "Level": "5",
                "Name": "彩·真刚·斩铁·痛击(伍)",
                "Type": "五彩石"
            },
            "Desc": "约践横刀绕尘霾，君当共解古人心。",
            "DetailType": "傲霜刀",
            "EquipBoxStrengthLevel": "6",
            "EquipItemStrengthLevel": "0",
            "FiveStone": [{
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功会心效果等级提高80",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsCriticalDamagePowerBase",
                "EnchantId": "6224",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_7a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 140,
                "Level": "7",
                "Name": "五行石（七级）",
                "Param1Max": "80",
                "Param1Min": "80",
                "SlotIdx": 0
            }, {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功会心等级提高80",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsCriticalStrike",
                "EnchantId": "6224",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_7a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 140,
                "Level": "7",
                "Name": "五行石（七级）",
                "Param1Max": "80",
                "Param1Min": "80",
                "SlotIdx": 1
            }, {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高18",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "EnchantId": "6216",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 21,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "18",
                "Param1Min": "18",
                "SlotIdx": 2
            }],
            "FiveStoneScore": 2124,
            "ID": "29558",
            "Icon": {
                "FileName": "https://dl.pvp.xoyo.com/prod/icons/wpn_21_10_28_32.png?v=2",
                "Kind": "武器",
                "SubKind": "霸刀"
            },
            "IncreaseQuality": 476,
            "JinglianScore": 1027,
            "Level": "110",
            "MaxDurability": "4800",
            "MaxEquipBoxStrengthLevel": "8",
            "MaxStrengthLevel": "6",
            "ModifyType": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "体质提高2740",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atVitalityBase",
                "Increase": 206,
                "Param1Max": "2740",
                "Param1Min": "2740",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高531",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "Increase": 40,
                "Param1Max": "531",
                "Param1Min": "531",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高2057",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "Increase": 154,
                "Param1Max": "2057",
                "Param1Min": "2057",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功会心等级提高2666",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsCriticalStrike",
                "Increase": 200,
                "Param1Max": "2666",
                "Param1Min": "2666",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "无双等级提高2843",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrainBase",
                "Increase": 213,
                "Param1Max": "2843",
                "Param1Min": "2843",
                "Param2Max": "",
                "Param2Min": ""
            }],
            "Name": "解尘",
            "Quality": "6340",
            "Score": 13694,
            "SetID": "",
            "SetName": "",
            "StrengthLevel": "6",
            "TabType": "6",
            "UID": "179503",
            "UcPos": "0",
            "WDurability": "4119",
            "WPermanentEnchant": {
                "Attributes": [{
                    "Attrib": {
                        "GeneratedBase": "",
                        "GeneratedMagic": "外功攻击提高130",
                        "HorseBase": "",
                        "HorseMagic": "",
                        "Type": "Attribute",
                        "percent": false
                    },
                    "Attribute1Value1": "130",
                    "Attribute1Value2": "130",
                    "Compare": "",
                    "Desc": "atPhysicsAttackPowerBase",
                    "DiamondCount": "",
                    "DiamondIntensity": "",
                    "DiamondType": ""
                }],
                "ID": "11061",
                "Icon": null,
                "Level": "",
                "Name": "奉天·兵·甲（外攻）",
                "Type": ""
            },
            "equipBelongs": [{
                "getDesc": "",
                "getType": "副本",
                "mapName": "",
                "order": "",
                "source": "副本：25人英雄雷域大泽 — 巨型尖吻凤,桑乔,悉达罗摩,尤珈罗摩,月泉淮；25人普通河阳之战 — 勒齐那,阿阁诺,周通忌,周贽"
            }]
        }, {
            "Base1Type": {
                "Attrib": {
                    "GeneratedBase": "外功防御等级提高200",
                    "GeneratedMagic": "外功防御等级提高200",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base1Max": "200",
                "Base1Min": "200",
                "Desc": "atPhysicsShieldBase"
            },
            "Base2Type": {
                "Attrib": {
                    "GeneratedBase": "内功防御等级提高160",
                    "GeneratedMagic": "内功防御等级提高160",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base2Max": "160",
                "Base2Min": "160",
                "Desc": "atMagicShield"
            },
            "Base3Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base3Max": "",
                "Base3Min": "",
                "Desc": "atInvalid"
            },
            "BelongForce": "霸刀",
            "BelongKungfu": "beiao",
            "BelongSchool": "霸刀",
            "Color": "4",
            "ColorActivateLevel": 0,
            "Desc": "覆手狂澜惊寰宇，长风浩荡洗乾坤。",
            "DetailType": "",
            "EquipBoxStrengthLevel": "6",
            "EquipItemStrengthLevel": "0",
            "FiveStone": [{
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高18",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "EnchantId": "6223",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 21,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "18",
                "Param1Min": "18",
                "SlotIdx": 0
            }, {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高36",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "EnchantId": "6223",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 42,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "36",
                "Param1Min": "36",
                "SlotIdx": 1
            }],
            "FiveStoneScore": 1465,
            "ID": "58585",
            "Icon": {
                "FileName": "https://dl.pvp.xoyo.com/prod/icons/item_22_4_15_13.png?v=2",
                "Kind": "防具",
                "SubKind": "护臂"
            },
            "IncreaseQuality": 525,
            "JinglianScore": 662,
            "Level": "110",
            "MaxDurability": "3360",
            "MaxEquipBoxStrengthLevel": "6",
            "MaxStrengthLevel": "6",
            "ModifyType": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "体质提高1764",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atVitalityBase",
                "Increase": 132,
                "Param1Max": "1764",
                "Param1Min": "1764",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高342",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "Increase": 26,
                "Param1Max": "342",
                "Param1Min": "342",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高555",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "Increase": 42,
                "Param1Max": "555",
                "Param1Min": "555",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功会心等级提高1717",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsCriticalStrike",
                "Increase": 129,
                "Param1Max": "1717",
                "Param1Min": "1717",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "破招值提高1526",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atSurplusValueBase",
                "Increase": 114,
                "Param1Max": "1526",
                "Param1Min": "1526",
                "Param2Max": "",
                "Param2Min": ""
            }],
            "Name": "承霁·惊寰护手",
            "Quality": "7000",
            "Score": 8820,
            "Set": [{
                "Analysis": "",
                "Attrib": {
                    "Desc": "装备：施展外功伤害招式，一定几率提高自身外功会心几率4%，会心效果4%，持续6秒。",
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Skillevent",
                    "percent": false
                },
                "Desc": "atSkillEventHandler",
                "Increase": 0,
                "Param1Max": "1925",
                "Param1Min": "1925",
                "Param2Max": "0",
                "Param2Min": "0",
                "SetNum": "2"
            }, {
                "Analysis": "",
                "Attrib": {
                    "Desc": "“刀啸风吟”伤害提高10%,“项王击鼎”伤害提高10%",
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Equipmentrecipe",
                    "percent": false
                },
                "Desc": "atSetEquipmentRecipe",
                "Increase": 0,
                "Param1Max": "4290",
                "Param1Min": "4290",
                "Param2Max": "1",
                "Param2Min": "1",
                "SetNum": "4"
            }],
            "SetID": "4638",
            "SetList": ["承霁·惊寰冠", "承霁·惊寰衣", "承霁·惊寰护手", "承霁·惊寰腰带", "承霁·惊寰靴"],
            "SetListMap": ["承霁·惊寰腰带", "承霁·惊寰护手", "承霁·惊寰靴"],
            "SetName": "承霁·惊寰",
            "StrengthLevel": "6",
            "TabType": "7",
            "UID": "184795",
            "UcPos": "12",
            "WCommonEnchant": {
                "Desc": "释放招式有10%几率对目标额外造成一次少量伤害效果，该效果每10秒最多触发一次。不在名剑大会中生效。",
                "ID": "11510"
            },
            "WDurability": "3192",
            "WPermanentEnchant": {
                "Attributes": [{
                    "Attrib": {
                        "GeneratedBase": "",
                        "GeneratedMagic": "外功破防等级提高491",
                        "HorseBase": "",
                        "HorseMagic": "",
                        "Type": "Attribute",
                        "percent": false
                    },
                    "Attribute1Value1": "491",
                    "Attribute1Value2": "491",
                    "Compare": "",
                    "Desc": "atPhysicsOvercomeBase",
                    "DiamondCount": "",
                    "DiamondIntensity": "",
                    "DiamondType": ""
                }],
                "ID": "11462",
                "Icon": null,
                "Level": "",
                "Name": "奉天·腕·绣（外破）",
                "Type": ""
            },
            "equipBelongs": []
        }, {
            "Base1Type": {
                "Attrib": {
                    "GeneratedBase": "外功防御等级提高143",
                    "GeneratedMagic": "外功防御等级提高143",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base1Max": "143",
                "Base1Min": "143",
                "Desc": "atPhysicsShieldBase"
            },
            "Base2Type": {
                "Attrib": {
                    "GeneratedBase": "内功防御等级提高115",
                    "GeneratedMagic": "内功防御等级提高115",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Base2Max": "115",
                "Base2Min": "115",
                "Desc": "atMagicShield"
            },
            "Base3Type": {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "percent": false
                },
                "Base3Max": "",
                "Base3Min": "",
                "Desc": "atInvalid"
            },
            "BelongForce": "霸刀",
            "BelongKungfu": "beiao",
            "BelongSchool": "霸刀",
            "Color": "4",
            "ColorActivateLevel": 0,
            "Desc": "覆手狂澜惊寰宇，长风浩荡洗乾坤。",
            "DetailType": "",
            "EquipBoxStrengthLevel": "6",
            "EquipItemStrengthLevel": "0",
            "FiveStone": [{
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高36",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "EnchantId": "6223",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 42,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "36",
                "Param1Min": "36",
                "SlotIdx": 0
            }, {
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功破防等级提高80",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsOvercomeBase",
                "EnchantId": "6223",
                "Icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/five_element_stone_6a.png?v=2",
                    "Kind": "道具",
                    "SubKind": "95级生活技能"
                },
                "IncreaseGeneratedMagic": 93,
                "Level": "6",
                "Name": "五行石（六级）",
                "Param1Max": "80",
                "Param1Min": "80",
                "SlotIdx": 1
            }],
            "FiveStoneScore": 602,
            "ID": "58641",
            "Icon": {
                "FileName": "https://dl.pvp.xoyo.com/prod/icons/item_22_4_15_17.png?v=2",
                "Kind": "防具",
                "SubKind": "鞋"
            },
            "IncreaseQuality": 525,
            "JinglianScore": 662,
            "Level": "110",
            "MaxDurability": "2400",
            "MaxEquipBoxStrengthLevel": "6",
            "MaxStrengthLevel": "6",
            "ModifyType": [{
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "体质提高1764",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atVitalityBase",
                "Increase": 132,
                "Param1Max": "1764",
                "Param1Min": "1764",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "力道提高342",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrengthBase",
                "Increase": 26,
                "Param1Max": "342",
                "Param1Min": "342",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功攻击提高555",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsAttackPowerBase",
                "Increase": 42,
                "Param1Max": "555",
                "Param1Min": "555",
                "Param2Max": "0",
                "Param2Min": "0"
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "外功破防等级提高1717",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atPhysicsOvercomeBase",
                "Increase": 129,
                "Param1Max": "1717",
                "Param1Min": "1717",
                "Param2Max": "",
                "Param2Min": ""
            }, {
                "Analysis": "1",
                "Attrib": {
                    "GeneratedBase": "",
                    "GeneratedMagic": "无双等级提高1526",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Attribute",
                    "percent": false
                },
                "Desc": "atStrainBase",
                "Increase": 114,
                "Param1Max": "1526",
                "Param1Min": "1526",
                "Param2Max": "",
                "Param2Min": ""
            }],
            "Name": "承霁·惊寰靴",
            "Quality": "7000",
            "Score": 8820,
            "Set": [{
                "Analysis": "",
                "Attrib": {
                    "Desc": "装备：施展外功伤害招式，一定几率提高自身外功会心几率4%，会心效果4%，持续6秒。",
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Skillevent",
                    "percent": false
                },
                "Desc": "atSkillEventHandler",
                "Increase": 0,
                "Param1Max": "1925",
                "Param1Min": "1925",
                "Param2Max": "0",
                "Param2Min": "0",
                "SetNum": "2"
            }, {
                "Analysis": "",
                "Attrib": {
                    "Desc": "“刀啸风吟”伤害提高10%,“项王击鼎”伤害提高10%",
                    "GeneratedBase": "",
                    "GeneratedMagic": "",
                    "HorseBase": "",
                    "HorseMagic": "",
                    "Type": "Equipmentrecipe",
                    "percent": false
                },
                "Desc": "atSetEquipmentRecipe",
                "Increase": 0,
                "Param1Max": "4290",
                "Param1Min": "4290",
                "Param2Max": "1",
                "Param2Min": "1",
                "SetNum": "4"
            }],
            "SetID": "4638",
            "SetList": ["承霁·惊寰护手", "承霁·惊寰腰带", "承霁·惊寰靴", "承霁·惊寰冠", "承霁·惊寰衣"],
            "SetListMap": ["承霁·惊寰腰带", "承霁·惊寰护手", "承霁·惊寰靴"],
            "SetName": "承霁·惊寰",
            "StrengthLevel": "6",
            "TabType": "7",
            "UID": "184851",
            "UcPos": "11",
            "WDurability": "2280",
            "WPermanentEnchant": {
                "Attributes": [{
                    "Attrib": {
                        "GeneratedBase": "",
                        "GeneratedMagic": "外功攻击提高221",
                        "HorseBase": "",
                        "HorseMagic": "",
                        "Type": "Attribute",
                        "percent": false
                    },
                    "Attribute1Value1": "221",
                    "Attribute1Value2": "221",
                    "Compare": "",
                    "Desc": "atPhysicsAttackPowerBase",
                    "DiamondCount": "",
                    "DiamondIntensity": "",
                    "DiamondType": ""
                }],
                "ID": "11466",
                "Icon": null,
                "Level": "",
                "Name": "奉天·鞋·绣（外攻）",
                "Type": ""
            },
            "equipBelongs": []
        }],
        "Kungfu": {
            "Attrib": {
                "atDecriticalDamagePowerBase": "1725",
                "atPhysicsAttackPowerBase": "1648",
                "atPhysicsShieldBase": "442"
            },
            "KungfuID": "10464",
            "Level": "12",
            "Name": "beiao"
        },
        "KungfuDisplayType": "strength",
        "MatchDetail": {
            "Level": 0,
            "atAgilityBase": 0,
            "atAgilityBasePercentAdd": 0,
            "atAllTypeCriticalDamagePowerBase": 0,
            "atAllTypeCriticalStrike": 0,
            "atAllTypeHitValue": 0,
            "atBasePotentialAdd": 0,
            "atCriticalDamagePowerBaseLevel": 188.97,
            "atCriticalStrikeLevel": 36.08,
            "atDecriticalDamagePowerBase": 0,
            "atDecriticalDamagePowerBaseLevel": 26.5,
            "atDodge": 0,
            "atHasteBase": 0,
            "atHasteBaseLevel": 289,
            "atLifeAdditional": 0,
            "atLunarAttackPowerBase": 0,
            "atLunarCriticalDamagePowerBase": 0,
            "atLunarCriticalStrike": 0,
            "atLunarHitValue": 0,
            "atLunarOvercomeBase": 0,
            "atMagicAttackPowerBase": 0,
            "atMagicCriticalDamagePowerBase": 0,
            "atMagicCriticalStrike": 0,
            "atMagicHitValue": 0,
            "atMagicOvercome": 0,
            "atMagicShield": 0,
            "atMagicShieldLevel": 6.84,
            "atMaxLifeAdditional": 0,
            "atNeutralAttackPowerBase": 0,
            "atNeutralCriticalDamagePowerBase": 0,
            "atNeutralCriticalStrike": 0,
            "atNeutralHitValue": 0,
            "atNeutralOvercomeBase": 0,
            "atOvercomeBaseLevel": 44.84,
            "atParryBase": 0,
            "atParryValueBase": 0,
            "atPhysicsAttackPowerBase": 0,
            "atPhysicsCriticalDamagePowerBase": 0,
            "atPhysicsCriticalStrike": 0,
            "atPhysicsHitValue": 0,
            "atPhysicsOvercomeBase": 0,
            "atPhysicsShieldAdditional": 0,
            "atPhysicsShieldBase": 0,
            "atPhysicsShieldBaseLevel": 9.88,
            "atPoisonAttackPowerBase": 0,
            "atPoisonCriticalDamagePowerBase": 0,
            "atPoisonCriticalStrike": 0,
            "atPoisonHitValue": 0,
            "atPoisonOvercomeBase": 0,
            "atSolarAndLunarAttackPowerBase": 0,
            "atSolarAndLunarCriticalDamagePowerBase": 0,
            "atSolarAndLunarCriticalStrike": 0,
            "atSolarAndLunarHitValue": 0,
            "atSolarAndLunarOvercomeBase": 0,
            "atSolarAttackPowerBase": 0,
            "atSolarCriticalDamagePowerBase": 0,
            "atSolarCriticalStrike": 0,
            "atSolarHitValue": 0,
            "atSolarOvercomeBase": 0,
            "atSpiritBase": 0,
            "atSpiritBasePercentAdd": 0,
            "atSpunkBase": 0,
            "atSpunkBasePercentAdd": 0,
            "atStrainBase": 0,
            "atStrainBaseLevel": 43.4,
            "atStrengthBase": 4947,
            "atStrengthBasePercentAdd": 0,
            "atSurplusValueBase": 5402,
            "atTherapyPowerBase": 0,
            "atToughnessBase": 0,
            "atToughnessBaseLevel": 0,
            "atVitalityBase": 0,
            "atVitalityBasePercentAdd": 0,
            "baseAttack": 11886,
            "name": "beiao",
            "score": 121984,
            "totalAttack": 19552,
            "totalLift": 320649,
            "type": "strength"
        },
        "Person": {
            "atAgilityBase": "38",
            "atLifeAdditional": "23766",
            "atManaAdditional": "5400",
            "atSpiritBase": "38",
            "atSpunkBase": "37",
            "atStrengthBase": "37",
            "atVitalityBase": "38",
            "body": "1",
            "experience": "8528970",
            "level": "110",
            "maxAssistExp": "100",
            "maxAssistTimes": "18",
            "maxStamina": "3000",
            "maxThew": "3400",
            "parryBaseRate": "300",
            "qiXueId": ["16728"],
            "qixueList": [{
                "desc": "“项王击鼎”会心几率提高10%，会心效果提高10%。",
                "icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/daoj_16_10_17_33.png?v=2",
                    "Kind": "技能",
                    "SubKind": "霸刀"
                },
                "level": 1,
                "name": "虎踞",
                "skill_id": "16692"
            }, {
                "desc": "“项王击鼎”三段触发的持续范围伤害每跳伤害递增10%。",
                "icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/daoj_16_10_17_24.png?v=2",
                    "Kind": "技能",
                    "SubKind": "霸刀"
                },
                "level": 2,
                "name": "沧雪",
                "skill_id": "16777"
            }, {
                "desc": "“项王击鼎”“破釜沉舟”命中目标数量降低为3个，伤害提高20%且无视目标50%的外功防御，“项王击鼎”第三段招式持续过程中，若自身未移动，则受到的伤害降低15%。",
                "icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/daoj_16_10_17_4.png?v=2",
                    "Kind": "技能",
                    "SubKind": "霸刀"
                },
                "level": 3,
                "name": "冥鼔",
                "skill_id": "26904"
            }, {
                "desc": "“醉斩白蛇”增加3跳，调息时间提高10秒，运功期间每跳伤害递增25%。",
                "icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/daoj_16_10_17_139.png?v=2",
                    "Kind": "技能",
                    "SubKind": "霸刀"
                },
                "level": 4,
                "name": "化蛟",
                "skill_id": "16779"
            }, {
                "desc": "“刀啸风吟”施展后使自身会心几率提高4%，会心效果提高4%，“雪絮金屏”套路招式伤害提高4%，气劲消耗降低10%，持续18秒，可以叠加2层。",
                "icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/daoj_16_10_17_10.png?v=2",
                    "Kind": "技能",
                    "SubKind": "霸刀"
                },
                "level": 5,
                "name": "含风",
                "skill_id": "25633"
            }, {
                "desc": "“项王击鼎”“割据秦宫”每段命中目标使“项王击鼎”三段触发的持续群体伤害持续时间增加1秒，“项王击鼎”第三段持续效果最多持续6.5秒。",
                "icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/daoj_16_10_17_36.png?v=2",
                    "Kind": "技能",
                    "SubKind": "霸刀"
                },
                "level": 6,
                "name": "逐鹿",
                "skill_id": "16748"
            }, {
                "desc": "“醉斩白蛇”施展后立刻对目标造成伤害，并提高自身的15%外功基础攻击力，持续16秒。",
                "icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/daoj_16_10_17_19.png?v=2",
                    "Kind": "技能",
                    "SubKind": "霸刀"
                },
                "level": 7,
                "name": "斩纷",
                "skill_id": "16733"
            }, {
                "desc": "体质和力道提高10%。",
                "icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/daoj_16_10_17_138.png?v=2",
                    "Kind": "技能",
                    "SubKind": "霸刀"
                },
                "level": 8,
                "name": "星火",
                "skill_id": "16728"
            }, {
                "desc": "命中气血值低于40%的目标，招式会心提高10%，会心效果提高10%。",
                "icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/daoj_16_10_17_140.png?v=2",
                    "Kind": "技能",
                    "SubKind": "霸刀"
                },
                "level": 9,
                "name": "楚歌",
                "skill_id": "16737"
            }, {
                "desc": "“闹须弥”流血效果伤害提高50%且自身伤害招式命中目标后有30%的几率刷新流血效果持续时间。",
                "icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/daoj_16_10_17_15.png?v=2",
                    "Kind": "技能",
                    "SubKind": "霸刀"
                },
                "level": 10,
                "name": "绝期",
                "skill_id": "17056"
            }, {
                "desc": "切换至“雪絮金屏”体态，使得自身下个“醉斩白蛇”运功时间降低25%，伤害提高25%。",
                "icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/daoj_16_10_17_72.png?v=2",
                    "Kind": "技能",
                    "SubKind": "霸刀"
                },
                "level": 11,
                "name": "砺锋",
                "skill_id": "26735"
            }, {
                "desc": "“坚壁清野”刀气残留时间增加5秒，“雪絮金屏”体态下施展伤害招式命中处于自身“坚壁清野”区域内的敌对目标，招式会心几率提高20%，会心效果提高20%。",
                "icon": {
                    "FileName": "https://dl.pvp.xoyo.com/prod/icons/daoj_16_10_17_50.png?v=2",
                    "Kind": "技能",
                    "SubKind": "霸刀"
                },
                "level": 12,
                "name": "心镜",
                "skill_id": "16912"
            }],
            "title": "0"
        },
        "PersonalPanel": [{
            "name": "攻击力",
            "percent": false,
            "value": 19552
        }, {
            "name": "会心",
            "percent": true,
            "value": 36.08
        }, {
            "name": "破防",
            "percent": true,
            "value": 44.84
        }, {
            "name": "化劲",
            "percent": true,
            "value": 26.5
        }, {
            "name": "气血",
            "percent": false,
            "value": 320649
        }, {
            "name": "基础攻击力",
            "percent": false,
            "value": 11886
        }, {
            "name": "会心效果",
            "percent": true,
            "value": 188.97
        }, {
            "name": "破招",
            "percent": false,
            "value": 5402
        }, {
            "name": "加速",
            "percent": false,
            "value": 289
        }, {
            "name": "无双",
            "percent": true,
            "value": 43.4
        }, {
            "name": "内功防御",
            "percent": true,
            "value": 6.84
        }, {
            "name": "外功防御",
            "percent": true,
            "value": 9.88
        }, {
            "name": "御劲",
            "percent": true,
            "value": 0
        }, {
            "name": "力道",
            "percent": false,
            "value": 4947
        }],
        "PveEquipsScore": 0,
        "PvpEquipsScore": 0,
        "Set": {
            "4230": 1,
            "4638": 3,
            "4766": 2
        },
        "TotalEquipsScore": 121984
    }
`
	templateData := map[string]interface{}{
		"name":   "123",
		"server": "123123",
		"data":   util.JsonToMap(jsonObj)}
	fmt.Println(util.Template2html("equip.html", templateData))

	// } else {
	//	ctx.SendChain(message.Text("输入区服有误，请检查qaq~"))
	//}
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

func news() {
	var msg []News
	count, _ := db.Count(dbNews)
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
		canFind := db.CanFind(dbNews, fmt.Sprintf("WHERE id = '%s'", href))
		data := News{
			ID:    href,
			Date:  date,
			Title: title,
			Kind:  kind,
		}
		if canFind {
			continue
		}
		err := db.Insert(dbNews, &data)
		if err != nil {
			continue
		}
		msg = append(msg, data)
	}
	if count == 0 {
		return
	}
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		for _, g := range ctx.GetGroupList().Array() {
			grp := g.Get("group_id").Int()
			isEnable, _ := isEnable(grp)
			if isEnable {
				for _, data := range msg {
					ctx.SendGroupMessage(grp, fmt.Sprintf("有新的资讯请查收:\n%s\n%s\n%s\n%s", data.Kind, data.Title, data.ID, data.Date))
				}
			}
		}
		return true
	})
}

func sendNotice(payload gjson.Result) {
	var rsp []message.MessageSegment
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		for _, g := range ctx.GetGroupList().Array() {
			grp := g.Get("group_id").Int()
			isEnable, bindArea := isEnable(grp)
			switch payload.Get("action").Int() {
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
