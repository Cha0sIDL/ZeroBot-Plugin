package util

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/smallnest/rpcx/client"

	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/golang-module/carbon/v2"

	log "github.com/sirupsen/logrus"

	"github.com/FloatTech/ZeroBot-Plugin/config"
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
	// req, err := http.NewRequest(method, httpUrl, bytes.NewReader(body))
	// if err != nil {
	//	panic("request error")
	//}
	// client := &http.Client{}
	// req.Header.Set("User-Agent", "User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)")
	// response, err := client.Do(req)
	// log.Errorln(response.Body, response.StatusCode)
	// if response.StatusCode != http.StatusOK {
	//	panic("请求失败")
	//}
	// return io.ReadAll(response.Body)
}

func ProxyHttp(client *http.Client, url, method, referer, ua string, body io.Reader) (data []byte, err error) {
	var request *http.Request
	request, err = http.NewRequest(method, "https://http-go-http-proxy-jvuuzynfbg.cn-hangzhou.fcapp.run", body)
	if err == nil {
		// 增加header选项
		if referer != "" {
			request.Header.Add("Referer", referer)
		}
		if ua != "" {
			request.Header.Add("User-Agent", ua)
		}
		request.Header.Add("proxy", url)
		var response *http.Response
		response, err = client.Do(request)
		if err == nil {
			if response.StatusCode != http.StatusOK {
				s := fmt.Sprintf("status code: %d", response.StatusCode)
				err = errors.New(s)
				return
			}
			data, err = io.ReadAll(response.Body)
			response.Body.Close()
		}
	}
	return
}

func Rand(min, max int) int {
	if min >= max || max == 0 {
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

func Interface2String(value interface{}) string {
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

// SplitSpace 按空格分隔
func SplitSpace(text string) []string {
	return strings.Fields(strings.TrimSpace(text))
}

// Deprecated: Use lo.Shuffle instead.
func Shuffle(slice interface{}) { // 切片乱序
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

// Deprecated: Use lo.Sample instead.
func RandSlice(slice interface{}) interface{} { // 随机取切片
	rv := reflect.ValueOf(slice)
	if rv.Type().Kind() != reflect.Slice {
		return slice
	}
	length := rv.Len()
	return rv.Index(rand.Intn(length)).Interface()
}

// Unicode2Zh Unicode转中文
func Unicode2Zh(sText string) []byte {
	textQuoted := strconv.QuoteToASCII(sText)
	str, _ := strconv.Unquote(strings.Replace(strconv.Quote(textQuoted), `\\u`, `\u`, -1))
	return []byte(str)
}

func Unicode2utf8(source string) string {
	var res = []string{""}
	sUnicode := strings.Split(source, "\\u")
	var context = ""
	for _, v := range sUnicode {
		var additional = ""
		if len(v) < 1 {
			continue
		}
		if len(v) > 4 {
			rs := []rune(v)
			v = string(rs[:4])
			additional = string(rs[4:])
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			context += v
		}
		context += fmt.Sprintf("%c", temp)
		context += additional
	}
	res = append(res, context)
	return strings.Join(res, "")
}

func GetWeek() string {
	s := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	intWeek := carbon.Now().Week()
	return s[intWeek]
}

// PrettyPrint 格式化打印
func PrettyPrint(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println(v)
		return ""
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		fmt.Println(v)
		return ""
	}

	return out.String()
}

// ConvertStrSlice2Map 将字符串 slice 转为 map[string]struct{}。
func ConvertStrSlice2Map(sl []string) map[string]struct{} {
	set := make(map[string]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// BytesCombine 将[]byte合并
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// RandStr 返回指定长度随机字符串
func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

// SliceDeduplicate 任意类型切片去重
// Deprecated: Use lo.Uniq instead.
func SliceDeduplicate(data interface{}) {
	dataVal := reflect.ValueOf(data)
	if dataVal.Kind() != reflect.Ptr {
		fmt.Println("input data.kind is not pointer")
		return
	}
	tmpData := Deduplicate(dataVal.Elem().Interface())
	tmpDataVal := reflect.ValueOf(tmpData)
	dataVal.Elem().Set(tmpDataVal)
}

func Deduplicate(data interface{}) interface{} {
	inArr := reflect.ValueOf(data)
	if inArr.Kind() != reflect.Slice && inArr.Kind() != reflect.Array {
		return data
	}

	existMap := make(map[interface{}]bool)
	outArr := reflect.MakeSlice(inArr.Type(), 0, inArr.Len())

	for i := 0; i < inArr.Len(); i++ {
		iVal := inArr.Index(i)

		if _, ok := existMap[iVal.Interface()]; !ok {
			outArr = reflect.Append(outArr, inArr.Index(i))
			existMap[iVal.Interface()] = true
		}
	}

	return outArr.Interface()
}

func DiffTime(start int64, end int64) string {
	diff := end - start
	if diff > 60 {
		min := diff / 60
		sec := diff % 60
		if sec > 0 {
			return fmt.Sprintf("%d分%d秒", min, sec)
		}
		return fmt.Sprintf("%d分", min)
	} else {
		return fmt.Sprintf("%d秒", diff)
	}
}

func HttpError() []message.MessageSegment {
	var msg []message.MessageSegment
	msg = append(msg, message.Text("垃圾服务器又抽风了，稍后再试试吧,,Ծ‸Ծ,,"))
	return msg
}

func GetChinese(text string) string {
	var s string
	for _, char := range text {
		if unicode.Is(unicode.Han, char) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(char))) {
			s += string(char)
		}
	}
	return s
}
