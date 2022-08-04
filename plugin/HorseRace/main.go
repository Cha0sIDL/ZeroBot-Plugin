package HorseRace

import (
	"fmt"
	"sync"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/file"
	"github.com/golang-module/carbon/v2"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type rwRace struct {
	sync.Mutex
	race map[int64]*globalGame
}

var race = rwRace{
	race: make(map[int64]*globalGame),
}

func init() {
	engine := control.Register("horseRace", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "群内赛马小游戏\n",
		PrivateDataFolder: "hr",
	})
	cacheDir := file.BOTPATH + "/" + engine.DataFolder()
	initConfig(cacheDir)
	engine.OnFullMatch("赛马创建", zero.OnlyGroup).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			race.Lock()
			defer race.Unlock()
			if val, ok := race.race[ctx.Event.GroupID]; ok {
				if val.start == 0 && carbon.Now().Timestamp()-val.time < settingOverTime {
					ctx.SendChain(message.Text(fmt.Sprintf("> 创建赛马比赛失败!\n> 原因:赫尔正在打扫赛马场...\n> 解决方案:等赫尔打扫完...\n> 可以在%d秒后开始下局比赛", settingOverTime-(carbon.Now().Timestamp()-val.time))))
				} else if val.start == 1 {
					ctx.SendChain(message.Text("一场赛马正在进行中..."))
				} else {
					race.race[ctx.Event.GroupID] = new(globalGame)
					ctx.SendChain(message.Text("创建赛马成功"))
				}
			} else {
				race.race[ctx.Event.GroupID] = new(globalGame)
				ctx.SendChain(message.Text("创建赛马成功"))
			}
		})
	engine.OnPrefix("赛马加入", zero.OnlyGroup).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			race.Lock()
			defer race.Unlock()
			horseName := ctx.State["args"].(string)
			if val, ok := race.race[ctx.Event.GroupID]; !ok {
				ctx.SendChain(message.Text("赛马活动未开始，请输入“赛马创建”开场"))
			} else {
				if len(val.players) >= maxPlayer {
					ctx.SendChain(message.Text("> 加入失败\n> 原因:赛马场就那么大，满了满了！"))
				} else if len(horseName) <= 0 {
					ctx.SendChain(message.Text("> 加入失败\n> 原因:您没有输入马儿的名字"))
				} else if val.isPlayerIn(ctx.Event.UserID) {
					ctx.SendChain(message.Text(">  加入失败\n> 原因:您已经加入了赛马场!"))
				} else if val.isHorseIn(horseName) {
					ctx.SendChain(message.Text("> 加入失败\n> 原因:有一匹同名的马儿加入了赛马场!"))
				} else if len([]rune(horseName)) > nameMaxLen {
					ctx.SendChain(message.Text(fmt.Sprintf("> 加入失败\n> 原因:马儿名字太长了\n> 不可以超过%d个字哦", nameMaxLen)))
				} else if string(horseName[:1]) == "." {
					ctx.SendChain(message.Text("> 加入失败\n> 原因:马儿名字不可以用“.”开头哦"))
				} else if string(horseName[len(horseName)-1:]) == "." {
					ctx.SendChain(message.Text("> 加入失败\n> 原因:马儿名字不可以用“.”结尾哦"))
				} else {
					val.addPlayer(&horse{
						horseName:  horseName,
						playerName: ctx.Event.Sender.NickName,
						playerUid:  ctx.Event.UserID,
					})
					ctx.SendChain(message.Text(fmt.Sprintf("> 加入赛马成功\n> 赌上马儿性命的一战即将开始!\n> 赛马场位置:%d/%d", val.queryOfPlayer(), maxPlayer)))
				}
			}
		})
	engine.OnFullMatch("赛马开始", zero.OnlyGroup).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			race.Lock()
			defer race.Unlock()
			if val, ok := race.race[ctx.Event.GroupID]; !ok {
				ctx.SendChain(message.Text("赛马活动未开始，请输入“赛马创建”开场"))
			} else {
				if len(val.players) < minPlayer {
					ctx.SendChain(message.Text(fmt.Sprintf("> 开始失败\n> 原因:赛马开局需要最少%d人参与", minPlayer)))
				} else if val.start == 1 {
					ctx.SendChain(message.Text("一场赛马正在进行中..."))
					return
				}
				val.changStatus(1)
				for race.race[ctx.Event.GroupID].start == 1 {
					display := ""
					val.roundAdd()
					display += val.eventStart()
					val.move()
					display += val.display()
					// ctx.SendChain(message.Text(display))
					fmt.Println(display)
					if val.isDieAll() {
						delete(race.race, ctx.Event.GroupID)
						ctx.SendChain(message.Text("比赛已结束，鉴定为无马生还"))
						return
					}
					winner := val.Winner()
					if len(winner) != 0 {
						ctx.SendChain(message.Text("> 比赛结束\n> 赫尔正在为您生成战报..."))
						time.Sleep(time.Second * 2)
						delete(race.race, ctx.Event.GroupID)
						ctx.SendChain(message.Text("比赛已结束，胜者为：" + winner))
						return
					}
					time.Sleep(time.Second * 4)
				}
			}
		})
	engine.OnFullMatch("测试赛马", zero.OnlyGroup).SetBlock(true).Handle(
		func(ctx *zero.Ctx) {
			race.Lock()
			defer race.Unlock()
			var players []horse
			players = append(players, horse{
				horseName:  "test1",
				playerName: "test1",
				playerUid:  123456,
			}, horse{
				horseName:  "test2",
				playerName: "test2",
				playerUid:  654321,
			})
			race.race[123456] = &globalGame{
				players:      players,
				round:        0,
				start:        1,
				time:         0,
				raceOnlyKeys: nil,
				events:       nil,
			}
			val := race.race[123456]
			for race.race[123456].start == 1 {
				display := ""
				val.roundAdd()
				display += val.eventStart()
				val.move()
				display += val.display()
				log.Println(display)
				if val.isDieAll() {
					delete(race.race, ctx.Event.GroupID)
					return
				}
				winner := val.Winner()
				if len(winner) != 0 {
					time.Sleep(time.Second * 2)
					delete(race.race, ctx.Event.GroupID)
					return
				}
				time.Sleep(time.Second * 4)
			}
		})
}
