package picture

import (
	"github.com/FloatTech/ZeroBot-Plugin/util"
	"math/rand"
)

var (
	modeMap = func() (m map[string]Picture) {
		setReplyMap := func(m map[string]Picture, r Picture) {
			m[r.String()] = r
		}
		m = make(map[string]Picture)
		setReplyMap(m, &XiaoGuo{})
		setReplyMap(m, &Al{})
		return
	}()
)

type Picture interface {
	Picture(msg string) (data []string, err error)
	String() string
}

func NewPicture(mode string) Picture {
	return modeMap[mode]
}

func GetPicture(msg string) (url string) {
	kind := []string{"木小果", "Al"}
	util.Shuffle(kind)
	for _, k := range kind {
		p := modeMap[k]
		data, err := p.Picture(msg)
		if err != nil {
			continue
		}
		return data[rand.Intn(len(data))]
	}
	return
}
