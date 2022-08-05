package HorseRace

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/FloatTech/ZeroBot-Plugin/util"
)

func eventMain(race globalGame, mover int, ev event, eventDelayKey int) string {
	//	eventName := ev.EventName
	if race.players[mover].isDie() && eventDelayKey == 0 {
		return ""
	} else if race.players[mover].isAway() && eventDelayKey == 0 {
		return ""
	} else if race.players[mover].findBuff("vertigo") && eventDelayKey == 0 {
		return ""
	}
	describe := ev.Describe
	target := ev.Target
	var targets []int
	var targetName1 string
	targetName0 := race.players[mover].horseName
	for i := 0; i < len(race.players); i++ {
		targets = append(targets, i)
	}
	switch target {
	case 0: // 目标自己
		targets = []int{mover}
		targetName1 = race.players[mover].horseName
	case 1: //# 1为随机选择一个非自己的目标（即<1>）
		targets = append(targets[:mover], targets[mover+1:]...)
		util.Shuffle(targets)
		targetName1 = race.players[targets[0]].horseName
	case 2:
		targetName1 = "所有马儿"
	case 3:
		targets = append(targets[:mover], targets[mover+1:]...)
		targetName1 = "其他所有马儿"
	case 4: // 随机一个目标
		util.Shuffle(targets)
		targets = []int{targets[0]}
		targetName1 = race.players[targets[0]].horseName
	case 5:
		targets = append(targets[:mover], targets[mover+1:]...)
		util.Shuffle(targets)
		targets = []int{targets[0], mover}
		targetName1 = race.players[targets[0]].horseName
	}
	targetIsBuff := ev.TargetIsBuff
	if targetIsBuff != "" {
		var a []int
		for _, i := range targets {
			if race.players[i].findBuff(targetIsBuff) {
				a = append(a, i)
			}
		}
		targets = a
	}
	targetNoBuff := ev.TargetNoBuff
	if targetNoBuff != "" {
		var a []int
		for _, i := range targets {
			if !race.players[i].findBuff(targetNoBuff) {
				a = append(a, i)
			}
		}
		targets = a
	}
	if len(targets) == 0 {
		return ""
	}
	describe = strings.Replace(describe, "<0>", targetName0, -1)
	describe = strings.Replace(describe, "<1>", targetName1, -1)
	if ev.Live == 1 {
		eventLive(race, targets)
	}
	if ev.Move != 0 {
		eventMove(race, targets, ev.Move)
	}
	if ev.TrackRandomLocation == 1 {
		eventTrackRandomLocation(race, targets)
	}
	if ev.BuffTimeAdd != 0 {
		eventBuffTimeAdd(race, targets, ev.BuffTimeAdd)
	}
	if ev.DelBuff != "" {
		eventDelBuffs(race, targets, ev.DelBuff)
	}
	if ev.TrackExchangeLocation == 1 && (target == 1 || target == 6) {
		eventTrackExchangeLocation(race, mover, targets[0])
	}
	var eventOnce event
	if len(ev.RandomEventOnce) != 0 {
		randomEventOnce := ev.RandomEventOnce
		randomEventOnceNum := len(randomEventOnce)
		for _, i := range targets {
			for j := 0; j < randomEventOnceNum; j++ {
				// TODO 可能需要优化
				randomEventOnceRate := util.Rand(0, randomEventOnce[randomEventOnceNum-1].Probability)
				if randomEventOnceRate <= randomEventOnce[j].Probability {
					eventOnce = randomEventOnce[j].Other
					break
				}
			}
			describe += eventMain(race, i, eventOnce, 1)
		}
	}
	if ev.Die == 1 {
		eventDie(race, targets, ev.DieName)
	}
	if ev.Away == 1 {
		eventAway(race, targets, ev.AwayName)
	}
	if ev.Rounds > 0 {
		rounds := ev.Rounds
		buffName := ev.Name
		var buffs = []string{""}
		moveMax := ev.MoveMax
		moveMin := ev.MoveMin
		if ev.LocateLock == 1 {
			buffs = append(buffs, "locate_lock")
		}
		if ev.Vertigo == 1 {
			buffs = append(buffs, "locate_lock", "vertigo")
		}
		if ev.Hiding == 1 {
			buffs = append(buffs, "hiding")
		}
		if len(ev.OtherBuff) != 0 {
			buffs = append(buffs, ev.OtherBuff...)
		}
		eventInBuff := ev.RandomEvent
		for _, buf := range buffs {
			eventAddBuff(race, targets, buff{
				buffName:    buffName,
				roundStart:  race.round + 1,
				roundEnd:    race.round + rounds,
				moveMin:     moveMin,
				moveMax:     moveMax,
				eventInBuff: eventInBuff,
				buffTag:     buf,
			})
		}
	}
	if len(ev.DelayEvent) != 0 {
		eventDelayRounds := ev.DelayEvent[0].Round
		if eventDelayRounds > 1 {
			eventDelay := ev.DelayEvent[0].Other
			for _, i := range targets {
				race.players[i].delayEvent = append(race.players[i].delayEvent, eventDelay)
			}
		}
	}
	if len(ev.DelayEventSelf) != 0 {
		eventDelayRoundsSelf := ev.DelayEventSelf[0].Round
		if eventDelayRoundsSelf > 1 {
			race.players[mover].delayEvent = append(race.players[mover].delayEvent, ev.DelayEventSelf[0].Other)
		}
	}
	if ev.AnotherEvent != nil {
		for _, i := range targets {
			describe += eventMain(race, i, *ev.AnotherEvent, 1)
		}
	}
	if ev.AnotherEventSelf != nil {
		describe += eventMain(race, mover, *ev.AnotherEventSelf, 1)
	}
	if ev.AddHorse != nil {
		addHorseEvent := ev.AddHorse
		race.addPlayer(&horse{
			horseName:  addHorseEvent.Horsename,
			playerName: addHorseEvent.Owner,
			playerUid:  addHorseEvent.Uid,
			location:   addHorseEvent.Location,
			round:      race.round,
		})
	}
	if ev.ReplaceHorse != nil {
		replaceHorse := ev.AddHorse
		if target == 0 || target == 1 || target == 4 || target == 6 {
			race.players[targets[0]].replaceHorseEx(replaceHorse.Horsename, replaceHorse.Uid, replaceHorse.Owner)
		}
	}
	return describe
}

func eventLive(race globalGame, targets []int) string {
	var msg string
	for _, i := range targets {
		if race.players[i].findBuff("die") {
			race.players[i].delBuff("die")
		}
		msg += fmt.Sprintf("%s复活了", race.players[i].horseName)
	}
	log.Println(msg)
	return msg
}

func eventMove(race globalGame, targets []int, move int) string {
	var msg string
	for _, i := range targets {
		race.players[i].locationMoveEvent(move)
		msg += fmt.Sprintf("%s移动了%d", race.players[i].horseName, move)
	}
	log.Println(msg)
	return msg
}

// func eventTrackToLocation(race globalGame, targets []int, moveTo int) string {
//	var msg string
//	for _, i := range targets {
//		race.players[i].locationMoveToEvent(moveTo)
//		msg += fmt.Sprintf("%s移动到了指定位置%d", race.players[i].horseName, moveTo)
//	}
//	return msg
//}

// #移动对象至随机位置
func eventTrackRandomLocation(race globalGame, targets []int) string {
	var msg string
	for _, i := range targets {
		moveTo := util.Rand(settingRandomMinLength, settingRandomMaxLength)
		race.players[i].locationMoveToEvent(moveTo)
		msg += fmt.Sprintf("%s移动到了随机位置%d", race.players[i].horseName, moveTo)
	}
	log.Println(msg)
	return msg
}

func eventBuffTimeAdd(race globalGame, targets []int, round int) string {
	var msg string
	for _, i := range targets {
		race.players[i].buffAddTime(round)
		msg += fmt.Sprintf("%s 的buff事件增加了%d", race.players[i].horseName, round)
	}
	log.Println(msg)
	return msg
}

func eventDelBuffs(race globalGame, targets []int, buffName string) string {
	var msg string
	for _, i := range targets {
		//	for _, j := range buffName {
		race.players[i].delBuff(buffName)
		msg += fmt.Sprintf("%s 移除了buff：%s", race.players[i].horseName, buffName)
		//	}
	}
	log.Println(msg)
	return msg
}

func eventTrackExchangeLocation(race globalGame, a, b int) string {
	var msg string
	x := race.players[a].location
	race.players[a].locationMoveToEvent(race.players[b].location)
	race.players[a].locationMoveToEvent(x)
	msg += fmt.Sprintf("%s 和 %s互换位置", race.players[a].horseName, race.players[b].horseName)
	log.Println(msg)
	return msg
}

func eventDie(race globalGame, targets []int, buffName string) string {
	var msg string
	for _, i := range targets {
		if !race.players[i].findBuff("die") {
			race.players[i].addBuff(buff{
				buffName:    buffName,
				roundStart:  race.round + 1,
				roundEnd:    999,
				moveMin:     0,
				moveMax:     0,
				eventInBuff: nil,
				buffTag:     "die",
			})
			msg += fmt.Sprintf("死亡事件判定 %s 死了", race.players[i].horseName)
		}
	}
	log.Println(msg)
	return msg
}

func eventAway(race globalGame, targets []int, buffName string) string {
	var msg string
	for _, i := range targets {
		if !race.players[i].findBuff("away") {
			race.players[i].addBuff(buff{
				buffName:    buffName,
				roundStart:  race.round + 1,
				roundEnd:    999,
				moveMin:     0,
				moveMax:     0,
				eventInBuff: nil,
				buffTag:     "away",
			})
			msg += fmt.Sprintf(" %s 离开了跑道", race.players[i].horseName)
		}
	}
	log.Println(msg)
	return msg
}

func eventAddBuff(race globalGame, targets []int, buffs buff) string {
	var msg string
	for _, i := range targets {
		if !race.players[i].findBuff(buffs.buffName) {
			msg += fmt.Sprintf(" %s 增加了 buff: %s ,第 %d~%d回合", race.players[i].horseName, buffs.buffName, buffs.roundStart, buffs.roundEnd)
		} else {
			race.players[i].delBuff(buffs.buffName)
		}
		race.players[i].addBuff(buffs)
	}
	log.Println(msg)
	return msg
}
