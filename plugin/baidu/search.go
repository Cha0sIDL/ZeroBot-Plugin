// Package baidu 百度百科
package baidu

import (
	"encoding/json"
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"net/url"

	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	duURL   = "https://api.a20safe.com/api.php?api=21&key=%s&text=%s" // api地址
	wikiURL = "https://api.a20safe.com/api.php?api=23&key=%s&text=%s"
	key     = "7d06a110e9e20a684e02934549db1d3d"
)

type result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Content string `json:"content"`
	} `json:"data"`
}

func init() { // 主函数
	en := control.Register("baidu", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "baidu\n" +
			"- 百度下[xxx]\n" +
			"- 维基|百科[xxx]\n" +
			"- 百度百科[xxx]",
	})
	en.OnRegex(`^(百度|维基|百科|wiki)\s*(.+)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		var es []byte
		var err error
		switch ctx.State["regex_matched"].([]string)[1] {
		case "百度", "百科":
			es, err = web.GetData(fmt.Sprintf(duURL, key, ctx.State["regex_matched"].([]string)[2])) // 将网站返回结果赋值
		case "wiki", "维基":
			es, err = web.GetData(fmt.Sprintf(wikiURL, key, ctx.State["regex_matched"].([]string)[2])) // 将网站返回结果赋值
		}
		if err != nil {
			ctx.SendChain(message.Text("出现错误捏：", err))
			return
		}
		var r result                 // r数组
		err = json.Unmarshal(es, &r) // 填api返回结果，struct地址
		if err != nil {
			ctx.SendChain(message.Text("出现错误捏：", err))
			return
		}
		if r.Code == 0 && len(r.Data) > 0 {
			ctx.SendChain(message.Text(r.Data[0].Content)) // 输出提取后的结果
		} else {
			ctx.SendChain(message.Text("API访问错误"))
		}
	})
	en.OnPrefixGroup([]string{"维基", "百科"}, zero.SuperUserPermission).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(
		func(ctx *zero.Ctx) {
			txt := ctx.State["args"].(string)
			if txt != "" {
				pic := screenshot("https://zh.wikipedia.org/zh-cn/" + url.QueryEscape(txt))
				ctx.SendChain(message.ImageBytes(pic))
			}
		})
	en.OnPrefixGroup([]string{"百度百科"}).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(
		func(ctx *zero.Ctx) {
			txt := ctx.State["args"].(string)
			if txt != "" {
				pic := screenshot("https://baike.baidu.com/item/" + url.QueryEscape(txt))
				ctx.SendChain(message.ImageBytes(pic))
			}
		})
}

func screenshot(url string) []byte {
	pw, err := playwright.Run()
	var pic []byte
	if err != nil {
		log.Errorf("could not launch playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Errorf("could not launch Chromium: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Errorf("could not create page: %v", err)
	}
	if _, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Errorf("could not goto: %v", err)
	}
	Clip := util.PageScreenshotOptionsClip(
		playwright.PageScreenshotOptionsClip{
			X:      playwright.Float(10),
			Y:      playwright.Float(0),
			Width:  playwright.Float(1500),
			Height: playwright.Float(1400),
		})
	if pic, err = page.Screenshot(playwright.PageScreenshotOptions{
		Clip:     Clip,
		Type:     playwright.ScreenshotTypeJpeg,
		Quality:  playwright.Int(100),
		FullPage: playwright.Bool(true),
	}); err != nil {
		log.Errorf("could not create screenshot: %v", err)
	}
	if err = browser.Close(); err != nil {
		log.Errorf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Errorf("could not stop Playwright: %v", err)
	}
	return pic
}
