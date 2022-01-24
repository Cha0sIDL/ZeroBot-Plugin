package util

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"time"
)

type Args struct {
	Method string
	Url    string
	Body   []byte
}

func SendHttp(httpUrl string, body []byte) ([]byte, error) {
	d, _ := client.NewPeer2PeerDiscovery("tcp@"+"www.cha0sidl.xyz:8888", "")
	option := client.DefaultOption
	option.Heartbeat = true
	option.HeartbeatInterval = time.Second
	option.MaxWaitForHeartbeat = 2 * time.Second
	option.IdleTimeout = 3 * time.Second
	xclient := client.NewXClient("Http", client.Failtry, client.RandomSelect, d, option)
	defer xclient.Close()
	args := &Args{
		Method: "GET",
		Url:    httpUrl,
		Body:   body,
	}
	var Reply []byte
	err := xclient.Call(context.Background(), "Send", args, &Reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}
	return Reply, err
	//req, err := http.NewRequest(method, httpUrl, bytes.NewReader(body))
	//if err != nil {
	//	panic("request error")
	//}
	//client := &http.Client{}
	//req.Header.Set("User-Agent", "User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)")
	//response, err := client.Do(req)
	//log.Errorln(response.Body, response.StatusCode)
	//if response.StatusCode != http.StatusOK {
	//	panic("请求失败")
	//}
	//return io.ReadAll(response.Body)
}
