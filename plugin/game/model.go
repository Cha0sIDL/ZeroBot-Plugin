package game

import (
	"errors"
	"fmt"
	sql "github.com/FloatTech/sqlite"
)

type gameNotify struct {
	ID       string `db:"id"`
	QQ       int64  `db:"qq"`
	ChatType string `db:"chat_type"`
	GameType string `db:"game_type"` // 预留字段
	RobotId  int64  `db:"robot_id"`
}

const (
	notifyDbName = "gameNotify"
)

var db = &sql.Sqlite{}

func insertNotify(data gameNotify) error {
	isExist := db.CanFind(notifyDbName, fmt.Sprintf("where qq=%d", data.QQ))
	if isExist {
		return errors.New("已经订阅了~")
	}
	err := db.Insert(notifyDbName, &data)
	return err
}
