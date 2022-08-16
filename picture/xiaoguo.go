package picture

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	binutils "github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/tidwall/gjson"

	"github.com/FloatTech/ZeroBot-Plugin/config"
)

type XiaoGuo struct{}

const (
	xiaoguoURL = "https://api.muxiaoguo.cn/api/emoticons?tuname=%s&api_key=%s"
)

func (*XiaoGuo) String() string {
	return "木小果"
}

// Picture 取得图片信息
func (*XiaoGuo) Picture(msg string) (data []string, err error) {
	u := fmt.Sprintf(xiaoguoURL, url.QueryEscape(msg), config.Cfg.Picture.MuXiaoGuo)
	client := &http.Client{}
	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		return
	}
	// 自定义Header
	req.Header.Set("User-Agent", web.RandUA())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		s := fmt.Sprintf("status code: %d", resp.StatusCode)
		err = errors.New(s)
		return
	}
	jsondata, _ := io.ReadAll(resp.Body)
	gjson.Get(binutils.BytesToString(jsondata), "data").ForEach(
		func(key, value gjson.Result) bool {
			data = append(data, value.Get("imagelink").String())
			return true
		})
	return
}
