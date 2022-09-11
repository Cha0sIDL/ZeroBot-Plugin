package jx3

import (
	"net/http"
	"sync"
	"time"

	binutils "github.com/FloatTech/floatbox/binary"
	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/FloatTech/ZeroBot-Plugin/config"
)

type userWsClient struct {
	url        string          // ws连接地址
	serverName string          // 当前链接的区服名字
	conn       *websocket.Conn // The Conn type represents a WebSocket connection.
	header     http.Header
}

func startChatWs(chat config.Chat) {
	ws := &userWsClient{
		url:        chat.Url,
		serverName: chat.Name,
	}
	ws.header = make(http.Header)
	ws.header.Add("origin", "https://m.pvp.xoyo.com:18048")
	ws.header.Add("apiVersion", "3")
	ws.header.Add("deviceId", "GtyAYsYwtgU8a10RQr1NCw==")
	ws.header.Add("clientId", "1")
	ws.header.Add("deviceToken", chat.DeviceToken)
	ws.header.Add("deviceOS", "a")
	ws.header.Add("token", chat.Token)
	ws.header.Add("Host", "m.pvp.xoyo.com:18048")
	ws.connect()
	ws.listen()
}

func (ws *userWsClient) connect() {
	var err error
RETRY:
	conn, res, err := websocket.DefaultDialer.Dial(ws.url, ws.header)
	for err != nil {
		conn.Close()
		log.Warnf("连接JXChat Websocket服务器时出现错误: %v", err)
		time.Sleep(2 * time.Second) // 等待两秒后重新连接
		goto RETRY
	}
	ws.conn = conn
	defer res.Body.Close()
	log.Infof("连接JXChat Websocket服务器成功")
}

func (ws *userWsClient) listen() {
	var rw sync.RWMutex
	tableName := dbUser
	err := db.Create(tableName, &User{})
	if err != nil {
		log.Warn("jx User Db Create error", err)
		return
	}
	for {
		t, payload, err := ws.conn.ReadMessage()
		if err != nil { // reconnect
			log.Warn("JXChat Websocket服务器连接断开...")
			time.Sleep(2 * time.Second)
			ws.connect()
		}
		if t == websocket.TextMessage {
			rsp := gjson.Parse(binutils.BytesToString(payload))
			if rsp.Get("cmd").Int() != 100120 {
				continue
			}
			// log.Println("World Chat", binutils.BytesToString(payload))
			server := rsp.Get("body.msg.0.extra.CenterID.0").String()
			roleName := rsp.Get("body.msg.0.sName").String()
			if len(roleName) == 0 || len(server) == 0 {
				continue
			}
			rw.Lock()
			db.Insert(tableName, &User{
				ID:   roleName + "_" + server,
				Data: binutils.BytesToString(payload),
			})
			rw.Unlock()
			time.Sleep(time.Millisecond * 500)
		}
	}
}
