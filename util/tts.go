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

type TtsUserParam struct {
	F io.Writer
}

func onTaskFailed(text string, param interface{}) {
	_, ok := param.(*TtsUserParam)
	if !ok {
		log.Fatal("text filed:", text)
		return
	}
}

func onSynthesisResult(data []byte, param interface{}) {
	p, ok := param.(*TtsUserParam)
	if !ok {
		return
	}
	p.F.Write(data)
}

func onCompleted(text string, param interface{}) {
	_, ok := param.(*TtsUserParam)
	if !ok {
		log.Fatal("invalid logger")
		return
	}
}

func onClose(param interface{}) {
	_, ok := param.(*TtsUserParam)
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

func TTS(FileName string, text string, param nls.SpeechSynthesisStartParam, AppKEY string, AKid string, AkKey string) error {
	config, err := nls.NewConnectionConfigWithAKInfoDefault(nls.DEFAULT_URL, AppKEY, AKid, AkKey)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ttsUserParam := new(TtsUserParam)
		fout, err := os.OpenFile(FileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
		ttsUserParam.F = fout
		tts, err := nls.NewSpeechSynthesis(config, nil,
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
