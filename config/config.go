package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

const config = "config.json"

type Config struct {
	RpcHost string `json:"rpc_host"`
	TTS     *TTS   `json:"tts"`
}

type TTS struct {
	Appkey string   `json:"appkey"`
	Access string   `json:"access"`
	Secret string   `json:"secret"`
	Voice  []string `json:"voice"`
	Start  string   `json:"start"`
}

var Cfg Config

func Init() {
	tmp, err := ioutil.ReadFile(config)
	if err != nil {
		panic("读取文件失败")
	}
	Cfg = Config{TTS: &TTS{Start: time.Now().Format("2006-01-02")}}
	json.Unmarshal(tmp, &Cfg)
	log.Println("读取配置成功\n", Cfg.RpcHost, "\n", Cfg.TTS)
}
