// Package horserace  赛马小游戏
package horserace

import (
	"fmt"
	"time"

	"github.com/golang-module/carbon/v2"

	cmap "github.com/orcaman/concurrent-map/v2"
	log "github.com/sirupsen/logrus"

	"github.com/FloatTech/ZeroBot-Plugin/util"

	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/img/text"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// type rwRace struct {
//	sync.Mutex
//	race map[int64]*globalGame
//}

var race = cmap.New[*globalGame]()

func init() {
	engine := control.Register("horseRace", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "群内赛马小游戏",
		Help: "赛马创建\n" +
			"赛马加入xxx xxx为马儿的名字\n" +
			"赛马开始\n" +
			"赛马事件重载\n" +
			"赛马暂停\n" +
			"赛马继续\n",
		PrivateDataFolder: "hr",
	})
	cacheDir := file.BOTPATH + "/" + engine.DataFolder()
	initConfig(cacheDir)
	engine.OnFullMatch("赛马创建", zero.OnlyGroup).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			strKey := util.Interface2String(ctx.Event.GroupID)
			if val, ok := race.Get(strKey); ok {
				if val.start == 1 {
					ctx.SendChain(message.Text("一场赛马正在进行中..."))
				} else {
					race.Set(strKey, new(globalGame))
					ctx.SendChain(message.Text("创建赛马成功"))
				}
			} else {
				race.Set(strKey, new(globalGame))
				ctx.SendChain(message.Text("创建赛马成功"))
			}
		})
	engine.OnPrefix("赛马加入", zero.OnlyGroup).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			strKey := util.Interface2String(ctx.Event.GroupID)
			horseName := ctx.State["args"].(string)
			if val, ok := race.Get(strKey); !ok {
				ctx.SendChain(message.Text("赛马活动未开始，请输入“赛马创建”开场"))
			} else {
				switch {
				case len(val.players) >= maxPlayer:
					ctx.SendChain(message.Text("> 加入失败\n> 原因:赛马场就那么大，满了满了！"))
				case len(horseName) == 0:
					ctx.SendChain(message.Text("> 加入失败\n> 原因:您没有输入马儿的名字"))
				case val.isPlayerIn(ctx.Event.UserID):
					ctx.SendChain(message.Text(">  加入失败\n> 原因:您已经加入了赛马场!"))
				case val.isHorseIn(horseName):
					ctx.SendChain(message.Text("> 加入失败\n> 原因:有一匹同名的马儿加入了赛马场!"))
				case len([]rune(horseName)) > nameMaxLen:
					ctx.SendChain(message.Text(fmt.Sprintf("> 加入失败\n> 原因:马儿名字太长了\n> 不可以超过%d个字哦", nameMaxLen)))
				case string(horseName[:1]) == ".": //nolint:unconvert
					ctx.SendChain(message.Text("> 加入失败\n> 原因:马儿名字不可以用“.”开头哦"))
				case string(horseName[len(horseName)-1:]) == ".": //nolint:unconvert
					ctx.SendChain(message.Text("> 加入失败\n> 原因:马儿名字不可以用“.”结尾哦"))
				default:
					val.addPlayer(&horse{
						horseName:  horseName,
						playerName: ctx.Event.Sender.NickName,
						playerUID:  ctx.Event.UserID,
					})
					ctx.SendChain(message.Text(fmt.Sprintf("> 加入赛马成功\n> 赌上马儿性命的一战即将开始!\n> 赛马场位置:%d/%d", val.queryOfPlayer(), maxPlayer)))
				}
			}
		})
	engine.OnFullMatch("赛马开始", zero.OnlyGroup).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			strKey := util.Interface2String(ctx.Event.GroupID)
			val, ok := race.Get(strKey)
			if !ok {
				ctx.SendChain(message.Text("赛马活动未开始，请输入“赛马创建”开场"))
				return
			}
			if len(val.players) < minPlayer {
				ctx.SendChain(message.Text(fmt.Sprintf("> 开始失败\n> 原因:赛马开局需要最少%d人参与", minPlayer)))
				return
			} else if val.start == 1 {
				ctx.SendChain(message.Text("一场赛马正在进行中..."))
				return
			}
			val.time = carbon.Now().Timestamp()
			val.changStatus(1)
			for {
				v, _ := race.Get(strKey)
				switch {
				case carbon.Now().Timestamp()-v.time >= settingOverTime:
					race.Remove(strKey)
					ctx.SendChain(message.Text("赛马超时,已结束..."))
					return
				case v.start == 1:
					display := ""
					val.roundAdd()
					val.delBuffOvertime()
					display += val.eventStart()
					val.move()
					display += val.display()
					data, _ := text.RenderToBase64(display, text.FontFile, 250, 20)
					ctx.SendChain(message.Image("base64://" + helper.BytesToString(data)))
					if val.isDieAll() {
						race.Remove(strKey)
						ctx.SendChain(message.Text("比赛已结束，鉴定为无马生还"))
						return
					}
					winner := val.Winner()
					if len(winner) != 0 {
						ctx.SendChain(message.Text(fmt.Sprintf("> 比赛结束\n> %s正在为您生成战报...", zero.BotConfig.NickName[0])))
						time.Sleep(time.Second * 2)
						race.Remove(strKey)
						ctx.SendChain(message.Text("比赛已结束，胜者为：" + winner))
						return
					}
					time.Sleep(time.Second * 5)
				default:
					time.Sleep(time.Second * 1)
				}
			}
		})
	engine.OnFullMatch("赛马暂停", zero.OnlyGroup, zero.AdminPermission).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			strKey := util.Interface2String(ctx.Event.GroupID)
			val, ok := race.Get(strKey)
			if !ok {
				ctx.SendChain(message.Text("赛马活动未开始，请输入“赛马创建”开场"))
				return
			}
			val.changStatus(0)
			race.Set(strKey, val)
			ctx.SendChain(message.Text("赛马已暂停"))
		})
	engine.OnFullMatch("赛马继续", zero.OnlyGroup, zero.AdminPermission).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			strKey := util.Interface2String(ctx.Event.GroupID)
			val, ok := race.Get(strKey)
			if !ok {
				ctx.SendChain(message.Text("赛马活动未开始，请输入“赛马创建”开场"))
				return
			}
			val.changStatus(1)
			race.Set(strKey, val)
			ctx.SendChain(message.Text("赛马已继续"))
		})
	engine.OnFullMatch("赛马事件重载", zero.SuperUserPermission).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			initConfig(cacheDir)
			ctx.SendChain(message.Text("事件重载成功共加载：", len(events), "条事件"))
		})
	engine.OnFullMatch("测试赛马", zero.OnlyGroup).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			var players []horse
			players = append(players, horse{
				horseName:  "test1",
				playerName: "test1",
				playerUID:  123456,
			}, horse{
				horseName:  "test2",
				playerName: "test2",
				playerUID:  654321,
			},
				horse{
					horseName:  "test3",
					playerName: "test3",
					playerUID:  65432,
				},
				horse{
					horseName:  "test4",
					playerName: "test4",
					playerUID:  6543,
				})
			race.Set("123456", &globalGame{
				players:      players,
				round:        0,
				start:        1,
				time:         0,
				raceOnlyKeys: nil,
				events:       nil,
			})
			for {
				val, _ := race.Get("123456")
				if val.start == 1 {
					display := ""
					val.roundAdd()
					val.delBuffOvertime()
					display += val.eventStart()
					val.move()
					display += val.display()
					log.Println(display)
					if val.isDieAll() {
						race.Remove("123456")
						fmt.Println("die")
						return
					}
					winner := val.Winner()
					if len(winner) != 0 {
						time.Sleep(time.Second * 2)
						race.Remove("123456")
						fmt.Println("winner")
						return
					}
					time.Sleep(time.Second * 4)
				}
			}
		})
}
