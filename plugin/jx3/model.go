package jx3

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

// type JxDb gorm.DB

// nsMental的结构体
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
	Area    string `db:"area"`    // 绑定的区服
}

// Team 的结构体
type Team struct {
	TeamId    int    `db:"teamID"`
	LeaderId  int64  `db:"leaderId"`  // 团长id
	Dungeon   string `db:"dungeon"`   // 副本名
	StartTime int64  `db:"startTime"` // 团开始时间
	Comment   string `db:"comment"`   // 备注信息
	GroupId   int64  `db:"groupId"`   // 团所属群组
}

type Leader struct {
	Id       int64  `db:"id"`
	NickName string `db:"nick_name"`
	TeamName string `db:"team_name"`
	IsOk     int    `db:"is_ok"`
}

type Member struct {
	Id             int    `db:"id"`
	TeamId         int    `db:"team_id"`
	MemberQQ       int64  `db:"member_qq"`
	MemberNickName string `db:"member_nick_name"`
	MentalId       uint64 `db:"mental_id"`
	Double         int    `db:"double"`
	SignUp         int64  `db:"sign_up"` // 进团时间
}

func getMental(mentalName string) string {
	var mental mental
	var rwMutex sync.RWMutex
	rwMutex.RLock()
	arg := fmt.Sprintf("WHERE acceptName LIKE '%%%s%%' OR mentalName='%s'", mentalName, mentalName)
	db.Find(dbMental, &mental, arg)
	rwMutex.RUnlock()
	return mental.Name
}

func createNewTeam(time int64, dungeon string,
	comment string, leaderID int64, groupId int64) (int, error) {
	var Mutex sync.Mutex
	Mutex.Lock()
	all, err := db.Count(dbTeam)
	db.Insert(dbTeam, &Team{
		TeamId:    all,
		LeaderId:  leaderID,
		Dungeon:   dungeon,
		StartTime: time,
		Comment:   comment,
		GroupId:   groupId,
	})
	Mutex.Unlock()
	return all, err
}

func getTeamInfo(teamId int) Team {
	var c Team
	db.Find(dbTeam, &c, "WHERE teamID = "+fmt.Sprintln(teamId))
	return c
}

func isInTeam(teamId int, qq int64) bool {
	arg := fmt.Sprintf("WHERE team_id = '%d' ADN member_qq = '%d'", teamId, qq)
	return db.CanFind(dbMember, arg)
}

func isBelongGroup(teamId int, groupId int64) bool {
	arg := fmt.Sprintf("WHERE teamID = '%d' AND groupId = '%d'", teamId, groupId)
	return db.CanFind(dbTeam, arg)
}

// 返回未过期的团
func getEfficientTeamInfo(arg string) []Team {
	var c Team
	var cSlice []Team
	db.FindFor(dbTeam, &c, arg, func() error {
		cSlice = append(cSlice, c)
		return nil
	})
	return cSlice
}

// 返回我报的团id
func getSignUp(qq int64) []int {
	var c Member
	var team []int
	arg := fmt.Sprintf("WHERE member_qq = '%d'", qq)
	db.FindFor(dbMember, &c, arg, func() error {
		team = append(team, c.TeamId)
		return nil
	})
	return team
}

func delTeam(teamId int, leaderId int64) int {
	var c Team
	db.Find(dbTeam, &c, "WHERE teamID = "+fmt.Sprintln(teamId))
	if c.LeaderId != leaderId {
		return -1 // 这个团不是你的
	}
	db.Del(dbTeam, "WHERE teamID = "+fmt.Sprintln(teamId))
	return 0
}

func addMember(data *Member) error {
	var Mutex sync.Mutex
	Mutex.Lock()
	all, err := db.Count(dbMember)
	if err != nil {
		return err
	}
	data.Id = all
	db.Insert(dbMember, data)
	Mutex.Unlock()
	return nil
}

func deleteMember(teamId int, qq int64) error {
	var Mutex sync.Mutex
	Mutex.Lock()
	arg := fmt.Sprintf("WHERE team_id = '%d' AND member_qq = '%d'", teamId, qq)
	db.Del(dbMember, arg)
	Mutex.Unlock()
	return nil
}

func isOk(qq int64) bool {
	var c Leader
	db.Find(dbLeader, &c, "WHERE id = "+fmt.Sprintln(qq))
	return c.IsOk == 1
}

// 添加新团长
func newLeader(QQ int64, nickName string, permission int, teamName ...string) int {
	ok := db.CanFind(dbLeader, "where id="+fmt.Sprintln(QQ))
	if ok {
		return -1 // 数据库中存在记录
	}
	name := ""
	if len(teamName) > 0 {
		name = teamName[0]
	}
	db.Insert(dbLeader, &Leader{
		Id:       QQ,
		NickName: nickName,
		TeamName: name,
		IsOk:     1, // 新团长默认没有权限
	})
	return 0
}

// 同意审批
func acceptLeader(qq int64) string {
	var c Leader
	err := db.Find(dbLeader, &c, "WHERE id = "+fmt.Sprintln(qq))
	if err != nil {
		log.Errorln(err)
		return ""
	}
	c.IsOk = 1
	err = db.Insert(dbLeader, &c)
	return c.TeamName
}

func deleteLeader(qq int64) {
	err := db.Del(dbLeader, "WHERE id = "+fmt.Sprintln(qq))
	if err != nil {
		log.Errorln(err)
	}
}

func getMemberInfo(teamId int) (mSlice []Member) {
	var c Member
	arg := fmt.Sprintf("WHERE team_id = '%d'", teamId)
	db.FindFor(dbMember, &c, arg, func() error {
		mSlice = append(mSlice, c)
		return nil
	})
	sort.SliceStable(mSlice, func(i, j int) bool {
		if mSlice[i].SignUp < mSlice[j].SignUp {
			return true
		}
		return false
	})
	return
}

func getMentalData(mentalName string) mental {
	var m mental
	var rwMutex sync.RWMutex
	rwMutex.RLock()
	arg := fmt.Sprintf("WHERE acceptName LIKE '%%%s%%' OR mentalName='%s'", mentalName, mentalName)
	db.Find(dbMental, &m, arg)
	rwMutex.RUnlock()
	return m
}

func isEnable(Gid int64) (bool, string) {
	var control jxControl
	var rwMutex sync.RWMutex
	rwMutex.RLock()
	arg := "where gid = " + strconv.FormatInt(Gid, 10)
	db.Find(dbControl, &control, arg)
	rwMutex.RUnlock()
	return control.Disable, control.Area
}

func bind(Gid int64) string {
	var control jxControl
	var rwMutex sync.RWMutex
	rwMutex.RLock()
	arg := "where gid = " + strconv.FormatInt(Gid, 10)
	db.Find(dbControl, &control, arg)
	rwMutex.RUnlock()
	return control.Area
}

func bindArea(Gid int64, Area string) {
	var c jxControl
	var rwMutex sync.RWMutex
	rwMutex.RLock()
	err := db.Find(dbControl, &c, "WHERE gid = "+strconv.FormatInt(Gid, 10))
	rwMutex.RUnlock()
	if err != nil {
		c.GroupID = Gid
	}
	c.Area = Area
	rwMutex.Lock()
	err = db.Insert(dbControl, &c)
	rwMutex.Unlock()
	if err != nil {
		log.Error("jx push disable database error")
	}
}

func disable(Gid int64) {
	var c jxControl
	var rwMutex sync.RWMutex
	rwMutex.RLock()
	err := db.Find(dbControl, &c, "WHERE gid = "+strconv.FormatInt(Gid, 10))
	rwMutex.RUnlock()
	if err != nil {
		c.GroupID = Gid
	}
	c.Disable = false
	rwMutex.Lock()
	err = db.Insert(dbControl, &c)
	rwMutex.Unlock()
	if err != nil {
		log.Error("jx push disable database error")
	}
}

func enable(Gid int64) string {
	var c jxControl
	var rwMutex sync.RWMutex
	rwMutex.RLock()
	err := db.Find(dbControl, &c, "WHERE gid = "+strconv.FormatInt(Gid, 10))
	rwMutex.RUnlock()
	if err != nil {
		c.GroupID = Gid
	}
	c.Disable = true
	rwMutex.Lock()
	err = db.Insert(dbControl, &c)
	rwMutex.Unlock()
	if err != nil {
		log.Error("jx push enable database error")
	}
	return c.Area
}
