// Package weather 查询天气插件
package weather

import (
	"fmt"
	"net/url"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"

	binutils "github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/config"
	"github.com/FloatTech/ZeroBot-Plugin/util"
)

const (
	servicename       = "weather"
	geoURL            = "https://geoapi.qweather.com/v2/city/lookup?"
	weatherURL        = "https://devapi.qweather.com/v7/weather/"
	weatherWarningURL = "https://devapi.qweather.com/v7/warning/now?"
)

func init() {
	engine := control.Register(servicename, &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "查询天气",
		Help:             "- xxx天气\n",
		PublicDataFolder: "Weather",
	})
	datapath := file.BOTPATH + "/" + engine.DataFolder()
	engine.OnSuffix("天气").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			city := strings.ReplaceAll(ctx.State["args"].(string), " ", "")
			geo := getGeo(city) // geo数据
			lat := gjson.Get(geo, "location.0.lat").Float()
			lon := gjson.Get(geo, "location.0.lon").Float()
			cityName := gjson.Get(geo, "location.0.name").String()
			todayWeather := getWeather("now", lat, lon) // 当天预报
			dailyWeather := getWeather("7d", lat, lon)  // 七天预报
			warning := getWarning(lat, lon)
			JMap := util.MergeMap(util.JSONToMap(todayWeather), util.JSONToMap(dailyWeather), util.JSONToMap(warning))
			JMap["city"] = cityName
			html := util.Template2html("weather.html", JMap)
			finName, err := util.HTML2pic(datapath, util.TodayFileName(), html)
			if city == "" {
				ctx.SendChain(message.Text("你还没有输入城市名字呢！"))
				return
			}
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			ctx.SendChain(message.Image("file:///" + finName))
		})
}

func getGeo(city string) string {
	api := geoURL + fmt.Sprintf("key=%s", config.Cfg.Weather) + "&location=" + url.QueryEscape(city)
	data, _ := web.RequestDataWith(web.NewDefaultClient(), api, "GET", "", web.RandUA(), nil)
	return binutils.BytesToString(data)
}

func getWeather(apiType string, lat float64, lon float64) string {
	api := weatherURL + apiType + "?" + fmt.Sprintf("location=%.2f,%.2f&key=%s", lon, lat, config.Cfg.Weather)
	data, _ := web.RequestDataWith(web.NewDefaultClient(), api, "GET", "", web.RandUA(), nil)
	return binutils.BytesToString(data)
}

func getWarning(lat float64, lon float64) string {
	api := weatherWarningURL + fmt.Sprintf("location=%.2f,%.2f&key=%s", lon, lat, config.Cfg.Weather)
	data, _ := web.RequestDataWith(web.NewDefaultClient(), api, "GET", "", web.RandUA(), nil)
	return binutils.BytesToString(data)
}
