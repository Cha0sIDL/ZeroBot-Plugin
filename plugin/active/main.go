package active

import (
	"errors"
	"github.com/FloatTech/AnimeAPI/aireply"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/control/order"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
)

const serviceName = "active"

func init() {
	en := control.Register(serviceName, order.AcquirePrio(), &control.Options{
		DisableOnDefault: false,
		Help: "自动插话\n" +
			"- 设置活跃度 xx\n" +
			"- 查询活跃度",
	})
	en.OnMessage(func(ctx *zero.Ctx) bool {
		return util.Rand(1, 100) < getActive(ctx) && zero.OnlyGroup(ctx)
	}).SetBlock(false).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			msg := ctx.ExtractPlainText()
			r := aireply.NewAIReply("青云客")
			ctx.SendChain(message.Text(r.TalkPlain(msg, zero.BotConfig.NickName[0])))
		})
	en.OnRegex(`设置活跃度(\d+)`, zero.SuperUserPermission, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			arg := ctx.State["regex_matched"].([]string)[1]
			active, _ := strconv.Atoi(arg)
			if active > 100 || active < 0 {
				ctx.SendChain(message.Text("请输入1-100内的活跃值"))
				return
			}
			err := setActive(ctx, active)
			if err != nil {
				ctx.SendChain(message.Text("Err :", err))
			}
			ctx.SendChain(message.Text("设置成功"))
		})
	en.OnFullMatch("查询活跃度", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			active := getActive(ctx)
			ctx.SendChain(message.Text("本群当前活跃度为:", active))
		})
}

func setActive(ctx *zero.Ctx, active int) error {
	gid := ctx.Event.GroupID
	if gid == 0 {
		gid = -ctx.Event.UserID
	}
	var ok bool
	m, ok := control.Lookup(serviceName)
	if !ok {
		return errors.New("no such plugin")
	}
	return m.SetData(gid, int64(active))
}

func getActive(ctx *zero.Ctx) (active int) {
	gid := ctx.Event.GroupID
	if gid == 0 {
		gid = -ctx.Event.UserID
	}
	m, ok := control.Lookup(serviceName)
	if ok {
		return int(m.GetData(gid))
	}
	return 0
}
