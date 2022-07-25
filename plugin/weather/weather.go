package weather

import (
	"fmt"
	"net/url"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/config"
	"github.com/FloatTech/ZeroBot-Plugin/util"
)

//// result geo数据
// type geo struct {
//	Location []location `json:"location"`
//}
//
// type location struct {
//	Name string `json:"name"`
//	Id   string `json:"id"`
//	Lat  string `json:"lat"`
//	Lon  string `json:"lon"`
//}

const (
	servicename       = "weather"
	geoUrl            = "https://geoapi.qweather.com/v2/city/lookup?"
	weatherUrl        = "https://devapi.qweather.com/v7/weather/"
	weatherWarningUrl = "https://devapi.qweather.com/v7/warning/now?"
)

func init() {
	engine := control.Register(servicename, &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "- xxx天气\n",
		PublicDataFolder: "Weather",
	})
	datapath := file.BOTPATH + "/" + engine.DataFolder()
	engine.OnRegex("^[一-龥]{0,5}天气").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			city := strings.ReplaceAll(ctx.ExtractPlainText(), "天气", "")
			geo := getGeo(city) // geo数据
			lat := gjson.Get(geo, "location.0.lat").Float()
			lon := gjson.Get(geo, "location.0.lon").Float()
			cityName := gjson.Get(geo, "location.0.name").String()
			todayWeather := getWeather("now", lat, lon) // 当天预报
			dailyWeather := getWeather("7d", lat, lon)  // 七天预报
			warning := getWarning(lat, lon)
			JMap := util.MergeMap(util.JsonToMap(todayWeather), util.JsonToMap(dailyWeather), util.JsonToMap(warning))
			JMap["city"] = cityName
			html := util.Template2html("weather.html", JMap)
			finName, err := util.Html2pic(datapath, util.TodayFileName(), "weather.html", html)
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

// geo数据
// type T struct {
//	Code     string `json:"code"`
//	Location []struct {
//		Name      string `json:"name"`
//		Id        string `json:"id"`
//		Lat       string `json:"lat"`
//		Lon       string `json:"lon"`
//		Adm2      string `json:"adm2"`
//		Adm1      string `json:"adm1"`
//		Country   string `json:"country"`
//		Tz        string `json:"tz"`
//		UtcOffset string `json:"utcOffset"`
//		IsDst     string `json:"isDst"`
//		Type      string `json:"type"`
//		Rank      string `json:"rank"`
//		FxLink    string `json:"fxLink"`
//	} `json:"location"`
//	Refer struct {
//		Sources []string `json:"sources"`
//		License []string `json:"license"`
//	} `json:"refer"`
//}

func getGeo(city string) string {
	api := geoUrl + fmt.Sprintf("key=%s", config.Cfg.Weather) + "&location=" + url.QueryEscape(city)
	data, _ := web.RequestDataWith(web.NewDefaultClient(), api, "GET", "", web.RandUA())
	return binary.BytesToString(data)
}

func getWeather(apiType string, lat float64, lon float64) string {
	api := weatherUrl + apiType + "?" + fmt.Sprintf("location=%.2f,%.2f&key=%s", lon, lat, config.Cfg.Weather)
	data, _ := web.RequestDataWith(web.NewDefaultClient(), api, "GET", "", web.RandUA())
	return binary.BytesToString(data)
}

func getWarning(lat float64, lon float64) string {
	api := weatherWarningUrl + fmt.Sprintf("location=%.2f,%.2f&key=%s", lon, lat, config.Cfg.Weather)
	data, _ := web.RequestDataWith(web.NewDefaultClient(), api, "GET", "", web.RandUA())
	return binary.BytesToString(data)
}
