package jx3

import (
	"fmt"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/process"
	"github.com/glebarez/sqlite"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
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
	dbDaily     = "daily" //每个区服的日常，每天七点刷新
)

// 剑网db全局结构体
var jdb *jx3db

// jx3db 剑网三数据库结构体
type jx3db gorm.DB

var rangeDb = map[string]interface{}{
	dbMental:    &mental{},
	dbControl:   &jxControl{},
	dbTeam:      &Team{},
	dbLeader:    &Leader{},
	dbMember:    &Member{},
	dbAdventure: &Adventure{},
	dbTalk:      &Jokes{},
	dbNews:      &News{},
	dbDaily:     &Daily{},
	dbUser:      &User{},
}

// TableName 表名
func (mental) TableName() string {
	return dbMental
}
func (jxControl) TableName() string {
	return dbControl
}
func (Team) TableName() string {
	return dbTeam
}
func (Leader) TableName() string {
	return dbLeader
}
func (Member) TableName() string {
	return dbMember
}
func (Adventure) TableName() string {
	return dbAdventure
}
func (Jokes) TableName() string {
	return dbTalk
}
func (News) TableName() string {
	return dbNews
}
func (Daily) TableName() string {
	return dbDaily
}

func initialize() *jx3db {
	if file.IsNotExist(dbfile) {
		process.SleepAbout1sTo2s()
		_ = os.MkdirAll(dbpath, 0755)
		err := file.DownloadTo(fileUrl+"robotData.db", dbfile, false)
		if err != nil {
			panic(err)
		}
	}
	jdb, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("jx3 db err,%s", err))
	}
	for _, value := range rangeDb {
		jdb.AutoMigrate(value)
	}
	log.Infoln("[jx3]加载成功")
	return (*jx3db)(jdb)
}
