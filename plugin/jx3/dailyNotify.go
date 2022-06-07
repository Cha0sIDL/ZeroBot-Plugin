package jx3

import (
	"fmt"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/fumiama/cron"
	log "github.com/sirupsen/logrus"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/zbputils/control"
	"github.com/golang-module/carbon/v2"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const (
	ServiceName       = "dailyNotify"
	notify      int64 = 15
)

var date = map[int]map[string]string{
	1: {"10:00:00": "帮会跑商：阴山商路", "19:00:00": "阵营祭天：出征祭祀"},
	2: {"20:00:00": "阵营攻防：逐鹿中原"},
	3: {"20:00:00": "世界首领：少林·乱世，七秀·乱世"},
	4: {"20:00:00": "阵营攻防：逐鹿中原"},
	5: {"20:00:00": "世界首领：黑山林海，藏剑·乱世"},
	6: {"12:00:00": "攻防前置：南屏山", "13:00:00": "阵营攻防：浩气盟；奇袭：恶人谷", "19:00:00": "阵营攻防：浩气盟；奇袭：恶人谷"},
	0: {"12:00:00": "攻防前置：昆仑", "13:00:00": "阵营攻防：恶人谷；奇袭：浩气盟", "19:00:00": "阵营攻防：恶人谷；奇袭：浩气盟"},
}

// "一": "帮会跑商：阴山商路(10:00)\n阵营祭天：出征祭祀(19:00)\n",
// "二": "阵营攻防：逐鹿中原(20:00)\n",
// "三": "世界首领：少林·乱世，七秀·乱世(20:00)\n",
// "四": "阵营攻防：逐鹿中原(20:00)\n",
// "五": "世界首领：黑山林海，藏剑·乱世(20:00)\n",
//"六": "攻防前置：南屏山(12:00)\n阵营攻防：浩气盟；奇袭：恶人谷(13:00，19:00)\n",
//"日": "攻防前置：昆仑(12:00)\n阵营攻防：恶人谷；奇袭：浩气盟(13:00，19:00)\n"

// func init() { // 插件主体
//	engine := control.Register(ServiceName, &ctrl.Options[*zero.Ctx]{
//		DisableOnDefault: false,
//		Help:             "每周日常定时推送\n",
//	})
//	engine.OnFullMatch("剑网活动推送").Handle(
//		sendMessage,
//	)
//}
//
// func sendMessage(ctx *zero.Ctx) {
//	week := carbon.Now().Week()
//	daily := date[week]
//	for time, msg := range daily {
//		diff := carbon.Parse(carbon.Now().ToDateString() + " " + time).DiffInMinutes(carbon.Now())
//		if diff == -notify {
//			ctx.SendChain(message.AtAll(), message.Text(fmt.Sprintf(" 还有%d分钟 %s 活动就要开始了~", notify, msg)))
//		}
//	}
//}

func init() { // 插件主体
	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", func() { sendMessage() })
	if err == nil {
		c.Start()
	}
	control.Register(ServiceName, &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "每周日常定时推送\n",
	})
}

func sendMessage() {
	m, ok := control.Lookup(ServiceName)
	if !ok {
		log.Errorln("dailyNotify Err")
	}
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		for _, g := range ctx.GetGroupList().Array() {
			grp := g.Get("group_id").Int()
			if m.IsEnabledIn(grp) {
				week := carbon.Now().Week()
				daily := date[week]
				for time, msg := range daily {
					diff := carbon.Parse(carbon.Now().ToDateString() + " " + time).DiffInMinutes(carbon.Now())
					if diff == -notify {
						ctx.SendGroupMessage(grp, []message.MessageSegment{message.AtAll(), message.Text(fmt.Sprintf(" 还有%d分钟 %s 活动就要开始了~", notify, msg))})
					}
				}
			}
		}
		return true
	})
}
