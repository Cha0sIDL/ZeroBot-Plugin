package config

import (
	"encoding/json"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"io/ioutil"
	"log"
	"time"
)

const config = "config.json"

type Config struct {
	RpcHost   string   `json:"rpc_host"`
	TTS       *TTS     `json:"tts"`
	WsUrl     string   `json:"ws_url"`
	Weather   string   `json:"weather"`
	Ignore    []string `json:"ignore"`
	SecretId  string   `json:"secretId"`
	SecretKey string   `json:"secretKey"`
}

type TTS struct {
	Appkey string   `json:"appkey"`
	Access string   `json:"access"`
	Secret string   `json:"secret"`
	Voice  []string `json:"voice"`
	Start  string   `json:"start"`
}

var Cfg Config

func init() {
	control.Register("config", &control.Options{
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
