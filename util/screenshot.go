package util

import (
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
)

// ScreenShot 截取浏览器
func ScreenShot(url string, option ...playwright.PageScreenshotOptionsClip) []byte {
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
	PageScreenshotOptions :=
		playwright.PageScreenshotOptions{
			Type:     playwright.ScreenshotTypeJpeg,
			Quality:  playwright.Int(100),
			FullPage: playwright.Bool(true),
		}
	if len(option) != 0 {
		PageScreenshotOptions.Clip = PageScreenshotOptionsClip(option[0])
	}
	if pic, err = page.Screenshot(PageScreenshotOptions); err != nil {
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
