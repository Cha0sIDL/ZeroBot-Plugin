package sulian

import (
	"encoding/json"
	"github.com/FloatTech/ZeroBot-Plugin/order"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/file"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io/ioutil"
	"math/rand"
)

func init() { // 插件主体
	engine := control.Register("sulian", order.PrioSulian, &control.Options{
		DisableOnDefault: false,
		Help:             "苏联笑话\n",
	})
	// 被喊名字
	engine.OnFullMatch("苏联笑话").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if file.IsNotExist(dbfile) {
				down()
			}
			data, err := ioutil.ReadFile(dbfile)
			if err != nil {
				ctx.SendChain(message.Text("读取配置文件出错了！！！"))
			}
			var temp []string
			json.Unmarshal(data, &temp)
			r := rand.Intn(len(temp))
			ctx.SendChain(message.Text(temp[r]))
		})
}
