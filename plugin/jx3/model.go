package jx3

import (
	"errors"
	"sort"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Mental的结构体
// type mental struct {
//	ID         uint64 `gorm:"column:mentalID"`
//	Name       string `gorm:"column:mentalName"`
//	Accept     string `gorm:"column:acceptName"`
//	OfficialID int    `gorm:"column:officialID"`
//}

type jxControl struct {
	GroupID int64  `gorm:"column:gid"`     // GroupID 群号
	Disable bool   `gorm:"column:disable"` // Disable 是否启用推送
	Area    string `gorm:"column:area"`    // 绑定的区服
}

// Team 的结构体
type Team struct {
	TeamID   uint   `gorm:"primary_key;AUTO_INCREMENT"`
	LeaderID int64  `gorm:"column:leaderId"` // 团长id
	Dungeon  string `gorm:"column:dungeon"`  // 副本名
	Comment  string `gorm:"column:comment"`  // 备注信息
	GroupID  int64  `gorm:"column:groupId"`  // 团所属群组
}

// Member 开团团员记录结构体
type Member struct {
	ID             uint   `gorm:"primary_key;AUTO_INCREMENT"`
	TeamID         uint   `gorm:"column:team_id"`
	MemberQQ       int64  `gorm:"column:member_qq"`
	MemberNickName string `gorm:"column:member_nick_name"`
	MentalID       uint64 `gorm:"column:mental_id"`
	Double         int    `gorm:"column:double"`
	SignUp         int64  `gorm:"column:sign_up"` // 进团时间
}

// Adventure 一些图片的byte缓存
type Adventure struct {
	Name string `gorm:"column:name"`
	Pic  []byte `gorm:"column:pic"`
	Time int64  `gorm:"column:time"`
}

// Jokes 骚话数据
type Jokes struct {
	ID   int64  `gorm:"column:id"`
	Talk string `gorm:"column:talk"`
}

// News 官方新闻数据
type News struct {
	ID    string `gorm:"column:id"` // href
	Date  string `gorm:"column:date"`
	Title string `gorm:"column:title"`
	Kind  string `gorm:"column:kind"`
}

// User 爬取的内部数据
type User struct {
	ID   string `gorm:"column:id"`
	Data string `gorm:"column:data"` // 服务器的json数据
}

// Daily 日常任务缓存
type Daily struct {
	Server    string `gorm:"column:id"`
	DailyTask string `gorm:"column:dailyTask"`
	Time      int64  `gorm:"column:time"`
}

func (jdb *jx3db) createNewTeam(team *Team) (uint, error) {
	db := (*gorm.DB)(jdb)
	err := db.Create(team).Error
	return team.TeamID, err
}

func (jdb *jx3db) getTeamInfo(teamID int) Team {
	var c Team
	db := (*gorm.DB)(jdb)
	db.Where("team_id = ?", teamID).First(&c)
	return c
}

func (jdb *jx3db) isInTeam(query interface{}, args ...interface{}) bool {
	db := (*gorm.DB)(jdb)
	err := db.Where(query, args...).First(&Member{}).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// 返回我报的团id
func (jdb *jx3db) getSignUp(qq int64) (team []uint) {
	var c []Member
	db := (*gorm.DB)(jdb)
	db.Where("member_qq = ?", qq).Find(&c)
	for _, data := range c {
		team = append(team, data.TeamID)
	}
	return
}

func (jdb *jx3db) delTeam(teamID int, leaderID int64) error {
	var c Team
	db := (*gorm.DB)(jdb)
	db.Where("team_id = ?", teamID).First(&c)
	if c.LeaderID != leaderID {
		return errors.New("这个团队不是你的") // 这个团不是你的
	}
	return db.Where("team_id = ?", teamID).Delete(&Team{}).Error
}

func (jdb *jx3db) addMember(data *Member) error {
	db := (*gorm.DB)(jdb)
	return db.Create(data).Error
}

func (jdb *jx3db) deleteMember(teamID int, qq int64) error {
	db := (*gorm.DB)(jdb)
	return db.Where("team_id = ? and member_qq = ?", teamID, qq).Delete(&Member{}).Error
}

func (jdb *jx3db) getMemberInfo(teamID int) (mSlice []Member) {
	db := (*gorm.DB)(jdb)
	db.Where("team_id = ?", teamID).Find(&mSlice)
	sort.SliceStable(mSlice, func(i, j int) bool {
		return mSlice[i].SignUp < mSlice[j].SignUp
	})
	return
}

// func (jdb *jx3db) getMentalData(mentalName string) mental {
//	db := (*gorm.DB)(jdb)
//	var m mental
//	db.Where("acceptName LIKE ? OR mentalName = ?", fmt.Sprintf("%%%s%%", mentalName), mentalName).First(&m)
//	return m
//}

func (jdb *jx3db) isEnable() (servers map[int64]string) {
	var controls []jxControl
	db := (*gorm.DB)(jdb)
	db.Where("disable = ?", true).Find(&controls)
	servers = lo.Associate(controls, func(item jxControl) (int64, string) {
		return item.GroupID, item.Area
	})
	return
}

func (jdb *jx3db) bind(gid int64) string {
	var control jxControl
	db := (*gorm.DB)(jdb)
	db.Where("gid = ?", gid).First(&control)
	return control.Area
}

func (jdb *jx3db) bindArea(gid int64, area string) {
	var c jxControl
	db := (*gorm.DB)(jdb)
	if err := db.Model(&jxControl{}).First(&c, "gid = ?", gid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.GroupID = gid
			c.Area = area
			c.Disable = true // 默认开
			db.Model(&jxControl{}).Create(&c)
		}
	} else {
		db.Model(&jxControl{}).Where("gid = ?", gid).Update("area", area)
	}
}

func (jdb *jx3db) disable(gid int64) {
	db := (*gorm.DB)(jdb)
	var c jxControl
	if err := db.Model(&jxControl{}).First(&c, "gid = ?", gid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
	} else {
		_ = db.Model(&jxControl{}).Where("gid = ?", gid).Update("disable", false).Error
	}
}

func enable(gid int64) string {
	db := (*gorm.DB)(jdb)
	var c jxControl
	if err := db.Model(&jxControl{}).First(&c, "gid = ?", gid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Area
		}
	} else {
		_ = db.Model(&jxControl{}).Where("gid = ?", gid).Update("disable", true).Error
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
	jdb.insert(data) //nolint:errcheck
}

func (jdb *jx3db) findDaily(server string) (daily Daily) {
	db := (*gorm.DB)(jdb)
	db.Where("id = ?", server).First(&daily)
	return
}

func (jdb *jx3db) pick(out interface{}) (data interface{}) {
	db := (*gorm.DB)(jdb)
	db.Order("random()").Take(&out)
	return out
}

func (jdb *jx3db) insert(value interface{}) error {
	db := (*gorm.DB)(jdb)
	return db.Clauses(clause.OnConflict{UpdateAll: true}).Create(value).Error
}

func (jdb *jx3db) find(query, out interface{}, args ...interface{}) error {
	db := (*gorm.DB)(jdb)
	return db.Where(query, args...).First(out).Error
}

func (jdb *jx3db) findAll(query interface{}, args ...interface{}) error {
	db := (*gorm.DB)(jdb)
	return db.Find(query, args...).Error
}

func (jdb *jx3db) delete(query, value interface{}, args ...interface{}) error {
	db := (*gorm.DB)(jdb)
	return db.Where(query, args).Delete(value).Error
}

func (jdb *jx3db) count(value interface{}) (num int64, err error) {
	db := (*gorm.DB)(jdb)
	err = db.Model(value).Count(&num).Error
	return
}

func (jdb *jx3db) canFind(query, out interface{}, args ...interface{}) bool {
	db := (*gorm.DB)(jdb)
	err := db.Where(query, args...).First(out).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
