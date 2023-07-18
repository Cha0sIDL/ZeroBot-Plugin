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
			sendNotice(rsp)
		}
	}
}

func sendNotice(payload gjson.Result) {
	var rsp []message.MessageSegment
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		controls := jdb.isEnable()
		for _, g := range ctx.GetGroupList().Array() {
			grp := g.Get("group_id").Int()
			if _, ok := controls[grp]; ok {
				switch payload.Get("action").Int() {
				case 2004:
					rsp =
						[]message.MessageSegment{
							message.Text("【818】:\n"),
							message.Text(payload.Get("data.title").String() + "\n" +
								payload.Get("data.url").String() + "\n" + payload.Get("data.date").String()),
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
