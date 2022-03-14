package weather

import (
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/config"
	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/control/order"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"net/url"
	"strings"
)

const (
	servicename = "weather"
	geoUrl      = "https://geoapi.qweather.com/v2/city/lookup?"
	weatherUrl  = "https://devapi.qweather.com/v7/weather/now?"
)

func init() {
	engine := control.Register(servicename, order.AcquirePrio(), &control.Options{
		DisableOnDefault: false,
		Help:             "- xxx天气\n",
	})
	engine.OnRegex("^[一-龥]{0,5}天气").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			city := strings.ReplaceAll(ctx.ExtractPlainText(), "天气", "")
			geo := getGeo(city)
			lat := gjson.Get(geo, "location.0.lat").Float()
			lon := gjson.Get(geo, "location.0.lon").Float()
			cityName := gjson.Get(geo, "location.0.name").String()
			api := weatherUrl + fmt.Sprintf("location=%.2f,%.2f&key=%s", lon, lat, config.Cfg.Weather)
			data, err := web.RequestDataWith(web.NewDefaultClient(), api, "GET", "", web.RandUA())
			json := binary.BytesToString(data)
			if city == "" {
				ctx.SendChain(message.Text("你还没有输入城市名字呢！"))
				return
			}
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			if gjson.Get(json, "code").Int() != 200 {
				ctx.SendChain(message.Text("未查询到相关数据"))
				return
			}
			ctx.SendChain(message.Text(
				"数据更新时间：" + gjson.Get(json, "updateTime").String() + "\n" +
					cityName + "天气为：\n" +
					fmt.Sprintf("温度：%s 摄氏度 \n体感温度：%s 摄氏度 \n天气状况 : %s \n风向 ：%s \n风力 ：%s \n风速 ：%s公里/小时 \n相对湿度 ：%s %% \n大气压强 ：%s百帕 \n能见度 ：%s 公里 ", gjson.Get(json, "now.temp"),
						gjson.Get(json, "now.feelsLike"), gjson.Get(json, "now.txt"), gjson.Get(json, "now.windDir"), gjson.Get(json, "now.windScale"), gjson.Get(json, "now.windSpeed"), gjson.Get(json, "now.humidity"), gjson.Get(json, "now.pressure"), gjson.Get(json, "now.vis")),
			))
		})
}

func getGeo(city string) string {
	api := geoUrl + fmt.Sprintf("key=%s", config.Cfg.Weather) + "&location=" + url.QueryEscape(city)
	data, _ := web.RequestDataWith(web.NewDefaultClient(), api, "GET", "", web.RandUA())
	return binary.BytesToString(data)
}
