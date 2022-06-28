package game

import sql "github.com/FloatTech/sqlite"

type gameNotify struct {
	ID       int    `db:"id"`
	QQ       int64  `db:"qq"`
	ChatType string `db:"chat_type"`
	GameType string `db:"game_type"` //预留字段
	RobotId  int64  `db:"robot_id"`
}

const (
	notifyDbName = "gameNotify"
)

var db = &sql.Sqlite{}

func insertNotify(data gameNotify) error {
	primaryKey, _ := db.Count(notifyDbName)
	data.ID = primaryKey + 1
	err := db.Insert(notifyDbName, &data)
	return err
}
