package active

import (
	"errors"
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
)

const (
	serviceName = "active"
	pictureUrl  = "https://doutu.lccyy.com/doutu/items?"
)

func init() {
	en := control.Register(serviceName, &control.Options{
		DisableOnDefault: false,
		Help: "自动插话\n" +
			"- 设置活跃度 xx\n" +
			"- 查询活跃度",
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
	en.OnMessage(func(ctx *zero.Ctx) bool {
		return util.Rand(1, 100) < getActive(ctx) && zero.OnlyGroup(ctx) && util.Ignore(ctx)
	}).SetBlock(false).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if zero.HasPicture(ctx) {
				for _, elem := range ctx.Event.Message {
					if elem.Type == "image" {
						ocrTags := make([]string, 0)
						ocrResult := ctx.OCRImage(elem.Data["file"]).Get("texts.#.text").Array()
						if len(ocrResult) == 0 {
							return
						}
						for _, text := range ocrResult {
							ocrTags = append(ocrTags, text.Str)
						}
						text := fmt.Sprintf("%s", strings.Join(ocrTags, ""))
						url := pictureUrl + fmt.Sprintf("pageNum=%d&pageSize=%d&keyword=", 1, util.Rand(1, 100)) + url.QueryEscape(text)
						data, err := web.RequestDataWith(web.NewDefaultClient(), url, "GET", "", web.RandUA())
						if err != nil {
							log.Errorln("Active Err :", err)
							return
						}
						Items := gjson.Get(binary.BytesToString(data), "items").Array()
						ctx.SendChain(message.Image(Items[rand.Intn(len(Items))].Get("url").String()))
					}
				}
			} else {
				msg := ctx.ExtractPlainText()
				t := []string{"青云客", "腾讯", "小爱"}
				util.Shuffle(t)
				r := NewAIReply(t[0])
				ctx.SendChain(message.Text(r.TalkPlain(msg, zero.BotConfig.NickName[0])))
			}
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
