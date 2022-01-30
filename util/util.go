package util

import (
	"context"
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/config"
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
	d, _ := client.NewPeer2PeerDiscovery("tcp@"+config.Cfg.RpcHost, "")
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
		log.Errorln("failed to call: %v", err)
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

func Max(l []float64) (max float64) {
	max = l[0]
	for _, v := range l {
		if v > max {
			max = v
		}
	}
	return
}

func Min(l []float64) (min float64) {
	min = l[0]
	for _, v := range l {
		if v < min {
			min = v
		}
	}
	return
}

func AppendAny(a interface{}, b interface{}) string {
	return fmt.Sprintf("%v", a) + "-" + fmt.Sprintf("%v", b)
}
