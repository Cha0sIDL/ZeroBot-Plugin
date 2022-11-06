package arknights

import (
	"archive/zip"
	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"io"
	"os"
	"path/filepath"
)

var CharTable map[string]CharData
var Rarity2CharName [][]string
var arkNightDataPath string

type CharData struct {
	Name               string `json:"name"`
	Profession         string `json:"profession"`
	Rarity             int8   `json:"rarity"`
	ItemObtainApproach string `json:"itemObtainApproach"`
}

func init() {
	Fonts, _ = gg.LoadFontFace(text.FontFile, 18)
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
				arkNightDataPath+"arknight.zip", false)
			if err != nil {
				return
			}
			unzip(arkNightDataPath+"arknight.zip", arkNightDataPath)
			log.Println("加载ArkNight数据成功")
		}
	}()
	engine.OnRegex("^查公招$").SetBlock(true).Limit(ctxext.LimitByUser).Handle(recruit)
	engine.OnKeyword("方舟资源").SetBlock(true).Limit(ctxext.LimitByUser).Handle(daily)
	engine.OnFullMatch("方舟十连").SetBlock(true).Limit(ctxext.LimitByUser).Handle(gacha)
}

func unzip(zipFile, dest string) error {
	log.Println(zipFile, dest)
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
				log.Fatalf("[unzip]: close io %s", iErr.Error())
			}
			fErr := rc.Close()
			if fErr != nil {
				log.Fatalf("[unzip]: close io %s", fErr.Error())
			}
		}
	}
	return nil
}
