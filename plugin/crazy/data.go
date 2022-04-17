package crazy

import sql "github.com/FloatTech/sqlite"

var db = &sql.Sqlite{}

//nsMental的结构体
type crazy struct {
	Crazy string `db:"crazy"`
}

//
type menu struct {
	Menu string `db:"menu"`
}
