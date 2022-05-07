package util

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/config"
	log "github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"reflect"
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

func strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

func Image2Base64(image image.Image) []byte {
	buffer := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	var opt jpeg.Options
	opt.Quality = 95
	_ = jpeg.Encode(encoder, image, &opt)
	err := encoder.Close()
	if err != nil {
		return nil
	}
	return buffer.Bytes()
}

// 通过map主键唯一的特性过滤重复元素
func RemoveRepByMap(slc []int) []int {
	result := []int{}
	tempMap := map[int]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

func SplitSpace(text string) []string {
	return strings.Fields(strings.TrimSpace(text))
}

func Shuffle(slice interface{}) { //切片乱序
	rv := reflect.ValueOf(slice)
	if rv.Type().Kind() != reflect.Slice {
		return
	}

	length := rv.Len()
	if length < 2 {
		return
	}

	swap := reflect.Swapper(slice)
	rand.Seed(time.Now().Unix())
	for i := length - 1; i >= 0; i-- {
		j := rand.Intn(length)
		swap(i, j)
	}
	return
}
