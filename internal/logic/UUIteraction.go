package logic

import (
	"Project/BlogSystem/internal/dao/mysqls"
	"Project/BlogSystem/internal/dao/rediss"
	"Project/BlogSystem/internal/model"
	"strconv"
	"strings"
)

// GetFollowList 获取粉丝列表
func GetFollowList(idolId int) []int {
	strList, err := rediss.GetFollowList(idolId)
	if err != nil {
		return []int{}
	}

	var followList []int
	for _, val := range strList {
		v, _ := strconv.ParseInt(val, 10, 64)
		followList = append(followList, int(v))
	}
	return followList
}

// GenerateRecommendContent 生成文章的推荐内容
func GenerateRecommendContent(idolId int, article *model.UploadArticle) string {
	tagList := strings.Split(article.Tag, ".")
	tagStr := ""
	num := len(tagList)
	for k, v := range tagList {
		if k == num-1 {
			tagStr += v
		} else {
			tagStr += v + ","
		}
	}

	name, _ := mysqls.GetSenderNameById(idolId)

	combineString := ",我是" + name + ".我刚刚发表了一篇题目为《" + article.Title + "》的博文,这是一篇关于" + tagStr + "的文章，希望能够契合你的兴趣点，也希望你能够多予斧正，谢谢!"
	return combineString
}
