package HorseRace

import (
	"encoding/json"
	"os"
)

const (
	// 跑道长度
	settingTrackLength = 15
	// 随机位置事件，最小能到的跑道距离
	settingRandomMinLength = 0
	// 随机位置事件，最大能到的跑道距离
	settingRandomMaxLength = 10
	// 每回合基础移动力最小值
	baseMoveMin = 1
	// 每回合基础移动力最大值
	baseMoveMax = 3
	// 最大支持玩家数
	maxPlayer = 8
	// 最少玩家数
	minPlayer = 2
	// 赛马超时时间，秒
	settingOverTime = 120
	// 事件概率 = event_rate / 1000
	eventRate = 300
	// 马儿名字最大长度
	nameMaxLen = 8
)

var events []event

func initConfig(path string) {
	events = []event{}
	files, _ := os.ReadDir(path)
	for _, f := range files {
		var e []event
		content, err := os.ReadFile(path + f.Name())
		if err != nil {
			panic(err)
		}
		json.Unmarshal(content, &e)
		events = append(events, e...)
	}
}
