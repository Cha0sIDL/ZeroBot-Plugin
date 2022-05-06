package repeater

import (
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"sync"
)

const (
	serviceName     = "repeater"
	maxRepeat       = 4 // 触发自动复读次数
	msgTraceBackNum = 8 // 复读消息追溯条数，只根据前 5 条判定是否复读
)

type autoCopy struct {
	*sync.Mutex // 群组信息更新互斥锁
	*groupMsg
}

// 循环消息队列，记录最近的 msgTraceBackNum 条消息
type groupMsg struct {
	msgMap   map[int64][]string
	indexMap map[int64]int
	*sync.Mutex
}

var instance *autoCopy

type msgCompareFunc func(string, string) bool

func init() {
	engine := control.Register(serviceName, &control.Options{
		DisableOnDefault: false,
	})
	instance = &autoCopy{}
	instance.groupMsg = newGroupMsg()
	engine.OnMessage(zero.OnlyGroup).SetBlock(false).
		Handle(func(ctx *zero.Ctx) {
			instance.autoCopyAndJoinIn(ctx)
		})
}

// 判断群消息是否是复读，如果是，则加入
func (m *autoCopy) autoCopyAndJoinIn(ctx *zero.Ctx) {
	msgStr := ctx.ExtractPlainText()
	if len(msgStr) == 0 {
		return
	}
	// 没有复读消息，直接返回
	if !m.isMsgRepeat(ctx.Event.GroupID, msgStr, strictCompare) {
		return
	}
	// 复读，清空复读消息历史记录，并开始复读
	m.groupMsg.reset(ctx.Event.GroupID)
	ctx.SendChain(message.Text(msgStr))
}

// 判断用户消息是否属于复读，并更新消息队列
func (g *groupMsg) isMsgRepeat(groupCode int64, msg string, same msgCompareFunc) bool {
	g.Lock()
	defer g.Unlock()
	if _, ok := g.msgMap[groupCode]; !ok {
		g.msgMap[groupCode] = make([]string, msgTraceBackNum)
		g.indexMap[groupCode] = 0
	}
	// 遍历循环消息队列，计算前面已经出现过的相同的消息数量
	appearedTimes := 0
	for _, v := range g.msgMap[groupCode] {
		if same(v, msg) {
			appearedTimes++
			// 本次消息也要计算在内，所以减一
			if appearedTimes >= maxRepeat-1 {
				break
			}
		}
	}
	// 更新消息队列记录
	g.msgMap[groupCode][g.indexMap[groupCode]] = msg
	g.indexMap[groupCode]++
	if g.indexMap[groupCode] >= msgTraceBackNum {
		g.indexMap[groupCode] = 0
	}
	return appearedTimes >= maxRepeat-1
}

func (g *groupMsg) reset(groupCode int64) {
	g.Lock()
	defer g.Unlock()
	delete(g.msgMap, groupCode)
	delete(g.indexMap, groupCode)
	g.msgMap[groupCode] = make([]string, msgTraceBackNum)
	g.indexMap[groupCode] = 0
}

// 严格对比
func strictCompare(msg1, msg2 string) bool {
	return msg1 == msg2
}

func newGroupMsg() *groupMsg {
	g := groupMsg{
		msgMap:   make(map[int64][]string),
		indexMap: make(map[int64]int),
	}
	g.Mutex = &sync.Mutex{}
	return &g
}
