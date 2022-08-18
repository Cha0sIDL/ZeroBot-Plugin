package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/wdvxdr1123/ZeroBot/message"

	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const config = "config.json"

type Config struct {
	RpcHost   string        `json:"rpc_host"`  // http rpc的地址
	TTS       *TTS          `json:"tts"`       // 阿里tts的一些配置
	WsUrl     string        `json:"ws_url"`    // jxapi ws的地址
	Weather   string        `json:"weather"`   // 天气查询token
	Ignore    []string      `json:"ignore"`    // 忽略的触发列表
	SecretId  string        `json:"secretId"`  // 腾讯npl
	SecretKey string        `json:"secretKey"` // 腾讯npl
	Picture   *PictureToken `json:"picture"`
	JxChat    *[]Chat       `json:"jxChat"`
	SignKey   string        `json:"signKey"`
	KasKey    string        `json:"kasKey"` // 卡巴斯基软件检测的key
}

type TTS struct {
	Appkey string   `json:"appkey"`
	Access string   `json:"access"`
	Secret string   `json:"secret"`
	Voice  []string `json:"voice"`
	Start  string   `json:"start"`
}

type PictureToken struct {
	MuXiaoGuo string `json:"mu_xiao_guo,omitempty"`
	AlApi     string `json:"al_api,omitempty"`
}

type Chat struct {
	Url         string `json:"url,omitempty"`
	Name        string `json:"name,omitempty"`
	Token       string `json:"token,omitempty"`
	DeviceToken string `json:"deviceToken,omitempty"`
}

var Cfg Config

func init() {
	en := control.Register("config", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "- 加载配置文件",
	})
	en.OnKeyword("配置").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			initConfig()
		})
	en.OnFullMatch("当前配置").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(prettyPrint(Cfg)))
		})
	initConfig()
}

func initConfig() {
	tmp, err := ioutil.ReadFile(config)
	if err != nil {
		panic("读取文件失败")
	}
	Cfg = Config{TTS: &TTS{Start: time.Now().Format("2006-01-02")}}
	json.Unmarshal(tmp, &Cfg)
}

func prettyPrint(v interface{}) string {
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
