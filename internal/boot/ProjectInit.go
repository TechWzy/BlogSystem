package boot

import g "Project/BlogSystem/internal/global"

// InitProject  初始化各类型数据
func InitProject() {
	for _, tag := range g.TagLists {
		g.TagMap[tag] = true
	}
	// 初始化消息ID
	g.MessageID = 1
}
