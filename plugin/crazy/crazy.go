package crazy

import (
	"math/rand"
	"os"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/process"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/util"
)

func init() { // 插件主体
	engine := control.Register("crazy", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		PublicDataFolder: "Crazy",
		Help: "选择困难症帮手\n" +
			"- 喝什么\n" +
			"- 吃什么\n",
	})
	//"https://raw.fastgit.org/MinatoAquaCrews/nonebot_plugin_crazy_thursday/beta/nonebot_plugin_crazy_thursday/post.json"
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
		db.Open(time.Hour * 24)
		err := db.Create("crazy", &crazy{})
		db.Create("menu", &menu{})
		db.Create("drink", &drink{})
		if err != nil {
			panic(err)
		}
	}()
	engine.OnRegex(`疯狂星期(一|二|三|四|五|六|日|天)`).SetBlock(true).
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
	word := map[string][]string{
		"breakfast": {
			"7点啦，吃早餐啦！",
			"起床啦起床啦！现在还没起床的都是懒狗！",
			"哦哈哟米娜桑！今日も元気でね！🥳",
			"新的一天又是全气满满哦！",
			"一日之计在于晨，懒狗还不起床？",
		},
		"lunch": {
			"12点啦，吃午餐啦！",
			"恰饭啦恰饭啦！再不去食堂就没吃的啦！",
			"中午还不恰点好的？整点碳水大餐嗯造吧！",
		},
		"snack": {
			"三点了，饮茶了先！",
			"摸鱼时刻，整点恰滴先~",
			"做咩啊做，真给老板打工啊！快来摸鱼！",
		},
		"dinner": {
			"6点了！不会真有人晚上加班恰外卖吧？",
			"下班咯，这不开造？",
			"当务之急是下班！",
		},
		"midnight": {
			"10点啦，整个夜宵犒劳自己吧！",
			"夜宵这不来个外卖？",
			"夜宵这不整点好的？",
		},
	}
	switch {
	case now < 6: // 凌晨
		text = word["midnight"][rand.Intn(len(word["midnight"]))] + "\n 恰"
	case now >= 6 && now < 9:
		text = word["breakfast"][rand.Intn(len(word["breakfast"]))] + "\n 恰"
	case now >= 9 && now < 14:
		text = word["lunch"][rand.Intn(len(word["lunch"]))] + "\n 恰"
	case now >= 14 && now < 18:
		text = word["snack"][rand.Intn(len(word["snack"]))] + "\n 恰"
	case now >= 18 && now < 24:
		text = word["dinner"][rand.Intn(len(word["dinner"]))] + "\n 恰"
	}
	return text
}
