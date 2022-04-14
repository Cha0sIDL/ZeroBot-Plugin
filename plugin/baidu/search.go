// Package baidu 百度一下
package baidu

import (
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"net/url"

	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
)

func init() {
	en := control.Register("baidu", &control.Options{
		DisableOnDefault: false,
		Help: "baidu\n" +
			"- 百度下[xxx]",
	})
	en.OnPrefixGroup([]string{"百度下", "百度"}).SetBlock(true).Limit(ctxext.LimitByGroup).
		Handle(func(ctx *zero.Ctx) {
			txt := ctx.State["args"].(string)
			if txt != "" {
				ctx.SendChain(message.Text("https://buhuibaidu.me/?s=" + url.QueryEscape(txt)))
			}
		})
	en.OnPrefixGroup([]string{"维基", "百科"}).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(
		func(ctx *zero.Ctx) {
			txt := ctx.State["args"].(string)
			if txt != "" {
				pic := screenshot("https://zh.wikipedia.org/wiki/" + url.QueryEscape(txt))
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
	if pic, err = page.Screenshot(playwright.PageScreenshotOptions{
		Type:     playwright.ScreenshotTypeJpeg,
		Quality:  playwright.Int(70),
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
