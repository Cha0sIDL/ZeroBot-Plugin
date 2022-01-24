package jx3

import (
	"encoding/json"
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/order"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/zbputils/control"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
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
	en := control.Register("jx3", order.PrioJx3, &control.Options{
		DisableOnDefault: false,
		Help: "- 日常任务xxx(eg 日常任务绝代天骄)\n" +
			"- 开服检查xxx(eg 开服检查绝代天骄)\n" +
			"- 金价查询xxx(eg 金价查询绝代天骄)\n" +
			"- 花价|花价查询 xxx xxx xxx(eg 花价 绝代天骄 绣球花 广陵邑)\n" +
			"- 小药\n" +
			"- 配装xxx(eg 配装分山劲)\n" +
			"- 奇穴xxx(eg 奇穴分山劲)\n" +
			"- 宏xxx(eg 宏分山劲)\n" +
			"- 沙盘xxx(eg 沙盘绝代天骄)\n" +
			"- 装饰属性|装饰xxx(eg 装饰混沌此生)\n" +
			"- 奇遇条件xxx(eg 奇遇条件三山四海)\n" +
			"- 奇遇攻略xxx(eg 奇遇攻略三山四海)\n" +
			"- 维护公告\n" +
			"- JX骚话（不区分大小写）\n" +
			"- 舔狗\n" +
			"TODO:宏转图片",
	})
	en.OnRegex(`^(日常任务|日常)(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			log.Errorln("日常任务")
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
					"万宝楼: ", appendAny(min(jin.wanbaolou), max(jin.wanbaolou)), "\n",
					"贴吧: ", appendAny(min(jin.tieba), max(jin.tieba)), "\n",
					"其他平台: ", appendAny(min(jin.qita), max(jin.qita)), "\n",
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
	en.OnRegex(`^小药(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := ctx.State["regex_matched"].([]string)[1]
			data := map[string]string{"name": strings.Replace(name, " ", "", -1)}
			reqbody, err := json.Marshal(data)
			rsp, err := util.SendHttp(url+"heighten", reqbody)
			if err != nil {
				log.Errorln("jx3daily:", err)
			}
			json := gjson.ParseBytes(rsp)
			log.Errorln(json)
			ctx.SendChain(
				message.Image(
					json.Get("data.url").String()),
			)
		})
	en.OnRegex(`^配装(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := ctx.State["regex_matched"].([]string)[1]
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
	en.OnRegex(`^奇穴(.*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := ctx.State["regex_matched"].([]string)[1]
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
				log.Errorln(json)
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
	en.OnRegex(`^宏(.*)`).SetBlock(true).
		//TODO 图片
		Handle(func(ctx *zero.Ctx) {
			name := ctx.State["regex_matched"].([]string)[1]
			if len(name) == 0 {
				ctx.SendChain(message.Text("请输入职业！！！！"))
			} else {
				data := map[string]string{"name": getMental(strings.Replace(name, " ", "", -1))}
				reqbody, err := json.Marshal(data)
				rsp, err := util.SendHttp(url+"macro", reqbody)
				if err != nil {
					log.Errorln("jx3daily:", err)
				}
				json := gjson.ParseBytes(rsp)
				ctx.SendChain(
					message.Text("奇穴：\n", json.Get("data.qixue").String(), "\n", "宏：\n", json.Get("data.macro").String()),
				)
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
	en.OnRegex(`^(?i)jx骚话(.*)`).SetBlock(true).
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
}

func max(l []float64) (max float64) {
	max = l[0]
	for _, v := range l {
		if v > max {
			max = v
		}
	}
	return
}

func min(l []float64) (min float64) {
	min = l[0]
	for _, v := range l {
		if v < min {
			min = v
		}
	}
	return
}

func appendAny(a interface{}, b interface{}) string {
	return fmt.Sprintf("%v", a) + "-" + fmt.Sprintf("%v", b)
}
