// Package logo 随机头像
package logo

import (
	"fmt"
	"math/rand"

	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	servicename = "logo"
	logoURL     = "https://avatars.dicebear.com/api/"
)

//https://avatars.dicebear.com/

func init() {
	engine := control.Register(servicename, &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "- 头像\n",
	})
	engine.OnFullMatch("头像").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			url := logoURL + fmt.Sprintf("%s/%d.png", getRandArg(), rand.Intn(10000))
			ctx.SendChain(message.Image(url))
		})
}

func getRandArg() string {
	all := []string{"male", "female", "human", "identicon", "initials", "bottts", "avataaars", "jdenticon", "gridy", "micah"}
	return all[rand.Intn(len(all))]
}
