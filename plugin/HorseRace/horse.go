package HorseRace

import (
	"github.com/FloatTech/ZeroBot-Plugin/util"
)

// #=====替换为其他马,指定数据（用于天灾马系列事件）
func (h *horse) replaceHorseEx(horseName string, uid int64, playerName string) {
	h.horseName = horseName
	h.playerUid = uid
	h.playerName = playerName
	h.selfBuff = nil
	h.delayEvent = nil
	h.locationAdd = 0
	h.locationAddMove = 0
}

func (h *horse) addBuff(buffs buff) {
	min := buffs.moveMin
	max := buffs.moveMax
	if min > max {
		buffs.moveMax = buffs.moveMin
	}
	h.selfBuff = append(h.selfBuff, buffs)
}

// #=====马儿指定buff移除：
func (h *horse) delBuff(buffKey string) {
	var buffs []buff
	for i := 0; i < len(h.selfBuff); i++ {
		if h.selfBuff[i].buffTag == buffKey {
			continue
		}
		buffs = append(buffs, h.selfBuff[i])
	}
	h.selfBuff = buffs
}

func (h *horse) findBuff(buffKey string) bool {
	for i := 0; i < len(h.selfBuff); i++ {
		if h.selfBuff[i].buffTag == buffKey {
			return true
		}
	}
	return false
}

func (h *horse) delBuffOvertime(round int) {
	var buffs []buff
	for i := 0; i < len(h.selfBuff); i++ {
		if h.selfBuff[i].roundEnd < round {
			continue
		}
		buffs = append(buffs, h.selfBuff[i])
	}
	h.selfBuff = buffs
}

func (h *horse) buffAddTime(round int) {
	for i := 0; i < len(h.selfBuff); i++ {
		buffs := h.selfBuff[i]
		buffs.roundEnd += round
		h.selfBuff[i] = buffs
	}
}

// [buff_name, round_start, round_end, move_min, move_max, event_in_buff]

func (h *horse) isStop() bool {
	for i := 0; i < len(h.selfBuff); i++ {
		if h.selfBuff[i].buffTag == "locate_lock" && h.selfBuff[i].roundStart <= h.round {
			return true
		}
	}
	return false
}

func (h *horse) isAway() bool {
	for i := 0; i < len(h.selfBuff); i++ {
		if h.selfBuff[i].buffTag == "away" && h.selfBuff[i].roundStart <= h.round {
			return true
		}
	}
	return false
}

func (h *horse) isDie() bool {
	for i := 0; i < len(h.selfBuff); i++ {
		if h.selfBuff[i].buffTag == "die" && h.selfBuff[i].roundStart <= h.round {
			return true
		}
	}
	return false
}

func (h *horse) locationMoveEvent(move int) {
	h.locationAddMove = h.locationAddMove + move
}

func (h *horse) locationMoveToEvent(moveTo int) {
	h.locationAddMove = h.locationAddMove + moveTo - h.location
}

func (h *horse) locationMove() {
	if h.location != settingTrackLength {
		h.locationAdd = h.move() + h.locationAddMove
		h.location = h.location + h.locationAdd
		if h.location > settingTrackLength {
			h.locationAdd -= h.location - settingTrackLength
			h.location = settingTrackLength
		}
		if h.location < 0 {
			h.locationAdd -= h.location
			h.location = 0
		}
	}
}

func (h *horse) move() int {
	if h.isStop() || h.isDie() || h.isAway() {
		return 0
	}
	moveMin := 0
	moveMax := 0
	for i := 0; i < len(h.selfBuff); i++ {
		if h.selfBuff[i].roundStart <= h.round && h.selfBuff[i].roundEnd >= h.round {
			moveMin += h.selfBuff[i].moveMin
			moveMax += h.selfBuff[i].moveMax
		}
	}
	return util.Rand(moveMin+baseMoveMin, moveMax+baseMoveMax)
}

func (h *horse) display() string {
	dis := ""
	if !h.findBuff("hiding") {
		if h.locationAdd < 0 {
			dis += "[" + util.Interface2String(h.locationAdd) + "]"
		} else {
			dis += "[+" + util.Interface2String(h.locationAdd) + "]"
		}
		for i := 0; i < settingTrackLength-h.location; i++ {
			dis += "."
		}
		dis += h.horseName
		for i := settingTrackLength - h.location; i < settingTrackLength; i++ {
			dis += "."
		}
	} else {
		dis += "[+？]"
		for i := 0; i < settingTrackLength; i++ {
			dis += "."
		}
	}
	dis += "\n"
	return dis
}
