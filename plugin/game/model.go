package game

import (
	"errors"
	"fmt"
	"sync"

	sql "github.com/FloatTech/sqlite"
)

type gameNotify struct {
	ID       string `db:"id"`
	QQ       int64  `db:"qq"`
	ChatType string `db:"chat_type"`
	GameType string `db:"game_type"` // 预留字段
	RobotID  int64  `db:"robot_id"`
}

const (
	notifyDBName = "gameNotify"
)

var db = &sql.Sqlite{}

func insertNotify(data gameNotify) error {
	var mutex sync.RWMutex
	mutex.RLock()
	defer mutex.RUnlock()
	isExist := db.CanFind(notifyDBName, fmt.Sprintf("where qq=%d", data.QQ))
	if isExist {
		return errors.New("已经订阅了~")
	}
	err := db.Insert(notifyDBName, &data)
	return err
}
