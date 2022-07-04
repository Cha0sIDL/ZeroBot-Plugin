package jx3

import (
	"os"
	"time"

	sql "github.com/FloatTech/sqlite"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/process"
	log "github.com/sirupsen/logrus"
)

const (
	dbpath      = "data/jx3/"
	dbfile      = dbpath + "robotData.db"
	iconfile    = dbpath + "mental_icon/"
	fileUrl     = "https://raw.githubusercontent.com/Cha0sIDL/data/master/jx/"
	dbMental    = "ns_mental"
	dbControl   = "jxControl"
	dbTeam      = "team"
	dbLeader    = "leader"
	dbMember    = "member"
	dbAdventure = "adventure"
	dbTalk      = "talk"
	dbNews      = "news"
)

var rangeDb = map[string]interface{}{
	dbMental:    &mental{},
	dbControl:   &jxControl{},
	dbTeam:      &Team{},
	dbLeader:    &Leader{},
	dbMember:    &Member{},
	dbAdventure: &Adventure{},
	dbTalk:      &Jokes{},
	dbNews:      &News{},
}

var db = &sql.Sqlite{DBPath: dbfile}

func initialize() {
	if file.IsNotExist(dbfile) {
		process.SleepAbout1sTo2s()
		_ = os.MkdirAll(dbpath, 0755)
		err := file.DownloadTo(fileUrl+"robotData.db", dbfile, false)
		if err != nil {
			panic(err)
		}
	}
	db.Open(time.Hour * 24)
	for key, value := range rangeDb {
		err := db.Create(key, value)
		if err != nil {
			panic(err)
		}
	}
	log.Infoln("[jx3]加载成功")
}

//// 加载数据库
// func init() {
//	go down()
//}

// func down() *JxDb {
//	if file.IsNotExist(dbfile) {
//		process.SleepAbout1sTo2s()
//		_ = os.MkdirAll(dbpath, 0755)
//		err := file.DownloadTo(fileUrl, dbfile, false)
//		if err != nil {
//			panic(err)
//		}
//	}
//
//	for key, value := range rangeDb {
//		db.Create(key, value)
//	}
//	gdb, err := gorm.Open("sqlite3", dbfile)
//	if err != nil {
//		panic(err)
//	}
//	gdb.AutoMigrate(&Member{})
//	logrus.Infoln("[jx3]加载成功")
//	return (*JxDb)(gdb)
//}
