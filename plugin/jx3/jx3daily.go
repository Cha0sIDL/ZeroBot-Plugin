package jx3

import (
	"encoding/json"
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/config"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/control/order"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const url = "https://www.jx3api.com/app/"

var method = "GET"

type jinjia struct {
	server    string
	wanbaolou []float64
	tieba     []float64
	qita      []float64
}

func init() {
	go startWs()
	en := control.Register("jx", order.AcquirePrio(), &control.Options{
		DisableOnDefault: false,
		Help: "- 日常任务xxx(eg 日常任务绝代天骄)\n" +
			"- 开服检查xxx(eg 开服检查绝代天骄)\n" +
			"- 金价查询xxx(eg 金价查询绝代天骄)\n" +
			"- 花价|花价查询 xxx xxx xxx(eg 花价 绝代天骄 绣球花 广陵邑)\n" +
			"- 小药\n" +
			"- xxx配装(eg 分山劲配装)\n" +
			"- xxx奇穴(eg 分山劲奇穴)\n" +
			"- xxx宏(eg 分山劲宏)\n" +
			"- 沙盘xxx(eg 沙盘绝代天骄)\n" +
			"- 装饰属性|装饰xxx(eg 装饰混沌此生)\n" +
			"- 奇遇条件xxx(eg 奇遇条件三山四海)\n" +
			"- 奇遇攻略xxx(eg 奇遇攻略三山四海)\n" +
			"- 维护公告\n" +
			"- JX骚话（不区分大小写）\n" +
			"- 舔狗\n" +
			"-（开启|关闭）jx推送\n" +
			"TODO:宏转图片",
	})
	en.OnRegex(`^(日常任务|日常)(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			data, err := util.SendHttp(url+"daily", []byte(getMental(strings.Replace(str, " ", "", -1))))
			if err != nil {
				log.Errorln("jx3daily:", err)
				ctx.SendChain(message.Text("出错了！！！可能是参数不对"))
			}
			json := gjson.ParseBytes(data)
			ctx.SendChain(message.Text(
				"日期: ", json.Get("data.date"), "\n",
				"大战: ", json.Get("data.dayWar").Str, "\n",
				"战场: ", json.Get("data.dayBattle").Str, "\n",
				"公共日常: ", json.Get("data.dayPublic").Str, "\n",
				"美人图: ", json.Get("data.dayDraw").Str, "\n",
				"十人周常：", json.Get("data.weekFive").Str, "\n",
			))
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
	en.OnRegex(`^(金价|金价查询)(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			if len(str) == 0 {
				ctx.SendChain(message.Text(
					"请输入区服",
				))
			} else {
				data := map[string]string{"server": strings.Replace(str, " ", "", -1)}
				reqbody, err := json.Marshal(data)
				rsp, err := util.SendHttp(url+"demon", reqbody)
				if err != nil {
					log.Errorln("jx3daily:", err)
				}
				json := gjson.ParseBytes(rsp)
				jin := jinjia{}
				for _, value := range json.Get("data").Array() {
					value.ForEach(func(key, v gjson.Result) bool {
						switch key.String() {
						case "server":
							jin.server = v.String()
						case "wanbaolou":
							jin.wanbaolou = append(jin.wanbaolou, v.Float())
						case "tieba":
							jin.tieba = append(jin.tieba, v.Float())
						case "dd373", "uu898", "5173", "7881":
							jin.qita = append(jin.qita, v.Float())
						}
						return true
					})
				}
				dateStr := time.Now().Format("2006/01/02 15:04:05")
				ctx.SendChain(message.Text(
					"服务器: ", jin.server, "\n",
					"万宝楼: ", util.AppendAny(util.Min(jin.wanbaolou), util.Max(jin.wanbaolou)), "\n",
					"贴吧: ", util.AppendAny(util.Min(jin.tieba), util.Max(jin.tieba)), "\n",
					"其他平台: ", util.AppendAny(util.Min(jin.qita), util.Max(jin.qita)), "\n",
					"时间：", dateStr, "\n",
				))
			}
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
	en.OnRegex(`^沙盘(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			if len(str) == 0 {
				ctx.SendChain(message.Text(
					"请输入区服",
				))
			} else {
				data := map[string]string{"server": strings.Replace(str, " ", "", -1)}
				reqbody, err := json.Marshal(data)
				rsp, err := util.SendHttp(url+"sand", reqbody)
				if err != nil {
					log.Errorln("jx3daily:", err)
				}
				json := gjson.ParseBytes(rsp)
				sandUrl := json.Get("data").Array()[0]
				ctx.SendChain(message.Text(
					"服务器：", sandUrl.Get("server"), "\n",
					"时间：", time.Unix(sandUrl.Get("time").Int(), 0).Format("2006/01/02 15:04:05"),
				))
				ctx.SendChain(message.Image(
					sandUrl.Get("url").String(),
				))
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
	en.OnRegex(`^奇遇条件(.*)`).SetBlock(true).
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
			name := ctx.State["args"].(string)
			data := map[string]string{"name": getMental(strings.Replace(name, " ", "", -1))}
			reqbody, err := json.Marshal(data)
			rsp, err := util.SendHttp(url+"heighten", reqbody)
			if err != nil {
				log.Errorln("jx3daily:", err)
			}
			json := gjson.ParseBytes(rsp)
			ctx.SendChain(
				message.Image(
					json.Get("data.url").String()),
			)
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
	en.OnSuffix("宏").SetBlock(true).
		//TODO 图片
		Handle(func(ctx *zero.Ctx) {
			name := ctx.State["args"].(string)
			if len(name) == 0 {
				ctx.SendChain(message.Text("请输入职业！！！！"))
			} else {
				//data := map[string]string{"name": getMental(strings.Replace(name, " ", "", -1))}
				//reqbody, err := json.Marshal(data)
				//rsp, err := util.SendHttp(url+"macro", reqbody)
				//if err != nil {
				//	log.Errorln("jx3daily:", err)
				//}
				//json := gjson.ParseBytes(rsp)
				//ctx.SendChain(
				//	message.Text("奇穴：\n", json.Get("data.qixue").String(), "\n", "宏：\n", json.Get("data.macro").String()),
				//)
				mental := getData(strings.Replace(name, " ", "", -1))
				b, err := ioutil.ReadFile(dbpath + "macro/" + strconv.FormatUint(mental.ID, 10))
				if err != nil {
					ctx.SendChain(message.Text("请检查参数或通知管理员更新数据"))
				} else {
					ctx.SendChain(message.Text(string(b)))
				}
			}
		})
	en.OnRegex(`^奇遇攻略(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := ctx.State["regex_matched"].([]string)[1]
			if len(name) == 0 {
				ctx.SendChain(message.Text("输入参数有误！！！"))
			} else {
				data := map[string]string{"name": strings.Replace(name, " ", "", -1)}
				reqbody, err := json.Marshal(data)
				rsp, err := util.SendHttp(url+"strategy", reqbody)
				if err != nil {
					log.Errorln("jx3daily:", err)
				}
				json := gjson.ParseBytes(rsp)
				ctx.SendChain(
					message.Image(json.Get("data.url").String()),
				)
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
			rsp, err := util.SendHttp("https://www.jx3api.com/share/random", nil)
			if err != nil {
				log.Errorln("jx3daily:", err)
			}
			json := gjson.ParseBytes(rsp)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(json.Get("data.text")))
		})
	en.OnKeyword(`渣男`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			for _, QQ := range config.Cfg.At {
				ctx.SendChain(message.At(QQ))
			}
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
	en.OnFullMatch("关闭jx推送", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			disable(ctx.Event.GroupID)
			ctx.Send(message.Text("关闭成功"))
		})
	en.OnPrefix("绑定区服", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			area := strings.Replace(ctx.State["args"].(string), " ", "", -1)
			server := make(map[string]int64)
			rsp, _ := util.SendHttp(url+"check", nil)
			json := gjson.ParseBytes(rsp)
			for _, value := range json.Get("data").Array() {
				server[value.Get("server").String()] = value.Get("id").Int()
			}
			log.Errorln(server)
			if _, ok := server[area]; ok {
				bindArea(ctx.Event.GroupID, area)
				ctx.Send(message.Text("绑定成功"))
			} else {
				ctx.Send(message.Text("区服输入有误"))
			}
		})
}

func sendNotice(payload gjson.Result) {
	var rsp []message.MessageSegment
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		for _, g := range ctx.GetGroupList().Array() {
			grp := g.Get("group_id").Int()
			isEnable, bindArea := isEnable(grp)
			switch payload.Get("type").Int() {
			case 2011:
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
			case 2012:
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
