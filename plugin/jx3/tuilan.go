// Package jx3 推栏相关的接口
package jx3

import (
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/config"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	binutils "github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/go-resty/resty/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	datapath := file.BOTPATH + "/" + en.DataFolder()
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
