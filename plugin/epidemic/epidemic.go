// Package epidemic 城市疫情查询
package epidemic

import (
	"encoding/json"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/zbputils/file"
	"github.com/playwright-community/playwright-go"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/control/order"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
)

const (
	servicename = "epidemic"
	txurl       = "https://view.inews.qq.com/g2/getOnsInfo?name=disease_h5"
)

// result 疫情查询结果
type result struct {
	Data string `json:"data"`
}

// epidemic 疫情数据
type epidemic struct {
	LastUpdateTime string  `json:"lastUpdateTime"`
	AreaTree       []*area `json:"areaTree"`
}

// area 城市疫情数据
type area struct {
	Name  string `json:"name"`
	Today struct {
		Confirm int `json:"confirm"`
		Wzzadd  int `json:"wzz_add"`
	} `json:"today"`
	Total struct {
		NowConfirm int    `json:"nowConfirm"`
		Confirm    int    `json:"confirm"`
		Dead       int    `json:"dead"`
		Heal       int    `json:"heal"`
		Grade      string `json:"grade"`
		Wzz        int    `json:"wzz"`
	} `json:"total"`
	Children []*area `json:"children"`
}

func init() {
	engine := control.Register(servicename, order.AcquirePrio(), &control.Options{
		DisableOnDefault: false,
		Help: "城市疫情查询\n" +
			"- xxx疫情\n",
		PublicDataFolder: "Epidemic",
	})
	datapath := file.BOTPATH + "/" + engine.DataFolder()
	engine.OnSuffix("疫情").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			city := ctx.State["args"].(string)
			if city == "" {
				ctx.SendChain(message.Text("你还没有输入城市名字呢！"))
				return
			}
			data, time, err := queryEpidemic(city)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			if data == nil {
				ctx.SendChain(message.Text("没有找到【", city, "】城市的疫情数据."))
				return
			}
			d := map[string]interface{}{
				"NowConfirm": data.Total.NowConfirm, //现有确诊
				"Confirm":    data.Today.Confirm,    //新增人数
				"Heal":       data.Total.Heal,       //累计确诊
				"deadCount":  data.Total.Dead,       //死亡人数
				"Grade":      data.Total.Grade,      //
				"name":       data.Name,
				"time":       time,
				"Wzz":        data.Total.Wzz,    //无症状人数
				"Wzzadd":     data.Today.Wzzadd} //新增无症状
			html := util.Template2html("yiqing.html", d)
			Clip := util.PageScreenshotOptionsClip(
				playwright.PageScreenshotOptionsClip{
					X:      playwright.Float(10),
					Y:      playwright.Float(0),
					Width:  playwright.Float(1100),
					Height: playwright.Float(450),
				})
			finName, err := util.Html2pic(datapath, util.TodayFileName(), "yiqing.html", html, Clip)
			ctx.SendChain(
				//message.Text(
				//	"【", data.Name, "】疫情数据\n",
				//	"新增人数：", data.Today.Confirm, "\n",
				//	"现有确诊：", data.Total.NowConfirm, "\n",
				//	"累计确诊：", data.Total.Confirm, "\n",
				//	"治愈人数：", data.Total.Heal, "\n",
				//	"死亡人数：", data.Total.Dead, "\n",
				//	"无症状人数：", data.Total.Wzz, "\n",
				//	"新增无症状：", data.Today.Wzzadd, "\n",
				//	"更新时间：\n『", time, "』",
				//),
				message.Image("file:///" + finName),
			)
		})
}

// rcity 查找城市
func rcity(a *area, cityName string) *area {
	if a == nil {
		return nil
	}
	if a.Name == cityName {
		return a
	}
	for _, v := range a.Children {
		if v.Name == cityName {
			return v
		}
		c := rcity(v, cityName)
		if c != nil {
			return c
		}
	}
	return nil
}

// queryEpidemic 查询城市疫情
func queryEpidemic(findCityName string) (citydata *area, times string, err error) {
	data, err := web.GetData(txurl)
	if err != nil {
		return
	}
	var r result
	err = json.Unmarshal(data, &r)
	if err != nil {
		return
	}
	var e epidemic
	err = json.Unmarshal(binary.StringToBytes(r.Data), &e)
	if err != nil {
		return
	}
	citydata = rcity(e.AreaTree[0], findCityName)
	return citydata, e.LastUpdateTime, nil
}
