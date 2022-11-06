package jx3

import (
	"fmt"

	"github.com/golang-module/carbon/v2"

	"github.com/FloatTech/floatbox/process"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/fumiama/cron"
	log "github.com/sirupsen/logrus"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const (
	ServiceName       = "dailyNotify"
	notify      int64 = 30
)

var date = map[int]map[string]string{
	1: {"10:00:00": "帮会跑商：阴山商路", "19:00:00": "阵营祭天：出征祭祀"},
	2: {"20:00:00": "阵营攻防：逐鹿中原"},
	3: {"20:00:00": "世界首领：百溪，烂柯山"},
	4: {"20:00:00": "阵营攻防：逐鹿中原"},
	5: {"20:00:00": "世界首领：楚州，晟江"},
	6: {"12:00:00": "攻防前置：南屏山", "13:00:00": "阵营攻防：浩气盟；奇袭：恶人谷", "19:00:00": "阵营攻防：浩气盟；奇袭：恶人谷"},
	0: {"12:00:00": "攻防前置：昆仑", "13:00:00": "阵营攻防：恶人谷；奇袭：浩气盟", "19:00:00": "阵营攻防：恶人谷；奇袭：浩气盟"},
}

func init() { // 插件主体
	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", func() { sendMessage() })
	if err == nil {
		c.Start()
	}
	control.Register(ServiceName, &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "剑网日常播报",
		Help:             "- 剑网每周日常定时推送\n",
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
						ctx.SendGroupMessage(grp, []message.MessageSegment{message.Text(fmt.Sprintf("还有%d分钟 %s 活动就要开始了~", notify, msg))})
					}
				}
			}
			process.SleepAbout1sTo2s()
		}
		return true
	})
}
