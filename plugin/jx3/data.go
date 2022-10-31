package jx3

import (
	"os"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/process"
	sql "github.com/FloatTech/sqlite"
	log "github.com/sirupsen/logrus"
)

const (
	dbpath      = "data/jx3/"
	dbfile      = dbpath + "robotData.db"
	iconfile    = dbpath + "mental_icon/"
	fileUrl     = "https://raw.githubusercontent.com/Cha0sIDL/data/master/jx/"
	dbMental    = "mental"
	dbControl   = "jxControl"
	dbTeam      = "team"
	dbLeader    = "leader"
	dbMember    = "member"
	dbAdventure = "adventure"
	dbTalk      = "talk"
	dbNews      = "news"
	dbUser      = "user"
	dbIp        = "ip"
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
	dbIp:        &Ip{},
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
