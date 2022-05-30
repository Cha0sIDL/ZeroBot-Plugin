package picture

import (
	"errors"
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/config"
	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
)

type Al struct{}

const (
	alURL = "https://v2.alapi.cn/api/doutu?keyword=%s&token=%s&type=%d"
)

func (*Al) String() string {
	return "Al"
}

// Picture 取得图片信息
func (*Al) Picture(msg string) (data []string, err error) {
	u := fmt.Sprintf(alURL, url.QueryEscape(msg), config.Cfg.Picture.AlApi, 7)
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
	for _, d := range gjson.Get(binary.BytesToString(jsondata), "data").Array() {
		data = append(data, d.String())
	}
	return
}
