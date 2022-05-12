package trending

import (
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
)

func init() { // 插件主体
	engine := control.Register("trending", &control.Options{
		DisableOnDefault: false,
		Help: "微博热搜\n" +
			"知乎热搜\n" +
			"github热搜",
	})
	engine.OnSuffix("热搜").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			kind := ctx.State["args"].(string)
			switch kind {
			case "微博":
				getWeiboTrending(ctx)
			case "知乎":
				getZhihuTrending(ctx)
			case "github":
				getGithubTrending(ctx)
			case "头条":
				getTouTiaoTrending(ctx)
			}
			return
		})
}
func getWeiboTrending(ctx *zero.Ctx) {
	rsp := "微博实时热榜:\n"
	url := "http://api.weibo.cn/2/guest/search/hot/word"
	data, err := web.RequestDataWith(web.NewDefaultClient(), url, "GET", "", web.RandUA())
	if err != nil {
		msg := message.Text("ERROR:", err)
		ctx.SendChain(msg)
		return
	}
	json := gjson.Get(binary.BytesToString(data), "data").Array()
	for idx, hot := range json {
		if hot.Get("word").String() == "" {
			continue
		}
		rsp = rsp + strconv.Itoa(idx+1) + ":" + hot.Get("word").String() + "\n"
	}
	ctx.SendChain(message.Text(rsp))
	return
}

func getZhihuTrending(ctx *zero.Ctx) {
	rsp := "知乎实时热榜:\n"
	url := "https://www.zhihu.com/api/v3/feed/topstory/hot-lists/total?limit=30&desktop=true"
	data, err := web.RequestDataWith(web.NewDefaultClient(), url, "GET", "", web.RandUA())
	if err != nil {
		msg := message.Text("ERROR:", err)
		ctx.SendChain(msg)
		return
	}
	json := gjson.Get(binary.BytesToString(data), "data").Array()
	for idx, hot := range json {
		if hot.Get("target.title").String() == "" {
			continue
		}
		rsp = rsp + strconv.Itoa(idx+1) + ":" + string(util.Unicode2Zh(hot.Get("target.title").String())) + "\n"
	}
	ctx.SendChain(message.Text(rsp))
	return
}

func getGithubTrending(ctx *zero.Ctx) {
	////rsp := "GitHub实时热榜:\n"
	//url := "https://github.com/trending"
	//data, err := web.RequestDataWith(web.NewDefaultClient(), url, "GET", "", web.RandUA())
	//if err != nil {
	//	msg := message.Text("ERROR:", err)
	//	ctx.SendChain(msg)
	//	return
	//}
	//log.Errorln(string(data))
	//doc, err := htmlquery.Parse(strings.NewReader(helper.BytesToString(data)))
	//if err != nil {
	//	return
	//}
	//htmlquery.Find(doc, "Box-row")
}

func getTouTiaoTrending(ctx *zero.Ctx) {
	rsp := "头条今日热搜:\n"
	url := "https://is-lq.snssdk.com/api/suggest_words/?business_id=10016"
	data, err := web.RequestDataWith(web.NewDefaultClient(), url, "GET", "", web.RandUA())
	if err != nil {
		msg := message.Text("ERROR:", err)
		ctx.SendChain(msg)
		return
	}
	if gjson.Get(binary.BytesToString(data), "msg").String() != "success" {
		return
	}
	json := gjson.Get(binary.BytesToString(data), "data").Array()
	for idx, hot := range json[0].Get("words").Array() {
		if hot.Get("word").String() == "" {
			continue
		}
		rsp = rsp + strconv.Itoa(idx+1) + "：" + hot.Get("word").String() + "\n"
	}
	ctx.SendChain(message.Text(rsp))
	return
}
