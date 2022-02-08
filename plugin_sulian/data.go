package sulian

import (
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/process"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	dbpath  = "data/sulian/"
	dbfile  = dbpath + "sulian.json"
	fileUrl = "https://raw.githubusercontent.com/Cha0sIDL/data/master/sulian/sulian.json"
)

// 加载数据库
func init() {
	go down()
	logrus.Infoln("[sulian]加载成功")
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
}
