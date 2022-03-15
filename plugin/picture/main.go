package picture

import (
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/control/order"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"net/url"
)

const (
	servicename = "picture"
	pictureUrl  = "https://doutu.lccyy.com/doutu/items?"
)

func init() {
	engine := control.Register(servicename, order.AcquirePrio(), &control.Options{
		DisableOnDefault: false,
		Help:             "- xxx表情\n" + "可能会偶尔抽风",
	})
	engine.OnSuffix("表情").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			arg := ctx.State["args"].(string)
			if arg == "" {
				return
			}
			url := pictureUrl + fmt.Sprintf("pageNum=%d&pageSize=%d&keyword=", util.Rand(1, 10), util.Rand(1, 100)) + url.QueryEscape(arg)
			data, err := web.RequestDataWith(web.NewDefaultClient(), url, "GET", "", web.RandUA())
			if err != nil {
				ctx.SendChain(message.Text("服务出错了请稍后重试！"))
			}
			Items := gjson.Get(binary.BytesToString(data), "items").Array()
			ctx.SendChain(message.Image(Items[rand.Intn(len(Items))].Get("url").String()))
		})
}
