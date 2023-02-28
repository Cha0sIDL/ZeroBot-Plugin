package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/wdvxdr1123/ZeroBot/message"
)

// ProxyHTTP 代理发送http
func ProxyHTTP(client *http.Client, url, method, referer, ua string, body io.Reader) (data []byte, err error) {
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

// Rand 获取min到max之间的随机数
func Rand(min, max int) int {
	if min >= max || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

// JSONToMap Json数据转换为map
func JSONToMap(jsonStr string) map[string]interface{} {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil
	}
	return m
}

// MergeMap 多个map数据合并
func MergeMap(mObj ...map[string]interface{}) map[string]interface{} {
	newObj := make(map[string]interface{})
	for _, m := range mObj {
		for k, v := range m {
			newObj[k] = v
		}
	}
	return newObj
}

// TodayFileName 返回今天文件的名称
func TodayFileName() string {
	t := time.Now()
	return fmt.Sprint(t.UnixMilli())
}

// Interface2String 任意类型转字符串
func Interface2String(value interface{}) string {
	return fmt.Sprint(value)
}

// Image2Base64 image转base64
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

// Unicode2Zh Unicode转中文
func Unicode2Zh(sText string) []byte {
	textQuoted := strconv.QuoteToASCII(sText)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(textQuoted), `\\u`, `\u`))
	return []byte(str)
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

// BytesCombine 将[]byte合并
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// DiffTime 格式化输出两个时间戳相差时间
func DiffTime(start int64, end int64) string {
	diff := end - start
	if diff > 60 {
		min := diff / 60
		sec := diff % 60
		if sec > 0 {
			return fmt.Sprintf("%d分%d秒", min, sec)
		}
		return fmt.Sprintf("%d分", min)
	}
	return fmt.Sprintf("%d秒", diff)
}

// HTTPError http错误的统一输出
func HTTPError() []message.MessageSegment {
	var msg []message.MessageSegment
	msg = append(msg, message.Text("垃圾服务器又抽风了，稍后再试试吧,,Ծ‸Ծ,,"))
	return msg
}

// GetChinese 获取字符串中的中文字符
func GetChinese(text string) string {
	var s string
	for _, char := range text {
		if unicode.Is(unicode.Han, char) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(char))) {
			s += string(char)
		}
	}
	return s
}

// IntersectArray 求两个切片的交集
func IntersectArray[T any](l, r []T) []T {
	var inter []T
	mp := make(map[any]bool)
	for _, s := range l {
		if _, ok := mp[s]; !ok {
			mp[s] = true
		}
	}
	for _, s := range r {
		if _, ok := mp[s]; ok {
			inter = append(inter, s)
		}
	}
	return inter
}
