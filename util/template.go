package util

import (
	"github.com/FloatTech/zbputils/file"
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

func Html2pic(dataPath string, fileName string, pageName string, html string) (finName string, err error) {
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
	_, err = page.Goto(pageName)
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
	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(finName),
		Type:     playwright.ScreenshotTypeJpeg,
		Quality:  playwright.Int(100),
		FullPage: playwright.Bool(true),
	})
	return finName, err
}
