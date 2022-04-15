package crazy

import (
	"encoding/json"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/process"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io/ioutil"
	"math/rand"
	"os"
)

func init() { // 插件主体
	engine := control.Register("sulian", &control.Options{
		DisableOnDefault: false,
		PublicDataFolder: "Crazy",
	})
	go func() {
		dbpath := engine.DataFolder()
		dbfile := dbpath + "crazy.json"
		if file.IsNotExist(dbfile) {
			process.SleepAbout1sTo2s()
			_ = os.MkdirAll(dbpath, 0755)
			err := file.DownloadTo("https://raw.githubusercontent.com/Cha0sIDL/data/master/crazy/crazy.json",
				dbfile, false)
			if err != nil {
				panic(err)
			}
			logrus.Infoln("[sulian]加载成功")
		}
	}()
	engine.OnFullMatch("Crazy").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			dbfile := engine.DataFolder() + "crazy.json"
			data, err := ioutil.ReadFile(dbfile)
			if err != nil {
				ctx.SendChain(message.Text("读取配置文件出错了！！！"))
			}
			var temp []string
			json.Unmarshal(data, &temp)
			r := rand.Intn(len(temp))
			ctx.SendChain(message.Text(temp[r]), message.AtAll())
		})
}
