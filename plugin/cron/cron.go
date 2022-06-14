package cron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/FloatTech/zbputils/binary"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/web"
	"github.com/fumiama/cron"
	"github.com/golang-module/carbon/v2"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/util"
)

const (
	ServiceName = "cron"
)

var history = make(map[int64]struct{}, 128)

var last = carbon.Now().Timestamp() - carbon.Now().Timestamp()%60

var provinces = map[string]struct{}{"河北": {}, "山西": {}, "辽宁": {}, "吉林": {}, "黑龙江": {}, "江苏": {}, "浙江": {}, "安徽": {}, "福建": {}, "江西": {}, "山东": {}, "河南": {}, "湖北": {}, "湖南": {}, "广东": {}, "海南": {}, "四川": {}, "贵州": {}, "云南": {}, "陕西": {}, "甘肃": {}, "青海": {}, "台湾": {}, "内蒙古": {}, "广西": {}, "西藏": {}, "宁夏": {}, "新疆": {}, "北京": {}, "天津": {}, "上海": {}, "重庆": {}, "香港": {}, "澳门": {}}

func init() { // 一些定时器
	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", func() { sendMessage1min() })
	if err == nil {
		c.Start()
	}
	control.Register(ServiceName, &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "一些定时任务\n",
	})
}

func sendMessage1min() {
	m, ok := control.Lookup(ServiceName)
	if !ok {
		log.Errorln("cron Notify Error")
	}
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		var grpIds []int64
		for _, g := range ctx.GetGroupList().Array() {
			grp := g.Get("group_id").Int()
			if m.IsEnabledIn(grp) {
				grpIds = append(grpIds, grp)
			}
		}
		sendEarthquake(ctx, grpIds)
		return true
	})
}

func sendEarthquake(ctx *zero.Ctx, grpIds []int64) {
	now := carbon.Now().Timestamp()
	data, _ := json.Marshal(map[string]string{
		"action":     "requestMonitorDataAction",
		"dataSource": "CEIC",
		"startTime":  util.Interface2String(last * 1000),
		"endTime":    util.Interface2String(now * 1000),
	})
	last = now
	// "http://www.ceic.ac.cn/ajax/speedsearch?page=1&&num=6"
	// data, _ := web.GetData("http://www.ceic.ac.cn/ajax/speedsearch?page=1&&num=6")
	//gjson.Get(strings.Trim(binary.BytesToString(data), "()"), "shuju")
	rspData, err := web.PostData("http://api.dizhensubao.getui.com/api.htm", "application/json", bytes.NewReader(data))
	if err != nil {
		log.Errorln("cron error ", err)
		return
	}
	strData := binary.BytesToString(rspData)
	log.Errorln("cron debug data", strData)
	for _, d := range gjson.Get(strData, "values").Array() {
		log.Errorln("cron debug", strData, d, "start", now, "last", last)
		_, ok := provinces[d.Get("loc_province").String()]
		_, hisOk := history[d.Get("time").Int()]
		lv := d.Get("mag").Float()
		if ok && lv >= 3.5 && !hisOk {
			for _, grpId := range grpIds {
				ctx.SendGroupMessage(grpId, []message.MessageSegment{
					message.Text(fmt.Sprintf("检测到 %s 发生 %.1f 级地震，请处于震中位置人员注意安全~", d.Get("loc_name").String(), lv)),
				})
			}
			history[d.Get("time").Int()] = struct{}{}
			if len(history) > 128 {
				history = nil // 防止无限增长
				runtime.GC()
			}
		}
	}
}
