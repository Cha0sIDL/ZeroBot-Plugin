package jx3

import (
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
)

func init() {
	// 万宝楼 外观/角色 编号
	en.OnPrefix("万宝楼", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(
		func(ctx *zero.Ctx) {
			var productType int
			commandPart := util.SplitSpace(ctx.State["args"].(string))
			product := commandPart[0]
			productNum := commandPart[1]
			switch {
			case len(commandPart) != 2:
				ctx.SendChain(message.Text("参数输入有误！\n" + "万宝楼 外观/角色 万宝楼编号"))
				return
			case checkNum(productNum):
				ctx.SendChain(message.Text("检测到商品编号出现了非数字的内容！\n请检查后重试~"))
				return
			case product == "外观":
				productType = 3
			case product == "角色":
				productType = 2
			default:
				ctx.SendChain(message.Text("发生了未知错误联系管理员看看吧"))
				return
			}
			url := fmt.Sprintf("https://api-wanbaolou.xoyo.com/api/buyer/goods/list?goods_type=%d&consignment_id=%s", productType, productNum)
			data, err := web.GetData(url)
			if err != nil {
				ctx.SendChain(message.Text("Err:", err))
				return
			}
			jsonData := gjson.ParseBytes(data)
			if jsonData.Get("code").Int() != 1 {
				ctx.SendChain(message.Text("请求出错了联系管理员看看吧～"))
				return
			}
			if len(jsonData.Get("data.list").Array()) <= 0 {
				ctx.SendChain(message.Text("没有查找到此角色/外观哦，检查下输入编号试试吧～"))
				return
			}
			wanbaolouData := jsonData.Get("data.list").Array()[0]
			zone := wanbaolouData.Get("zone_name").String()
			server := wanbaolouData.Get("server_name").String()
			seller := wanbaolouData.Get("seller_role_name").String()
			productData := wanbaolouData.Get("info").String()
			remain := time2person(wanbaolouData.Get("remaining_time").Int())
			price := wanbaolouData.Get("single_unit_price").Int() / 100
			followerCount := wanbaolouData.Get("followed_num").Int()
			if product == "外观" {
				productDetailType := wanbaolouData.Get("attrs.appearance_type_name").String()
				ctx.SendChain(message.Text("商品名称：", productData, "\n", "售卖者：", zone, server, seller, "\n", "点赞：", followerCount, "\n", "剩余时间：", remain, "\n", "商品类型：", productDetailType, "\n", "价格：", price))
			} else {
				tags := wanbaolouData.Get("tags").String()
				level := wanbaolouData.Get("attrs.role_level").Int()
				equip := wanbaolouData.Get("attrs.role_equipment_point").Int()
				experience := wanbaolouData.Get("attrs.role_experience_point").Int()
				sect := wanbaolouData.Get("attrs.role_sect").String()
				camp := wanbaolouData.Get("attrs.role_camp").String()
				shape := wanbaolouData.Get("attrs.role_shape").String()
				ctx.SendChain(message.Image(wanbaolouData.Get("thumb").String()), message.Text("商品名称：", productData, "\n售卖者：", zone, server, seller, "\n点赞：", followerCount, "\n剩余时间：", remain, "\n标签：", tags, "\n价格：", price, "\n体型：", sect, "-", camp, "-", shape, "\n装备分数：", level, "级", equip, "\n资历：", experience))
			}
		})
}

func checkNum(s string) bool {
	_, err := strconv.Atoi(s)
	return err != nil
}

func time2person(time int64) string {
	return fmt.Sprintf("%d时%d分", time/3600, time%3600/60)
}
