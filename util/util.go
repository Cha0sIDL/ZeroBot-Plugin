package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/config"
	log "github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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

func Rand(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

func Decimal(value float64, num int) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(num)+"f", value), 64)
	return value
}

func JsonToMap(jsonStr string) map[string]interface{} {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil
	}
	return m
}

func MergeMap(mObj ...map[string]interface{}) map[string]interface{} {
	newObj := make(map[string]interface{})
	for _, m := range mObj {
		for k, v := range m {
			newObj[k] = v
		}
	}
	return newObj
}

func GetCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

func TodayFileName() string {
	t := time.Now()
	return fmt.Sprint(t.Format("2006-01-02"))
}
