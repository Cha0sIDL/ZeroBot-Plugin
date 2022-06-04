package chat_record

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"os"
	"sync"

	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/file"
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/FloatTech/ZeroBot-Plugin/util"
)

var m sync.Mutex

func init() { // 插件主体
	engine := control.Register("record", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		PublicDataFolder: "Record",
	})
	go func() {
		dbpath := engine.DataFolder()
		db.DBPath = dbpath + "record.db"
		if file.IsNotExist(db.DBPath) {
			_ = os.MkdirAll(dbpath, 0755)
		}
		err := db.Create("record", &record{})
		if err != nil {
			panic(err)
		}
	}()
	engine.OnMessage(func(ctx *zero.Ctx) bool {
		return zero.OnlyGroup(ctx) && util.Ignore(ctx)
	}).SetBlock(false).Handle(
		func(ctx *zero.Ctx) {
			go func() {
				m.Lock()
				defer m.Unlock()
				var dbMsg string
				for _, msg := range ctx.Event.Message {
					dbMsg += msg.String() + "#Split#"
				}
				db.Insert("record", &record{
					MId:     ctx.Event.MessageID,
					Message: dbMsg,
					GroupId: ctx.Event.GroupID,
					Time:    ctx.Event.Time,
					UserID:  ctx.Event.UserID,
				})
			}()
		},
	)
}
