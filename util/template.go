package util

import (
	"github.com/FloatTech/floatbox/file"
	"github.com/flosch/pongo2/v5"
	"github.com/playwright-community/playwright-go"
)

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

func Html2pic(dataPath string, fileName string, html string, Clip ...*playwright.PageScreenshotOptionsClip) (finName string, err error) {
	pw, err := playwright.Run()
	finName = ""
	if err != nil {
		playwright.Install()
		playwright.Run()
	}
	defer pw.Stop()
	browser, err := pw.Chromium.Launch()
	if err != nil {
		playwright.Install()
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
	err = page.SetContent(html, playwright.PageSetContentOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return
	}
	page.WaitForTimeout(10)
	page.QuerySelector("#main")
	finName = dataPath + fileName + ".jpeg"
	PageScreenshotOptions := playwright.PageScreenshotOptions{
		Path:     playwright.String(finName),
		Type:     playwright.ScreenshotTypeJpeg,
		Quality:  playwright.Int(100),
		FullPage: playwright.Bool(true),
	}
	if len(Clip) != 0 {
		PageScreenshotOptions.Clip = Clip[0]
	}
	_, err = page.Screenshot(PageScreenshotOptions)
	return finName, err
}

func PageScreenshotOptionsClip(v playwright.PageScreenshotOptionsClip) *playwright.PageScreenshotOptionsClip {
	return &v
}
