package crazy

import (
	"math/rand"
	"os"
	"time"

	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/process"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/util"
)

func init() { // 插件主体
	engine := control.Register("crazy", &control.Options{
		DisableOnDefault: false,
		PublicDataFolder: "Crazy",
		Help: "喝什么\n" +
			"吃什么\n",
	})
	go func() {
		dbpath := engine.DataFolder()
		db.DBPath = dbpath + "crazy.db"
		if file.IsNotExist(db.DBPath) {
			process.SleepAbout1sTo2s()
			_ = os.MkdirAll(dbpath, 0755)
			err := file.DownloadTo("https://raw.githubusercontent.com/Cha0sIDL/data/master/crazy/crazy.db",
				db.DBPath, false)
			if err != nil {
				panic(err)
			}
		}
		err := db.Create("crazy", &crazy{})
		db.Create("menu", &menu{})
		db.Create("drink", &drink{})
		if err != nil {
			panic(err)
		}
	}()
	engine.OnFullMatch("Crazy").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var t crazy
			db.Pick("crazy", &t)
			ctx.SendChain(message.Text(t.Crazy), message.AtAll())
		})
	engine.OnKeyword("吃什么").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var t menu
			db.Pick("menu", &t)
			ctx.SendChain(message.Text(now() + t.Menu))
		})
	engine.OnKeyword("喝什么").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("制作中...."))
			time.Sleep(2 * time.Second)
			var msg string
			var d drink
			var dSlice []drink
			db.Find("drink", &d, "WHERE kind='temperature' ORDER BY RANDOM() limit 1") // 温度
			msg += d.Drink + "/"
			db.Find("drink", &d, "WHERE kind = 'sugar' ORDER BY RANDOM() limit 1") // 糖
			msg += d.Drink + "/"
			db.FindFor("drink", &d, "WHERE kind = 'addon'", func() error {
				dSlice = append(dSlice, d)
				return nil
			})
			util.Shuffle(dSlice)
			for i := 0; i < rand.Intn(len(dSlice)); i++ {
				msg += dSlice[i].Drink
			}
			db.Find("drink", &d, "WHERE kind = 'body' ORDER BY RANDOM() limit 1") // 主体
			msg += d.Drink
			ctx.SendChain(message.Text(msg))
		})
}

func now() string {
	var text string
	now := time.Now().Hour()
	switch {
	case now < 6: // 凌晨
		text = "凌晨了，还在这吃"
	case now >= 6 && now < 9:
		text = "早上吃"
	case now >= 9 && now < 18:
		text = "中午吃"
	case now >= 18 && now < 24:
		text = "晚上吃"
	}
	return text
}
