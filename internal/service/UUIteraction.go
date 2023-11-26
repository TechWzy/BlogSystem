package service

import (
	"Project/BlogSystem/internal/dao/mysqls"
	"Project/BlogSystem/internal/dao/rediss"
	g "Project/BlogSystem/internal/global"
	"Project/BlogSystem/internal/logic"
	"Project/BlogSystem/internal/model"
	"fmt"
	"github.com/go-redis/redis/v8"
)

// AddFollow 添加粉丝
func AddFollow(idolId int, followId int) (string, error) {
	flag, err := rediss.IsFollow(idolId, followId)

	//	err报出其他错误，不存在err!=nil&&flag==true的情况
	if err != nil && err != redis.Nil {
		return "添加失败", err
	} else if err == redis.Nil || flag == false {
		if e := rediss.AddFollow(idolId, followId); e != nil {
			return "添加失败或集合创建失败...", e
		}
	} else {
		return "粉丝已存在...", nil
	}

	// 粉丝数量添加1
	if e := rediss.IsIdolKeyExist(idolId); e != nil {
		return "无法创建key值", e
	}
	if e := rediss.AdjustFollowNumber(idolId, 1); e != nil {
		return "粉丝数量增添失败...", e
	}

	// 添加成功后，向粉丝发送一条感谢信
	var letter model.PrivateLetter
	letter.BaseMessage.SId = idolId
	letter.BaseMessage.RId = followId
	letter.BaseMessage.Expiration = 0
	letter.BaseMessage.ID = g.MessageID
	letter.Subject = "《粉丝感谢信》"

	message, e := mysqls.GetFollowMessage(idolId)
	if e != nil {
		return "无法获取粉丝感谢信...", nil
	} else {
		letter.Content = message
	}
	if e := rediss.InsertPrivateLetter(&letter); e != nil {
		return "无法插入私信", e
	}
	return "添加成功...", nil
}

// DelFollow 粉丝取消关注,无需通知Idol
func DelFollow(idolId int, followId int) (string, error) {
	flag, err := rediss.IsFollow(idolId, followId)
	if err != nil && err != redis.Nil {
		return "调用IsFollow之后报错了...", err
	} else if err == redis.Nil {
		return "集合压根不存在...", err
	} else if flag == false {
		followName, _ := mysqls.GetSenderNameById(followId)
		return followName + "不是你的粉丝...", nil
	} else {
		if e := rediss.DeleteFollow(idolId, followId); e != nil {
			return "删除失败", e
		} else {
			if e := rediss.AdjustFollowNumber(idolId, -1); e != nil {
				return "无法更新粉丝量...", e
			} else {
				return "取关成功...", e
			}
		}
	}
}

// RecommendArticle 向每一位粉丝推荐自己的文章
func RecommendArticle(idolId int, article *model.UploadArticle) {
	// 先获取所有的粉丝列表
	followList := logic.GetFollowList(idolId)
	if len(followList) == 0 {
		return
	}
	combineString := logic.GenerateRecommendContent(idolId, article)
	for _, v := range followList {
		name, _ := mysqls.GetSenderNameById(v)
		combineString = "你好," + name + combineString
		var letter model.PrivateLetter
		letter.BaseMessage.RId = v
		letter.BaseMessage.SId = idolId
		letter.BaseMessage.Expiration = 0
		letter.BaseMessage.ID = g.MessageID
		letter.Subject = "文章推荐信"
		letter.Content = combineString
		fmt.Println(letter)
		e := rediss.InsertPrivateLetter(&letter)
		fmt.Println(e)
	}
}
