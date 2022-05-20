package active

import "github.com/FloatTech/AnimeAPI/aireply"

var (
	ModeMap = func() (M map[string]aireply.AIReply) {
		setReplyMap := func(m map[string]aireply.AIReply, r aireply.AIReply) {
			m[r.String()] = r
		}
		M = make(map[string]aireply.AIReply, 3)
		setReplyMap(M, &aireply.QYKReply{})
		setReplyMap(M, &aireply.XiaoAiReply{})
		setReplyMap(M, &Tencent{})
		return
	}()
)

// NewAIReply 智能回复简单工厂
func NewAIReply(mode string) aireply.AIReply {
	return ModeMap[mode]
}
