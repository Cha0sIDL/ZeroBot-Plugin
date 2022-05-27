package picture

import (
	"fmt"
	"github.com/FloatTech/zbputils/web"
	"net/http"
	"net/url"
)

type XiaoGuo struct{}

const (
	xiaoguoURL = "https://api.muxiaoguo.cn/api/emoticons?tuname=%s"
)

func (*XiaoGuo) String() string {
	return "木小果"
}

// TalkPlain 取得图片信息
func (*XiaoGuo) Picture(msg string) string {
	u := fmt.Sprintf(xiaoguoURL, url.QueryEscape(msg))
	client := &http.Client{}
	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	// 自定义Header
	req.Header.Set("User-Agent", web.RandUA())
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "81.70.100.130")
	resp, err := client.Do(req)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	defer resp.Body.Close()
	return ""
}
