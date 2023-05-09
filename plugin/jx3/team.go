// Package jx3 团队相关
package jx3

import (
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/golang-module/carbon/v2"
	"github.com/samber/lo"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	"image"
	"strconv"
)

func init() {
	// 开团 副本名 备注
	en.OnPrefixGroup([]string{"开团", "新建团队", "创建团队"}, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) != 2 {
				ctx.SendChain(message.Text("开团参数输入有误！开团 副本名 备注"))
				return
			}
			dungeon := commandPart[0]
			comment := commandPart[1]
			teamID, err := jdb.createNewTeam(&Team{
				LeaderID: ctx.Event.UserID,
				Dungeon:  dungeon,
				Comment:  comment,
				GroupID:  ctx.Event.GroupID,
			})
			if err != nil {
				ctx.SendChain(message.Text("Error :", err))
				return
			}
			ctx.SendChain(message.Text("开团成功，团队id为：", teamID))
		})
	// 报团 团队ID 心法 角色名 [是否双休] 按照报名时间先后默认排序 https://docs.qq.com/doc/DUGJRQXd1bE5YckhB
	en.OnPrefixGroup([]string{"报名", "报团", "报名团队", "代报名"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			double := 0
			switch {
			case len(commandPart) == 3:
				double = 0
			case len(commandPart) == 4:
				double, _ = strconv.Atoi(commandPart[3])
			default:
				ctx.SendChain(message.Text("报团参数有误。"))
				return
			}
			teamID, err := strconv.Atoi(commandPart[0])
			if err != nil {
				ctx.SendChain(message.Text("团队编号输入有误"))
				return
			}
			mental := getMentalData(commandPart[1])
			nickName := commandPart[2]
			if mental.officialID == 0 {
				ctx.SendChain(message.Text("心法输入有误"))
				return
			}
			Team := jdb.getTeamInfo(teamID)
			if Team.GroupID != ctx.Event.GroupID {
				ctx.SendChain(message.Text("当前团队不存在。"))
				return
			}
			if []rune(ctx.MessageString())[0] == '代' {
				if jdb.isInTeam("team_id = ? and member_nick_name = ?", teamID, nickName) {
					ctx.SendChain(message.Text(nickName, "已经在团队中了。"))
					return
				}
			} else {
				if jdb.isInTeam("team_id = ? and member_qq = ?", teamID, ctx.Event.UserID) {
					ctx.SendChain(message.Text("你已经在团队中了。"))
					return
				}
			}
			var member = Member{
				TeamID:         uint(teamID),
				MemberQQ:       ctx.Event.UserID,
				MemberNickName: nickName,
				MentalID:       mental.officialID,
				Double:         double,
				SignUp:         carbon.Now().Timestamp(),
			}
			err = jdb.addMember(&member)
			if err != nil {
				ctx.SendChain(message.Text("数据库写入失败,Err:", err))
				return
			}
			ctx.SendChain(message.Text("报团成功"), message.Reply(ctx.Event.MessageID))
			ctx.SendChain(message.Text("当前团队:\n"), message.Image("base64://"+helper.BytesToString(util.Image2Base64(drawTeam(teamID)))))
		})
	en.OnPrefixGroup([]string{"撤销报团"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			teamID, _ := strconv.Atoi(commandPart[0])
			team := jdb.getTeamInfo(teamID)
			if team.GroupID != ctx.Event.GroupID {
				ctx.SendChain(message.Text("参数输入有误。"))
				return
			}
			err := jdb.deleteMember(teamID, ctx.Event.UserID)
			if err != nil {
				ctx.SendChain(message.Text("Err:", err))
				return
			}
			ctx.SendChain(message.Text("撤销成功"), message.Reply(ctx.Event.MessageID))
			ctx.SendChain(message.Text("当前团队:\n"), message.Image("base64://"+helper.BytesToString(util.Image2Base64(drawTeam(teamID)))))
		})
	en.OnFullMatchGroup([]string{"我报的团", "我的报名"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var sTeam []Team
			err := jdb.Find("groupId = ?", &sTeam, ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Text("Err:", err))
			}
			s := lo.Map(sTeam, func(item Team, _ int) uint {
				return item.TeamID
			})
			SignUp := lo.Uniq(jdb.getSignUp(ctx.Event.UserID))

			ctx.SendChain(message.Text("本群你报名过的团队id：\n", util.IntersectArray(s, SignUp)))
		})
	en.OnFullMatchGroup([]string{"我的开团"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var sTeam []Team
			err := jdb.Find("leaderId = ? and groupId = ?", &sTeam, ctx.Event.UserID, ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Text("Err:", err))
				return
			}
			out := ""
			for _, data := range sTeam {
				out += fmt.Sprintf("团队id：%d,团长 ：%d,副本：%s，备注：%s\n",
					data.TeamID, data.LeaderID, data.Dungeon, data.Comment)
			}
			ctx.SendChain(message.Text(out))
		})
	en.OnFullMatchGroup([]string{"查看全部团队"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var sTeam []Team
			err := jdb.Find("groupId = ?", &sTeam, ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Text("Err:", err))
				return
			}
			if len(sTeam) == 0 {
				ctx.SendChain(message.Text("本群没有有效团队哦"))
				return
			}
			out := ""
			for _, data := range sTeam {
				out += fmt.Sprintf("团队id：%d,团长 ：%d,副本：%s，备注：%s\n",
					data.TeamID, data.LeaderID, data.Dungeon, data.Comment)
			}
			ctx.SendChain(message.Text(out))
		})
	// 查看团队 teamid
	en.OnPrefixGroup([]string{"查看团队", "查询团队", "查团"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			teamID, _ := strconv.Atoi(commandPart[0])
			team := jdb.getTeamInfo(teamID)
			if team.GroupID != ctx.Event.GroupID {
				ctx.SendChain(message.Text("团队id输入有误。"))
				return
			}
			ctx.SendChain(message.Image("base64://" + helper.BytesToString(util.Image2Base64(drawTeam(teamID)))))
		})
	// 取消开团 团队id
	en.OnPrefixGroup([]string{"取消开团", "删除团队", "结束团队"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			if len(commandPart) < 1 {
				ctx.SendChain(message.Text("参数有误"))
			}
			teamID, err := strconv.Atoi(commandPart[0])
			team := jdb.getTeamInfo(teamID)
			if err != nil || team.GroupID != ctx.Event.GroupID || team.LeaderID != ctx.Event.UserID {
				ctx.SendChain(message.Text("团队id输入有误"))
				return
			}
			err = jdb.delTeam(teamID, ctx.Event.UserID)
			if err != nil {
				ctx.SendChain(message.Text(err))
			}
			ctx.SendChain(message.Text("取消成功"))
		})
}

func drawTeam(teamID int) image.Image {
	Fonts, err := gg.LoadFontFace(text.FontFile, 50)
	if err != nil {
		panic(err)
	}
	const W = 1200
	const H = 1200
	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	// 画直线
	for i := 0; i < 1200; {
		dc.SetRGBA(255, 255, 255, 11)
		dc.SetLineWidth(1)
		dc.DrawLine(0, float64(i), 1200, float64(i))
		dc.Stroke()
		i += 200
	}
	// 画直线
	for i := 200; i < 1200; {
		// dc.SetRGBA(255, 255, 255, 11)
		// dc.SetLineWidth(1)
		dc.DrawLine(float64(i), 200, float64(i), 1200)
		dc.Stroke()
		i += 200
	}
	dc.SetFontFace(Fonts)
	// 队伍
	for i := 1; i < 6; i++ {
		dc.DrawString(strconv.Itoa(i)+"队", 40, float64(100+200*i))
	}
	// 标题
	team := jdb.getTeamInfo(teamID)
	title := strconv.Itoa(int(team.TeamID)) + " " + team.Dungeon
	_, th := dc.MeasureString("哈")
	t := 1200/2 - (float64(len([]rune(title))) / 2)
	dc.DrawStringAnchored(title, t, th, 0.5, 0.5)
	dc.DrawStringAnchored(team.Comment, 1200/2-float64(len([]rune(team.Comment)))/2, 3*th, 0.5, 0.5)
	// 团队
	mSlice := jdb.getMemberInfo(teamID)
	dc.LoadFontFace(text.FontFile, 30) //nolint:errcheck
	_, th = dc.MeasureString("哈")
	start := 200
	for idx, m := range mSlice {
		x := float64(start + idx%5*200 + 10)
		y := float64(start+idx/5*200) + th*2
		dc.DrawString(m.MemberNickName, x, y)
		double := "单修"
		if m.Double == 1 {
			double = "双修"
		}
		dc.DrawString(double, x, y+th*2)
		back, _ := gg.LoadImage(util.IconFilePath + strconv.Itoa(int(m.MentalID)) + ".png")
		dc.DrawImage(back, int(x), int(y+th*3))
	}
	return dc.Image()
}
