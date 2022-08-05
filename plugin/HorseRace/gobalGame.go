package HorseRace

import (
	"math"

	"github.com/FloatTech/ZeroBot-Plugin/util"
)

func (global *globalGame) changStatus(status int) {
	global.start = status
}

func (global *globalGame) isPlayerIn(uid int64) bool {
	for _, item := range global.players {
		if item.playerUid == uid {
			return true
		}
	}
	return false
}

func (global *globalGame) isHorseIn(name string) bool {
	for _, item := range global.players {
		if item.horseName == name {
			return true
		}
	}
	return false
}

func (global *globalGame) addPlayer(horse *horse) {
	global.players = append(global.players, *horse)
}

func (global *globalGame) queryOfPlayer() int {
	return len(global.players)
}

func (global *globalGame) roundAdd() {
	global.round += 1
	for i := 0; i < len(global.players); i++ {
		global.players[i].locationAddMove = 0
		global.players[i].round = global.round
	}
}

func (global *globalGame) delBuffOvertime() {
	for i := 0; i < len(global.players); i++ {
		global.players[i].delBuffOvertime(global.round)
	}
}

// #所有马儿移动，移动计算已包含死亡/离开/止步判定
func (global *globalGame) move() {
	for i := 0; i < len(global.players); i++ {
		global.players[i].locationMove()
	}
}

// #所有马儿数据显示（须先移动)
func (global *globalGame) display() string {
	display := ""
	for i := 0; i < len(global.players); i++ {
		display += global.players[i].display()
	}
	return display
}

// #所有马儿是否死亡/离开
func (global *globalGame) isDieAll() bool {
	for i := 0; i < len(global.players); i++ {
		if global.players[i].isDie() == false && global.players[i].isAway() == false {
			return false
		}
	}
	return true
}

// Winner #所有马儿是否到终点
func (global *globalGame) Winner() string {
	winName := ""
	for i := 0; i < len(global.players); i++ {
		if global.players[i].location >= settingTrackLength {
			winName += "\n" + global.players[i].playerName
		}
	}
	return winName
}

// #事件触发
func (global *globalGame) eventStart() string {
	eventDisplay := ""
	var newDelayEvent []event
	for i := 0; i < len(global.players); i++ {
		if len(global.players[i].delayEvent) > 0 {
			for j := 0; j < len(global.players[i].delayEvent); j++ {
				if global.players[i].delayEvent[j].Rounds == global.round {
					display0 := eventMain(*global, i, global.players[i].delayEvent[j], 1) + "\n"
					if display0 != "\n" {
						eventDisplay += display0
					}
				}
				newDelayEvent = append(newDelayEvent, global.players[i].delayEvent[j])
			}
			global.players[i].delayEvent = newDelayEvent
		}
	}
	for i := 0; i < len(global.players); i++ {
		var eventInBuffX event
		for j := 0; j < len(global.players[i].selfBuff); j++ {
			if len(global.players[i].selfBuff[j].eventInBuff) != 0 {
				eventInBuff := global.players[i].selfBuff[j].eventInBuff
				eventInBuffNum := len(eventInBuff)
				eventInBuffRate := util.Rand(0, eventInBuff[eventInBuffNum-1].Probability)
				for k := 0; k < eventInBuffNum; k++ {
					if eventInBuffRate <= eventInBuff[k].Probability {
						eventInBuffX = eventInBuff[k].Other
						break
					}
				}
				display0 := eventMain(*global, i, eventInBuffX, 1) + "\n"
				if display0 != "\n" {
					eventDisplay += display0
				}
			}
		}
	}
	allEvents := len(events)
	for i := 0; i < len(global.players); i++ {
		eventId := util.Rand(0, int(math.Ceil(float64(1000*allEvents/eventRate))-1))
		if eventId < allEvents {
			display0 := eventMain(*global, i, events[eventId], 1) + "\n"
			if display0 != "\n" {
				eventDisplay += display0
			}
		}
	}
	return eventDisplay
}

// #事件唯一码查询
func (global *globalGame) isRaceOnlyKeyIn(key string) bool {
	for _, keys := range global.raceOnlyKeys {
		if keys == key {
			return true
		}
	}
	return false
}

// #事件唯一码增加
func (global *globalGame) addRaceOnlyKey(key string) {
	global.raceOnlyKeys = append(global.raceOnlyKeys, key)
}
