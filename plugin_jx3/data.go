package jx3

import (
	sql "github.com/FloatTech/sqlite"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/process"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	dbpath  = "data/jx3/"
	dbfile  = dbpath + "robotData.db"
	fileUrl = "https://raw.githubusercontent.com/Cha0sIDL/data/master/jx/robotData.db"
)

var db = &sql.Sqlite{DBPath: dbfile}

// 加载数据库
func init() {
	go down()
}

func down() {
	if file.IsNotExist(dbfile) {
		process.SleepAbout1sTo2s()
		_ = os.MkdirAll(dbpath, 0755)
		err := file.DownloadTo(fileUrl, dbfile, false)
		if err != nil {
			panic(err)
		}
	}
	db.Create("ns_mental", &mental{})
	db.Create("jxControl", &jxControl{})
	logrus.Infoln("[jx3]加载成功")
}
