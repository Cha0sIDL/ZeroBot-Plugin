// Package arknights 明日方舟相关插件
package arknights

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/process"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var charTable map[string]charData
var rarity2CharName [][]string
var arkNightDataPath string

type charData struct {
	Name               string `json:"name"`
	Profession         string `json:"profession"`
	Rarity             int8   `json:"rarity"`
	ItemObtainApproach string `json:"itemObtainApproach"`
}

func init() {
	fonts, _ = gg.LoadFontFace(text.FontFile, 18)
	engine := control.Register("arknight", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		PublicDataFolder: "ArkNights",
		Brief:            "明日方舟相关插件",
		Help: "- 查公招\n" +
			"- 方舟今日资源\n" + "- 方舟十连\n",
	})
	var err error
	arkNightDataPath = engine.DataFolder()
	go func() {
		if file.IsNotExist(arkNightDataPath+"arknight.zip") || file.IsNotExist(arkNightDataPath+"version") {
			err = file.DownloadTo("https://raw.githubusercontent.com/Cha0sIDL/data/master/arknight.zip",
				arkNightDataPath+"arknight.zip")
			if err != nil {
				return
			}
			unzip(arkNightDataPath+"arknight.zip", arkNightDataPath) //nolint:errcheck
			log.Println("加载ArkNight数据成功")
		}
	}()
	engine.OnRegex("^查公招$").SetBlock(true).Limit(ctxext.LimitByUser).Handle(recruit)
	engine.OnKeyword("方舟资源").SetBlock(true).Limit(ctxext.LimitByUser).Handle(daily)
	engine.OnFullMatchGroup([]string{"方舟十连", "方舟抽卡"}).SetBlock(true).Limit(ctxext.LimitByUser).Handle(gacha)
	engine.OnFullMatch("切换方舟卡池").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
			if !ok {
				ctx.SendChain(message.Text("找不到服务!"))
				return
			}
			gid := ctx.Event.GroupID
			if gid == 0 {
				gid = -ctx.Event.UserID
			}
			store := (storage)(c.GetData(gid))
			if store.setmode(!store.is6starsmode()) {
				process.SleepAbout1sTo2s()
				ctx.SendChain(message.Text("切换到六星卡池~"))
			} else {
				process.SleepAbout1sTo2s()
				ctx.SendChain(message.Text("切换到普通卡池~\n", "2%概率六星,10%概率5星,58%概率4星,30%概率3星"))
			}
			err := c.SetData(gid, int64(store))
			if err != nil {
				process.SleepAbout1sTo2s()
				ctx.SendChain(message.Text("ERROR: ", err))
			}
		})
}

func unzip(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			log.Fatalf("[unzip]: close reader %s", err.Error())
		}
	}()
	for _, f := range reader.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		filename := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(filename, 0755)
			if err != nil {
				return err
			}
		} else {
			w, err := os.Create(filename)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, rc)
			if err != nil {
				return err
			}
			iErr := w.Close()
			if iErr != nil {
				log.Panicf("[unzip]: close io %s", iErr.Error())
			}
			fErr := rc.Close()
			if fErr != nil {
				log.Panicf("[unzip]: close io %s", fErr.Error())
			}
		}
	}
	return nil
}
