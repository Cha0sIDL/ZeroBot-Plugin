package jx3

//JxApi Ws
import (
	"time"

	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	"github.com/FloatTech/ZeroBot-Plugin/config"
)

type wsClient struct {
	url  string // ws连接地址
	conn *websocket.Conn
}

func startWs() {
	ws := &wsClient{
		url: config.Cfg.WsUrl,
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
