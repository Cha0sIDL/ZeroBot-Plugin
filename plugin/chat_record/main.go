// Package chatrecord  聊天记录存储
package chatrecord

import (
	"os"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
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
		err := db.Open(time.Hour * 24)
		if err != nil {
			panic(err)
		}
	}()
	engine.OnMessage(zero.OnlyGroup).SetBlock(false).Handle(
		func(ctx *zero.Ctx) {
			go func() {
				m.Lock()
				defer m.Unlock()
				gidStr := strconv.FormatInt(ctx.Event.GroupID, 10)
				err := db.Create(gidStr, &record{})
				if err != nil {
					return
				}
				err = db.Insert(gidStr, &record{
					MId:     ctx.Event.MessageID,
					Message: ctx.Event.RawMessage,
					GroupID: ctx.Event.GroupID,
					Time:    ctx.Event.Time,
					UserID:  ctx.Event.UserID,
				})
				if err != nil {
					log.Errorln("Err", err)
				}
			}()
		},
	)
}
