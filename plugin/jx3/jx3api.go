package jx3

// JxApi Ws
import (
	"github.com/FloatTech/floatbox/process"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"time"

	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

type wsClient struct {
	url  string // ws连接地址
	conn *websocket.Conn
}

func startWs() {
	ws := &wsClient{
		url: "wss://socket.nicemoe.cn",
	}
	ws.connect()
	ws.listen()
}

func (ws *wsClient) connect() {
	var err error
RETRY:
	conn, res, err := websocket.DefaultDialer.Dial(ws.url, nil)
	for err != nil {
		log.Warnf("连接JXApi Websocket服务器时出现错误: %v", err)
		time.Sleep(2 * time.Second) // 等待两秒后重新连接
		goto RETRY
	}
	ws.conn = conn
	defer res.Body.Close()
	log.Infof("连接JXApi Websocket服务器成功")
}

func (ws *wsClient) listen() {
	for {
		t, payload, err := ws.conn.ReadMessage()
		if err != nil { // reconnect
			log.Warn("JXApi Websocket服务器连接断开...")
			time.Sleep(time.Millisecond * time.Duration(3))
			ws.connect()
		}
		if t == websocket.TextMessage {
			rsp := gjson.Parse(helper.BytesToString(payload))
			log.Println("收到JXApi推送", helper.BytesToString(payload))
			go sendNotice(rsp)
		}
	}
}

func sendNotice(payload gjson.Result) {
	var rsp []message.MessageSegment
	now := time.Now().Hour()
	if now >= 0 && now < 6 { //十二点之后不响应
		return
	}
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		controls := jdb.isEnable()
		log.Println("sendNotice controls ", controls, "data", payload.Get("data.server"))
		for _, g := range ctx.GetGroupList().Array() {
			grp := g.Get("group_id").Int()
			if server, ok := controls[grp]; ok {
				switch payload.Get("action").Int() {
				case 2004:
					log.Println("sendNotice grp ", controls[grp], "data", payload.Get("data.server").String(), "grp", grp)
					if server == payload.Get("data.server").String() || payload.Get("data.server").String() == "-" {
						rsp =
							[]message.MessageSegment{
								message.Text(payload.Get("data.title").String() + "\n" +
									payload.Get("data.url").String() + "\n" + payload.Get("data.date").String()),
							}
					}
				}
				if len(rsp) != 0 {
					ctx.SendGroupMessage(grp, rsp)
					process.SleepAbout1sTo2s()
				}
			}
		}
		return true
	})
}
