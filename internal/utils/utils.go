package utils

import (
	"Project/BlogSystem/internal/model"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	var builder strings.Builder
	for i := 0; i < length; i++ {
		num := rand.Intn(10) // 生成0到9之间的随机数
		builder.WriteString(strconv.Itoa(num))
	}
	result := builder.String()
	return result
}

// ParseKeyValue 将字符串解码
func ParseKeyValue(str string) model.KeyValue {
	parts := strings.Split(str, ":")
	if len(parts) != 2 {
		return model.KeyValue{}
	}
	return model.KeyValue{
		Key:   parts[0],
		Value: parts[1],
	}
}

// IntegrateArticleInfo 整合结构体数据字段为一个字符串
func IntegrateArticleInfo(artInfo *model.OtherArticleInfo) string {
	// 拼接字符串优化
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%v:%v:%v:%s:%s",
		artInfo.ReadAmount, artInfo.Like, artInfo.Dislike, artInfo.Class, artInfo.Tag))
	intString := builder.String()
	return intString
}

// ParseOtherArticleInfo 将ArticleInfo字符串解析为OtherArticleInfo
func ParseOtherArticleInfo(artId int, otherInfoString string) *model.OtherArticleInfo {
	// 分割字符串
	parts := strings.Split(otherInfoString, ":")
	if len(parts) != 5 {
		return &model.OtherArticleInfo{}
	}
	v1, _ := strconv.Atoi(parts[0])
	v2, _ := strconv.Atoi(parts[1])
	v3, _ := strconv.Atoi(parts[2])
	v4, _ := strconv.Unquote(`"` + parts[3] + `"`)
	v5, _ := strconv.Unquote(`"` + parts[4] + `"`)
	return &model.OtherArticleInfo{
		ArticleId:  artId,
		ReadAmount: v1,
		Like:       v2,
		Dislike:    v3,
		Class:      v4,
		Tag:        v5,
	}
}

func ParseTagFromRedisRecord(target string) *model.TagList {
	parts := strings.Split(target, ":")
	if len(parts) != 3 {
		return &model.TagList{}
	}
	v0, _ := strconv.Unquote(`"` + parts[0] + `"`)
	v1, _ := strconv.Atoi(parts[1])
	v2, _ := strconv.Atoi(parts[2])
	return &model.TagList{
		Title:     v0,
		ArticleId: v1,
		AuthorId:  v2,
	}
}

func ParseUserArticleResponse1(userId int, articleId int, value string) any {
	// 切割字符串为三部分
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return nil
	}
	v0, _ := strconv.Atoi(parts[0])
	v1, _ := strconv.ParseInt(parts[1], 10, 64)
	return &model.UserResponse1{
		UserId:        userId,
		ArticleId:     articleId,
		LikeCondition: v0,
		CCondition:    int(v1),
	}
}

// GetArticleTagSlice 获取文章的标签频次切片与Map
func GetArticleTagSlice(tagAmountMap map[string]int) ([]int, map[int]string) {
	mapLen := len(tagAmountMap)

	proList := make([]int, mapLen+1, mapLen+1) //前缀和数组
	proList[0] = 0
	proMap := make(map[int]string) //索引映射

	cnt := 1
	for k, v := range tagAmountMap {
		proList[cnt] = proList[cnt-1] + v
		proMap[cnt-1] = k
		cnt++
	}

	return proList, proMap
}

// GetRandomVal 随机获取一个标签
func GetRandomVal(proList []int) int {
	listLen := len(proList)
	maxVal := proList[listLen-1]     //	前缀和数组的最大值
	rand.Seed(time.Now().UnixNano()) //	随机数种子

	target := rand.Intn(maxVal)

	// 二分查找target左边的第一个元素
	l := -1
	r := listLen

	for l+1 != r {
		mid := (l + r) / 2
		if proList[mid] < target {
			l = mid
		} else {
			r = mid
		}
	}
	return l
}
