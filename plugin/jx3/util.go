// Package jx3 一些工具函数
package jx3

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/config"
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"github.com/FloatTech/floatbox/web"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/golang-module/carbon/v2"
	"github.com/tidwall/gjson"
	"strconv"
)

func ts() string {
	return carbon.Now().Layout("20060102150405", carbon.UTC) + util.Interface2String(carbon.Now(carbon.UTC).Millisecond())
}

func sign(data interface{}) string {
	bData, _ := json.Marshal(data)
	CombineData := util.BytesCombine(bData, []byte("@#?.#@"))
	key := []byte(config.Cfg.SignKey)
	h := hmac.New(sha256.New, key)
	h.Write(CombineData)
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

// 51.2345.67.89
func price2hRead(price int64) (readStr string) {
	strPrice := strconv.FormatInt(price, 10)
	l := len(strPrice)
	for idx, str := range strPrice {
		i := l - idx
		readStr += string(str)
		switch {
		case i == 9:
			readStr += "金砖"
		case i == 5:
			readStr += "金"
		case i == 3:
			readStr += "银"
		}
	}
	readStr += "铜"
	return
}

func newPage() *components.Page {
	p := components.NewPage()
	_, err := web.GetData("http://localhost:8083/assets")
	if err != nil {
		return p
	}
	p.AssetsHost = "http://localhost:8083/assets/"
	return p
}

func average(price gjson.Result) string {
	var a float64
	price.ForEach(
		func(key, value gjson.Result) bool {
			a += value.Float()
			return true
		})
	return fmt.Sprintf("%.2f", a/price.Get("#").Float())
}
