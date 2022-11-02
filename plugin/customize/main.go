package customize

import (
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/fumiama/unibase2n"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/ZeroBot-Plugin/util"

	"github.com/FloatTech/floatbox/process"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := control.Register("custom", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "自定义插件集合\n" +
			" - /kill\n" +
			" - /发送公告\n" +
			" - /restart\n" +
			" - @bot给主人留言<内容>",
	})
	engine.OnCommandGroup([]string{"pause", "kill"}, zero.OnlyToMe, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("正在自爆...(我会想你们的)"))
			time.Sleep(time.Second * 5)
			ctx.SendChain(message.Face(55))
			os.Exit(0)
		})
	engine.OnCommandGroup([]string{"restart"}, zero.OnlyToMe, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			os.Exit(1) // systemd会把服务拉起来
		})
	engine.OnCommand("发送公告", zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			next := zero.NewFutureEvent("message", 1, false, zero.CheckUser(ctx.Event.UserID), ctx.CheckSession())
			recv, stop := next.Repeat()
			defer stop()
			ctx.SendChain(message.Text("请输入公告内容"))
			var step int
			var origin string
			for {
				select {
				case <-time.After(time.Second * 120):
					ctx.SendChain(message.Text("时间太久啦！不发了！"))
					return
				case c := <-recv:
					switch step {
					case 0:
						origin = "来自开发者的信息：\n" + c.Event.RawMessage + "\n--------------------\n" + unibase2n.BaseRune.EncodeString(util.RandStr(rand.Intn(20)))
						ctx.SendChain(message.Text("请输入\"确定\"或者\"取消\"来决定是否发送此公告"))
						step++
					case 1:
						msg := c.Event.Message.ExtractPlainText()
						if msg != "确定" && msg != "取消" {
							ctx.SendChain(message.Text("请输入\"确定\"或者\"取消\"哟"))
							continue
						}
						if msg == "确定" {
							ctx.SendChain(message.Text("正在发送..."))
							var fail []int64
							zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
								grpList := ctx.GetGroupList().Array()
								for _, g := range grpList {
									time.Sleep(time.Second + time.Second*time.Duration(rand.Intn(20)))
									gid := g.Get("group_id").Int()
									if id := ctx.SendGroupMessage(gid, origin); id == 0 {
										fail = append(fail, gid)
									}
									process.SleepAbout1sTo2s()
								}
								return true
							})
							if len(fail) == 0 {
								ctx.SendChain(message.Text("发送成功"))
							} else {
								ctx.SendChain(message.Text("检测到公告发送失败,群号为:", util.PrettyPrint(fail)))
							}
							return
						}
						ctx.SendChain(message.Text("已经取消发送了哟~"))
						return
					}
				}
			}
		})
	engine.OnRegex(`给主人留言.*?(.*)`, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			su := zero.BotConfig.SuperUsers[0]
			now := time.Unix(ctx.Event.Time, 0).Format("2006-01-02 15:04:05")
			uid := ctx.Event.UserID
			gid := ctx.Event.GroupID
			username := ctx.CardOrNickName(uid)
			botid := ctx.Event.SelfID
			botname := zero.BotConfig.NickName[0]
			rawmsg := ctx.State["regex_matched"].([]string)[1]
			rawmsg = message.UnescapeCQCodeText(rawmsg)
			msg := make(message.Message, 10)
			msg = append(msg, message.CustomNode(botname, botid, "有人给你留言啦！\n在"+now))
			if gid != 0 {
				groupname := ctx.GetGroupInfo(gid, true).Name
				msg = append(msg, message.CustomNode(botname, botid, "来自群聊:["+groupname+"]("+strconv.FormatInt(gid, 10)+")\n来自群成员:["+username+"]("+strconv.FormatInt(uid, 10)+")\n以下是留言内容"))
			} else {
				msg = append(msg, message.CustomNode(botname, botid, "来自私聊:["+username+"]("+strconv.FormatInt(uid, 10)+")\n以下是留言内容:"))
			}
			msg = append(msg, message.CustomNode(username, uid, rawmsg))
			ctx.SendPrivateForwardMessage(su, msg)
		})
	engine.OnMessage(func(ctx *zero.Ctx) bool {
		msg := ctx.Event.Message
		if len(msg) > 0 {
			if msg[0].Type != "reply" {
				return false
			}
			for _, elem := range msg {
				if elem.Type == "text" {
					text := elem.Data["text"]
					text = strings.ReplaceAll(text, " ", "")
					text = strings.ReplaceAll(text, "\r", "")
					text = strings.ReplaceAll(text, "\n", "")
					if text == "撤回" {
						return true
					}
				}
			}
		}
		return false
	}, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(
		func(ctx *zero.Ctx) {
			ctx.DeleteMessage(message.NewMessageIDFromString(ctx.Event.Message[0].Data["id"]))
			ctx.DeleteMessage(message.NewMessageIDFromInteger(ctx.Event.MessageID.(int64)))
			return
		})
}
