package active

var (
	ModeMap = func() (M map[string]AIReply) {
		setReplyMap := func(m map[string]AIReply, r AIReply) {
			m[r.String()] = r
		}
		M = make(map[string]AIReply, 3)
		setReplyMap(M, &QYKReply{})
		setReplyMap(M, &XiaoAiReply{})
		setReplyMap(M, &Tencent{})
		return
	}()
)

// AIReply 公用智能回复类
type AIReply interface {
	// Talk 取得带 CQ 码的回复消息
	Talk(msg, nickname string) string
	// Talk 取得文本回复消息
	TalkPlain(msg, nickname string) string
	// String 获得实际使用的回复服务名
	String() string
}

// NewAIReply 智能回复简单工厂
func NewAIReply(mode string) AIReply {
	return ModeMap[mode]
}
