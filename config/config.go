package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

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
	control.Register("config", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "- 加载配置文件",
	}).OnKeyword("配置").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			initConfig()
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
	log.Println("读取配置成功\n", Cfg.RpcHost, "\n", Cfg.TTS, "\n", Cfg.WsUrl)
}
