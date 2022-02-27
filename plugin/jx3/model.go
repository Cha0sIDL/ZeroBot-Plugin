package jx3

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

//nsMental的结构体
type mental struct {
	ID          uint64 `db:"mentalID"`
	Name        string `db:"mentalName"`
	MentalIcon  string `db:"mentalIcon"`
	Accept      string `db:"acceptName"`
	MentalColor string `db:"mentalColor"`
	Works       int    `db:"works"`
	Relation    int    `db:"relation"`
}

//
type jxControl struct {
	GroupID int64  `db:"gid"`     // GroupID 群号
	Disable bool   `db:"disable"` // Disable 是否启用推送
	Area    string `db:"area"`    //绑定的区服
}

func getMental(mentalName string) string {
	var mental mental
	var rwMutex sync.RWMutex
	rwMutex.RLock()
	arg := fmt.Sprintf("WHERE acceptName LIKE '%%%s%%' OR mentalName='%s'", mentalName, mentalName)
	db.Find("ns_mental", &mental, arg)
	rwMutex.RUnlock()
	return mental.Name
}

func getData(mentalName string) mental {
	var m mental
	var rwMutex sync.RWMutex
	rwMutex.RLock()
	arg := fmt.Sprintf("WHERE acceptName LIKE '%%%s%%' OR mentalName='%s'", mentalName, mentalName)
	db.Find("ns_mental", &m, arg)
	rwMutex.RUnlock()
	return m
}

func isEnable(Gid int64) (bool, string) {
	var control jxControl
	var rwMutex sync.RWMutex
	rwMutex.RLock()
	arg := "where gid = " + strconv.FormatInt(Gid, 10)
	db.Find("jxControl", &control, arg)
	rwMutex.RUnlock()
	return control.Disable, control.Area
}

func bindArea(Gid int64, Area string) {
	var c jxControl
	var rwMutex sync.RWMutex
	rwMutex.RLock()
	err := db.Find("jxControl", &c, "WHERE gid = "+strconv.FormatInt(Gid, 10))
	rwMutex.RUnlock()
	if err != nil {
		c.GroupID = Gid
	}
	c.Area = Area
	rwMutex.Lock()
	err = db.Insert("jxControl", &c)
	rwMutex.Unlock()
	if err != nil {
		log.Error("jx push disable database error")
	}
}

func disable(Gid int64) {
	var c jxControl
	var rwMutex sync.RWMutex
	rwMutex.RLock()
	err := db.Find("jxControl", &c, "WHERE gid = "+strconv.FormatInt(Gid, 10))
	rwMutex.RUnlock()
	if err != nil {
		c.GroupID = Gid
	}
	c.Disable = false
	rwMutex.Lock()
	err = db.Insert("jxControl", &c)
	rwMutex.Unlock()
	if err != nil {
		log.Error("jx push disable database error")
	}
}

func enable(Gid int64) string {
	var c jxControl
	var rwMutex sync.RWMutex
	rwMutex.RLock()
	err := db.Find("jxControl", &c, "WHERE gid = "+strconv.FormatInt(Gid, 10))
	rwMutex.RUnlock()
	if err != nil {
		c.GroupID = Gid
	}
	c.Disable = true
	rwMutex.Lock()
	err = db.Insert("jxControl", &c)
	rwMutex.Unlock()
	if err != nil {
		log.Error("jx push enable database error")
	}
	return c.Area
}
