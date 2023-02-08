// Package util 工具函数tts
package util

import (
	"errors"
	"io"
	"os"
	"sync"
	"time"

	nls "github.com/aliyun/alibabacloud-nls-go-sdk"
	log "github.com/sirupsen/logrus"
)

type ttsUserParam struct {
	f io.Writer
}

func onTaskFailed(text string, param interface{}) {
	_, ok := param.(*ttsUserParam)
	if !ok {
		log.Fatal("text filed:", text)
		return
	}
}

func onSynthesisResult(data []byte, param interface{}) {
	p, ok := param.(*ttsUserParam)
	if !ok {
		return
	}
	p.f.Write(data) //nolint:errcheck
}

func onCompleted(text string, param interface{}) {
	_, ok := param.(*ttsUserParam)
	if !ok {
		log.Fatal("invalid logger")
		return
	}
}

func onClose(param interface{}) {
	_, ok := param.(*ttsUserParam)
	if !ok {
		log.Fatal("invalid logger")
		return
	}
}

func waitReady(ch chan bool) error {
	select {
	case done := <-ch:
		{
			if !done {
				return errors.New("wait failed")
			}
		}
	case <-time.After(60 * time.Second):
		{
			return errors.New("wait timeout")
		}
	}
	return nil
}

var lk sync.Mutex
var fail = 0
var reqNum = 0

// TTS 阿里tts接口
func TTS(fileName string, text string, param nls.SpeechSynthesisStartParam, appKEY string, aKid string, akKey string) error {
	config, err := nls.NewConnectionConfigWithAKInfoDefault(nls.DEFAULT_URL, appKEY, aKid, akKey)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ttsUserParam := new(ttsUserParam)
		fout, err := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755) //nolint
		ttsUserParam.f = fout
		tts, err := nls.NewSpeechSynthesis(config, nil, false,
			onTaskFailed, onSynthesisResult, nil,
			onCompleted, onClose, ttsUserParam)
		if err != nil {
			return
		}
		lk.Lock()
		reqNum++
		lk.Unlock()
		ch, err := tts.Start(text, param, nil)
		if err != nil {
			lk.Lock()
			fail++
			lk.Unlock()
			tts.Shutdown()
		}

		err = waitReady(ch)
		if err != nil {
			lk.Lock()
			fail++
			lk.Unlock()
			tts.Shutdown()
		}
		tts.Shutdown()
	}()
	wg.Wait()
	return nil
}
