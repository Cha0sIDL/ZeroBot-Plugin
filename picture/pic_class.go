package picture

var (
	modeMap = func() (m map[string]Picture) {
		setReplyMap := func(m map[string]Picture, r Picture) {
			m[r.String()] = r
		}
		m = make(map[string]Picture)
		setReplyMap(m, &XiaoGuo{})
		return
	}()
)

type Picture interface {
	Picture(msg string) string
	String() string
}

func NewPicture(mode string) Picture {
	return modeMap[mode]
}
