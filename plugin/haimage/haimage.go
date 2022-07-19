package haimage

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const url = "https://cdn.seovx.com/ha/?mom=302"

func init() {
	control.Register("haimage", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "- 古风小姐姐图片\n",
	}).OnKeywordGroup([]string{"古风"}).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		r, err := web.RequestDataWith(web.NewDefaultClient(), url, "GET", "", web.RandUA())
		if err != nil {
			ctx.SendChain(message.Text("出错了稍后再试吧"))
			return
		}
		ctx.SendChain(message.ImageBytes(r))
	})
}
