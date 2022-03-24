package jx3

import (
	"encoding/json"
	"fmt"
	"github.com/DanPlayer/timefinder"
	"github.com/FloatTech/ZeroBot-Plugin/config"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/control/order"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/FloatTech/zbputils/math"
	"github.com/fogleman/gg"
	"github.com/golang-module/carbon/v2"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	"image"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const url = "https://www.jx3api.com/app/"

const realizeUrl = "https://www.jx3api.com/realize/"

const cloudUrl = "https://www.jx3api.com/cloud/"

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
			"- 骚话（不区分大小写）\n" +
			"- 舔狗\n" +
			"-（开启|关闭）jx推送\n" +
			"- /roll随机roll点\n" +
			"TODO:宏转图片",
	}).ApplySingle(ctxext.DefaultSingle)
	go func() {
		initialize()
	}()
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
			json := gjson.ParseBytes(rsp)
			if err != nil {
				log.Errorln("jx3daily:", err)
			} else {
				ctx.SendChain(
					message.Text(
						"辅助食品：", json.Get("data.auxiliary_food"), "\n",
						"辅助药品：", json.Get("data.auxiliary_drug"), "\n",
						"增强食品：", json.Get("data.heighten_food"), "\n",
						"增强药品：", json.Get("data.heighten_drug"), "\n",
					))
			}
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
			if utf8.RuneCountInString(name) > 5 {
				log.Println("name len max")
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
				mental := getMentalData(strings.Replace(name, " ", "", -1))
				b, err := ioutil.ReadFile(dbpath + "macro/" + strconv.FormatUint(mental.ID, 10))
				if err == nil {
					ctx.SendChain(message.Text(string(b)))
				}
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
			rsp, err := util.SendHttp(realizeUrl+"random", nil)
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
	//开团 时间 副本名 备注
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
			teamId, err := createNewTeam(startTime, dungeon, comment, leaderId)
			if err != nil {
				ctx.SendChain(message.Text("Error :", err))
				return
			}
			ctx.SendChain(message.Text("开团成功，团队id为：", teamId))
		})
	//报团 团队ID 心法 角色名 [是否双休] 按照报名时间先后默认排序
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
			if carbon.Now().TimestampWithSecond() >= Team.StartTime {
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
				SignUp:         carbon.Now().TimestampWithSecond(),
			}
			addMember(&member)
			ctx.SendChain(message.Text("报团成功"), message.Reply(ctx.Event.MessageID))
			ctx.SendChain(message.Text("当前团队:\n"), message.Image("base64://"+helper.BytesToString(util.Image2Base64(drawTeam(teamId)))))
		})
	en.OnPrefixGroup([]string{"撤销报团"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			teamId, _ := strconv.Atoi(commandPart[0])
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
					fmt.Sprintf("WHERE teamID = '%d' AND startTime > '%d'", d, carbon.Now().TimestampWithSecond()))
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
				fmt.Sprintf("WHERE leaderId = '%d' AND startTime > '%d'", ctx.Event.UserID, carbon.Now().TimestampWithSecond()))
			out := ""
			for _, data := range InfoSlice {
				out = out + fmt.Sprintf("团队id：%d,团长 ：%d,副本：%s，开始时间：%s，备注：%s\n",
					data.TeamId, data.LeaderId, data.Dungeon, carbon.CreateFromTimestamp(data.StartTime).ToDateTimeString(), data.Comment)
			}
			ctx.SendChain(message.Text(out))
		})
	//查看团队 teamid
	en.OnPrefixGroup([]string{"查看团队", "查询团队", "查团"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			teamId, _ := strconv.Atoi(commandPart[0])
			ctx.SendChain(message.Image("base64://" + helper.BytesToString(util.Image2Base64(drawTeam(teamId)))))
		})
	//申请团长 团牌
	en.OnPrefixGroup([]string{"申请团长"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			permission := 0
			var teamName string
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) == 1 {
				teamName = commandPart[0]
			}
			if ctx.Event.Sender.Role != "member" {
				permission = 1
			}
			teamName = ""
			err := newLeader(ctx.Event.UserID, ctx.Event.Sender.NickName, permission, teamName)
			if err == 0 {
				ctx.SendChain(message.Text("申请团长成功，请管理员同意审批。"))
			}
			if err == -1 {
				ctx.SendChain(message.Text("贵人多忘事，你已经申请过了"))
			}
		})
	//取消开团 团队id
	en.OnPrefixGroup([]string{"取消开团", "删除团队", "撤销团队", "撤销开团"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) < 1 {
				ctx.SendChain(message.Text("撤销开团参数有误"))
			}
			teamId, err := strconv.Atoi(commandPart[0])
			if err != nil {
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
	//同意审批@qq
	en.OnRegex(`^同意审批.*?(\d+)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			qq := math.Str2Int64(ctx.State["regex_matched"].([]string)[1])
			teamName := acceptLeader(qq)
			ctx.SendChain(message.At(qq), message.Text("已成为团长,团队名称为："),
				message.Text(teamName))
		})
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

func parseDate(msg string) int64 {
	extract := timefinder.New().TimeExtract(msg)
	return carbon.Time2Carbon(extract[0]).TimestampWithSecond()
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
	//画直线
	for i := 0; i < 1200; {
		dc.SetRGBA(255, 255, 255, 11)
		dc.SetLineWidth(1)
		dc.DrawLine(0, float64(i), 1200, float64(i))
		dc.Stroke()
		i += 200
	}
	//画直线
	for i := 200; i < 1200; {
		//dc.SetRGBA(255, 255, 255, 11)
		//dc.SetLineWidth(1)
		dc.DrawLine(float64(i), 200, float64(i), 1200)
		dc.Stroke()
		i += 200
	}
	dc.SetFontFace(Fonts)
	//队伍
	for i := 1; i < 6; i++ {
		dc.DrawString(strconv.Itoa(i)+"队", 40, float64(100+200*i))
	}
	//标题
	team := getTeamInfo(teamId)
	title := strconv.Itoa(team.TeamId) + " " + team.Dungeon
	_, th := dc.MeasureString("哈")
	t := 1200/2 - (float64(len([]rune(title))) / 2)
	dc.DrawStringAnchored(title, t, th, 0.5, 0.5)
	dc.DrawStringAnchored(team.Comment, 1200/2-float64(len([]rune(team.Comment)))/2, 3*th, 0.5, 0.5)
	//团队
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
