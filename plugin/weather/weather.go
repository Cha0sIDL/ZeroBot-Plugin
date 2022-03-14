package weather

import (
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/control/order"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
	zero "github.com/wdvxdr1123/ZeroBot"
	"strings"
)

const (
	servicename = "weather"
	txurl       = "https://view.inews.qq.com/g2/getOnsInfo?name=disease_h5"
)

func init() {
	engine := control.Register(servicename, order.AcquirePrio(), &control.Options{
		DisableOnDefault: false,
		Help:             "- xxx天气\n",
	})
	engine.OnRegex("^[一-龥]{0,5}天气").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			//city := ctx.State["regex_matched"]
			city := strings.ReplaceAll(ctx.ExtractPlainText(), "天气", "")
			geo := getGeo(city)
			//if city == "" {
			//	ctx.SendChain(message.Text("你还没有输入城市名字呢！"))
			//	return
			//}
			//data, time, err := queryEpidemic(city)
			//if err != nil {
			//	ctx.SendChain(message.Text("ERROR: ", err))
			//	return
			//}
			//if data == nil {
			//	ctx.SendChain(message.Text("没有找到【", city, "】城市的疫情数据."))
			//	return
			//}
			ctx.SendChain()
		})
}

func getGeo(city string) string {
	data, err := web.RequestDataWith(web.NewDefaultClient(), api, "GET", "", ua)
}
