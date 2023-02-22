package jx3

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/FloatTech/floatbox/web"
)

// NewTimeOutDefaultClient 超时延长的客户端
func NewTimeOutDefaultClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 60 * time.Second,
		},
	}
}

// RequestDataWith 剑网小黑特殊客户端
func RequestDataWith(url string) (data []byte, err error) {
	var request *http.Request
	client := NewTimeOutDefaultClient()
	request, err = http.NewRequest("POST", url, nil)
	if err == nil {
		// 增加header选项
		request.Header.Add("X-Token", "")
		request.Header.Add("User-Agent", web.RandUA())
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
