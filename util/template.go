// Package util 工具函数
package util

import (
	"github.com/FloatTech/floatbox/file"
	"github.com/flosch/pongo2/v5"
	"github.com/playwright-community/playwright-go"
)

// Template2html 模板转html
func Template2html(templateName string, arg pongo2.Context) string {
	tpl, err := pongo2.FromFile("./template/" + templateName)
	if err != nil {
		panic(err)
	}
	out, err := tpl.Execute(arg)
	if err != nil {
		panic(err)
	}
	return out
}

// HTML2pic  html转图片
func HTML2pic(dataPath string, fileName string, html string, clip ...*playwright.PageScreenshotOptionsClip) (finName string, err error) {
	pw, err := playwright.Run()
	finName = ""
	if err != nil {
		playwright.Install() //nolint:errcheck
		playwright.Run()     //nolint:errcheck
	}
	defer pw.Stop() //nolint:errcheck
	browser, err := pw.Chromium.Launch()
	if err != nil {
		playwright.Install() //nolint:errcheck
	}
	page, err := browser.NewPage(playwright.BrowserNewContextOptions{
		BaseURL: playwright.String("file:///" + file.Pwd() + "/template/"),
	})
	if err != nil {
		return
	}
	_, err = page.Goto("")
	if err != nil {
		return
	}
	page.SetDefaultNavigationTimeout(60 * 1000)
	err = page.SetContent(html, playwright.PageSetContentOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return
	}
	page.WaitForTimeout(10)
	page.QuerySelector("#main") //nolint:errcheck
	finName = dataPath + fileName + ".jpeg"
	PageScreenshotOptions := playwright.PageScreenshotOptions{
		Path:     playwright.String(finName),
		Type:     playwright.ScreenshotTypeJpeg,
		Quality:  playwright.Int(100),
		FullPage: playwright.Bool(true),
	}
	if len(clip) != 0 {
		PageScreenshotOptions.Clip = clip[0]
	}
	_, err = page.Screenshot(PageScreenshotOptions)
	return finName, err
}

// PageScreenshotOptionsClip 返回PageScreenshotOptionsClip的指针
func PageScreenshotOptionsClip(v playwright.PageScreenshotOptionsClip) *playwright.PageScreenshotOptionsClip {
	return &v
}
