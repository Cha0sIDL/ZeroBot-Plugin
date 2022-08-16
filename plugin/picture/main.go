package picture

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/picture"
)

const (
	servicename = "picture"
)

func init() {
	engine := control.Register(servicename, &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "- xxx表情\n" + "可能会偶尔抽风",
	})
	engine.OnSuffixGroup([]string{"表情", "表情包"}).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			arg := ctx.State["args"].(string)
			if arg == "" {
				return
			}
			url := picture.GetPicture(arg)
			if len(url) == 0 {
				ctx.SendChain(message.Text("出错了稍后再试试吧~"))
				return
			}
			ctx.SendChain(message.Image(url))
		})
}
