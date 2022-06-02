package chat_record

import (
	sql "github.com/FloatTech/sqlite"
)

var db = &sql.Sqlite{}

// 聊天记录存储结构体
type record struct {
	MId     interface{}
	Message string
	GroupId int64
	Time    int64
	UserID  int64
}
