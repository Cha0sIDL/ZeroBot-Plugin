package jx3

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sort"
)

// Mental的结构体
type mental struct {
	ID     uint64 `gorm:"column:mentalID"`
	Name   string `gorm:"column:mentalName"`
	Accept string `gorm:"column:acceptName"`
}

type jxControl struct {
	GroupID int64  `gorm:"column:gid"`     // GroupID 群号
	Disable bool   `gorm:"column:disable"` // Disable 是否启用推送
	Area    string `gorm:"column:area"`    // 绑定的区服
}

// Team 的结构体
type Team struct {
	TeamId    uint   `gorm:"primary_key;AUTO_INCREMENT"`
	LeaderId  int64  `gorm:"column:leaderId"`  // 团长id
	Dungeon   string `gorm:"column:dungeon"`   // 副本名
	StartTime int64  `gorm:"column:startTime"` // 团开始时间
	Comment   string `gorm:"column:comment"`   // 备注信息
	GroupId   int64  `gorm:"column:groupId"`   // 团所属群组
}

type Leader struct {
	Id       uint   `gorm:"primary_key;AUTO_INCREMENT"`
	NickName string `gorm:"column:nick_name"`
	TeamName string `gorm:"column:team_name"`
	IsOk     int    `gorm:"column:is_ok"`
}

type Member struct {
	Id             uint   `gorm:"primary_key;AUTO_INCREMENT"`
	TeamId         int    `gorm:"column:team_id"`
	MemberQQ       int64  `gorm:"column:member_qq"`
	MemberNickName string `gorm:"column:member_nick_name"`
	MentalId       uint64 `gorm:"column:mental_id"`
	Double         int    `gorm:"column:double"`
	SignUp         int64  `gorm:"column:sign_up"` // 进团时间
}

type Adventure struct {
	Name string `gorm:"column:name"`
	Pic  []byte `gorm:"column:pic"`
	Time int64  `gorm:"column:time"`
}

type Jokes struct {
	ID   int64  `gorm:"column:id"`
	Talk string `gorm:"column:talk"`
}

type News struct {
	ID    string `gorm:"column:id"` // href
	Date  string `gorm:"column:date"`
	Title string `gorm:"column:title"`
	Kind  string `gorm:"column:kind"`
}

type User struct {
	ID   string `gorm:"column:id"`
	Data string `gorm:"column:data"` // 服务器的json数据
}

type Daily struct {
	Server    string `gorm:"column:id"`
	DailyTask string `gorm:"column:dailyTask"`
	Time      int64  `gorm:"column:time"`
}

func (jdb *jx3db) createNewTeam(time int64, dungeon string,
	comment string, leaderID int64, groupId int64) (uint, error) {
	db := (*gorm.DB)(jdb)
	team := &Team{
		LeaderId:  leaderID,
		Dungeon:   dungeon,
		StartTime: time,
		Comment:   comment,
		GroupId:   groupId,
	}
	err := db.Create(team).Error
	return team.TeamId, err
}

func (jdb *jx3db) getTeamInfo(teamId int) Team {
	var c Team
	db := (*gorm.DB)(jdb)
	db.Where("team_id = ?", teamId).First(&c)
	return c
}

func (jdb *jx3db) isInTeam(teamId int, qq int64) bool {
	db := (*gorm.DB)(jdb)
	var c Member
	err := db.Where("team_id = ? and member_qq = ?", teamId, qq).First(&c).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (jdb *jx3db) isBelongGroup(teamId int, groupId int64) bool {
	db := (*gorm.DB)(jdb)
	var c Team
	err := db.Where("team_id = ? and groupId = ?", teamId, groupId).First(&c).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// 返回未过期的团
func (jdb *jx3db) getEfficientTeamInfo(query interface{}, args ...interface{}) (cSlice []Team) {
	db := (*gorm.DB)(jdb)
	db.Where(query, args...).Find(&cSlice)
	return
}

// 返回我报的团id
func (jdb *jx3db) getSignUp(qq int64) (team []int) {
	var c []Member
	db := (*gorm.DB)(jdb)
	db.Where("member_qq = ?", qq).Find(&c)
	for _, data := range c {
		team = append(team, data.TeamId)
	}
	return
}

func (jdb *jx3db) delTeam(teamId int, leaderId int64) error {
	var c Team
	db := (*gorm.DB)(jdb)
	db.Where("team_id = ?", teamId).First(&c)
	if c.LeaderId != leaderId {
		return errors.New("这个团队不是你的") // 这个团不是你的
	}
	return db.Where("team_id = ?", teamId).Delete(&Team{}).Error
}

func (jdb *jx3db) addMember(data *Member) error {
	db := (*gorm.DB)(jdb)
	db.Create(data)
	return nil
}

func (jdb *jx3db) deleteMember(teamId int, qq int64) error {
	db := (*gorm.DB)(jdb)
	return db.Where("team_id = ? and member_qq = ?", teamId, qq).Delete(&Member{}).Error
}

func (jdb *jx3db) getMemberInfo(teamId int) (mSlice []Member) {
	db := (*gorm.DB)(jdb)
	db.Where("team_id = ?", teamId).Find(&mSlice)
	sort.SliceStable(mSlice, func(i, j int) bool {
		if mSlice[i].SignUp < mSlice[j].SignUp {
			return true
		}
		return false
	})
	return
}

func (jdb *jx3db) getMentalData(mentalName string) mental {
	db := (*gorm.DB)(jdb)
	var m mental
	db.Where("acceptName LIKE ? OR mentalName = ?", fmt.Sprintf("%%%s%%", mentalName), mentalName).First(&m)
	return m
}

func (jdb *jx3db) isEnable(Gid int64) (bool, string) {
	var control jxControl
	db := (*gorm.DB)(jdb)
	db.Where("gid = ?", Gid).First(&control)
	return control.Disable, control.Area
}

func (jdb *jx3db) bind(Gid int64) string {
	var control jxControl
	db := (*gorm.DB)(jdb)
	db.Where("gid = ?", Gid).First(&control)
	return control.Area
}

func (jdb *jx3db) bindArea(Gid int64, Area string) {
	var c jxControl
	db := (*gorm.DB)(jdb)
	if err := db.Model(&jxControl{}).First(&c, "gid = ?", Gid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.GroupID = Gid
			c.Area = Area
			c.Disable = true // 默认开
			err = db.Model(&jxControl{}).Create(&c).Error
		}
	} else {
		c.Area = Area
		err = db.Model(&jxControl{}).Where("gid = ?", Gid).Updates(&c).Error
	}
}

func (jdb *jx3db) disable(Gid int64) {
	db := (*gorm.DB)(jdb)
	var c jxControl
	if err := db.Model(&jxControl{}).First(&c, "gid = ?", Gid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
	} else {
		err = db.Model(&jxControl{}).Where("gid = ?", Gid).Update("disable", false).Error
	}
}

func enable(Gid int64) string {
	db := (*gorm.DB)(jdb)
	var c jxControl
	if err := db.Model(&jxControl{}).First(&c, "gid = ?", Gid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Area
		}
	} else {
		err = db.Model(&jxControl{}).Where("gid = ?", Gid).Update("disable", true).Error
		return c.Area
	}
	return ""
}

func (jdb *jx3db) getAdventure(name string) Adventure {
	var data Adventure
	db := (*gorm.DB)(jdb)
	db.Where("name = ?", name).First(&data)
	return data
}

func (jdb *jx3db) updateAdventure(data *Adventure) {
	jdb.Insert(data)
}

func (jdb *jx3db) findDaily(server string) (daily Daily) {
	db := (*gorm.DB)(jdb)
	db.Where("id = ?", server).First(&daily)
	return
}

func (jdb *jx3db) Pick(out interface{}) (data interface{}) {
	db := (*gorm.DB)(jdb)
	db.Order("random()").Take(&out)
	return out
}

func (jdb *jx3db) Insert(value interface{}) error {
	db := (*gorm.DB)(jdb)
	return db.Clauses(clause.OnConflict{UpdateAll: true}).Create(value).Error
}

func (jdb *jx3db) Find(query, out interface{}, args ...interface{}) error {
	db := (*gorm.DB)(jdb)
	return db.Where(query, args).Find(out).Error
}

func (jdb *jx3db) Count(value interface{}) (num int64, err error) {
	db := (*gorm.DB)(jdb)
	err = db.Model(value).Count(&num).Error
	return
}

func (jdb *jx3db) CanFind(query, out interface{}, args ...interface{}) bool {
	db := (*gorm.DB)(jdb)
	err := db.Where(query, args).First(out).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
