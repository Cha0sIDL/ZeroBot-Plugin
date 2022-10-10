package cron

import (
	"fmt"
	"github.com/samber/lo"
	"strings"

	"github.com/golang-module/carbon/v2"

	binutils "github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/process"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/fumiama/cron"
	"github.com/gogo/protobuf/sortkeys"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/util"
)

const (
	ServiceName = "cron"
)

var history []int64

var last = carbon.Now().Timestamp() - carbon.Now().Timestamp()%60

var provinces = map[string]struct{}{"河北": {}, "山西": {}, "辽宁": {}, "吉林": {}, "黑龙江": {}, "江苏": {}, "浙江": {}, "安徽": {}, "福建": {}, "江西": {}, "山东": {}, "河南": {}, "湖北": {}, "湖南": {}, "广东": {}, "海南": {}, "四川": {}, "贵州": {}, "云南": {}, "陕西": {}, "甘肃": {}, "青海": {}, "台湾": {}, "内蒙古": {}, "广西": {}, "西藏": {}, "宁夏": {}, "新疆": {}, "北京": {}, "天津": {}, "上海": {}, "重庆": {}, "香港": {}, "澳门": {}}

func init() { // 一些定时器
	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	_, err := c.AddFunc("@every 30s", func() { sendMessage30s() })
	if err == nil {
		c.Start()
	}
	control.Register(ServiceName, &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: true,
		Help:             "- 一些定时任务\n- 目前有地震播报功能",
	})
}

func sendMessage30s() {
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

// func sendEarthquake(ctx *zero.Ctx, grpIds []int64) {
//	now := carbon.Now().Timestamp()
//	data, _ := json.Marshal(map[string]string{
//		"action":     "requestMonitorDataAction",
//		"dataSource": "CEIC",
//		"startTime":  util.Interface2String(last * 1000),
//		"endTime":    util.Interface2String(now * 1000),
//	})
//	// "http://www.ceic.ac.cn/ajax/speedsearch?page=1&&num=6"
//	// data, _ := web.GetData("http://www.ceic.ac.cn/ajax/speedsearch?page=1&&num=6")
//	// gjson.Get(strings.Trim(binutils.BytesToString(data), "()"), "shuju")
//	rspData, err := web.PostData("http://api.dizhensubao.getui.com/api.htm", "application/json", bytes.NewReader(data))
//	if err != nil {
//		log.Errorln("cron error ", err)
//		return
//	}
//	strData := binutils.BytesToString(rspData)
//	for _, d := range gjson.Get(strData, "values").Array() {
//		last = now
//		_, ok := provinces[d.Get("loc_province").String()]
//		_, hisOk := history[d.Get("time").Int()]
//		lv := d.Get("mag").Float()
//		if ok && lv >= 3.5 && !hisOk {
//			for _, grpId := range grpIds {
//				ctx.SendGroupMessage(grpId, []message.MessageSegment{
//					message.Text(fmt.Sprintf("检测到 %s 于 "+carbon.CreateFromTimestamp(d.Get("time").Int()/1000).ToDateTimeString()+" 发生 %.1f 级地震，请处于震中位置人员注意安全~", d.Get("loc_name").String(), lv)),
//				})
//				process.SleepAbout1sTo2s()
//			}
//			history[d.Get("time").Int()] = struct{}{}
//			if len(history) > 128 {
//				history = nil // 防止无限增长
//				runtime.GC()
//			}
//		}
//	}
// } 883734530

func sendEarthquake(ctx *zero.Ctx, grpIds []int64) {
	// "http://www.ceic.ac.cn/ajax/speedsearch?page=1&&num=6"
	data, err := util.ProxyHttp(web.NewDefaultClient(), "http://www.ceic.ac.cn/ajax/speedsearch?page=1&&num=1", "GET", "", web.RandUA(), nil)
	if err != nil {
		return
	}
	earth := gjson.Get(strings.Trim(binutils.BytesToString(data), "()"), "shuju").Array()
	count := len(history)
	for _, earthData := range earth {
		id := earthData.Get("id").Int()
		if count != 0 && !sliceFind(id) && regionFind(earthData.Get("LOCATION_C").String()) && earthData.Get("M").Float() > 4.0 {
			log.Errorln("corn history:", history)
			for _, grpId := range grpIds {
				ctx.SendGroupMessage(grpId, []message.MessageSegment{
					message.Text(fmt.Sprintf("检测到 %s 于 %s 发生 %.1f 级地震，请处于震中位置人员注意安全~", earthData.Get("LOCATION_C").String(), earthData.Get("O_TIME").String(), earthData.Get("M").Float())),
				})
				process.SleepAbout1sTo2s()
			}
		}
		history = append(history, id)
	}
	history = lo.Uniq(history)
	if len(history) > 20 {
		sortkeys.Int64s(history)
		history = history[len(history)-15:]
	}
}

func sliceFind(id int64) bool {
	for _, h := range history {
		if h == id {
			return true
		}
	}
	return false
}

func regionFind(str string) bool {
	for key := range provinces {
		if strings.Contains(str, key) {
			return true
		}
	}
	return false
}
