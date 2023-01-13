package haimage

import (
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const url = "https://cdn.seovx.com/ha/?mom=302"

func init() {
	control.Register("haimage", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "古风小姐姐",
		Help:             "- 古风小姐姐图片\n",
	}).OnKeywordGroup([]string{"古风"}).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		r, err := web.RequestDataWith(web.NewDefaultClient(), url, "GET", "", web.RandUA(), nil)
		if err != nil {
			ctx.SendChain(message.Text("出错了稍后再试吧"))
			return
		}
		ctx.SendChain(message.ImageBytes(r))
	})
}
