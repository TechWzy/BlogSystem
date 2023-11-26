package rediss

import (
	g "Project/BlogSystem/internal/global"
	"Project/BlogSystem/internal/model"
	"Project/BlogSystem/internal/utils"
	"context"
	"fmt"
	"strconv"
	"strings"
)

// InsertArticleInfo 将articleId:otherArticleInfo插入数据库
func InsertArticleInfo(artInfo *model.OtherArticleInfo) error {
	// 将文章的ID值作为键，文章字段内容拼接为字符串(值)
	key := "ArticleId"
	key = fmt.Sprintf("%s:%v", key, artInfo.ArticleId)
	intString := utils.IntegrateArticleInfo(artInfo)
	// 设置key:value值
	if err := g.Client.Set(context.Background(), key, intString, 0).Err(); err != nil {
		fmt.Println("Error setting value:", err)
		return err
	}
	return nil
}

// GetArticleInfo 获取redis集合中的OtherArticleInfo数据(已知道文章的ID)
func GetArticleInfo(articleId int) (*model.OtherArticleInfo, error) {
	key := "ArticleId"
	key = fmt.Sprintf("%s:%v", key, articleId)
	// 获取key的value值
	value, err := g.Client.Get(context.Background(), key).Result()
	if err != nil {
		return &model.OtherArticleInfo{}, err
	}
	otherArticleInfo := utils.ParseOtherArticleInfo(articleId, value)
	return otherArticleInfo, nil
}

// RedisInsertTag 将title:articleId:authId插入集合Tag:tag
func RedisInsertTag(tag string, title string, articleId int, authorId int) error {
	// 将tag数据切割，要求tag的数据格式为:"xx.xx.xx"
	tagList := strings.Split(tag, ".")

	//判断tagVal是否存在于tagList当中
	for _, tagTarget := range tagList {
		if g.TagMap[tagTarget] == true {
			tagSet := fmt.Sprintf("%s:%s", "Tag", tagTarget)
			valString := fmt.Sprintf("%s:%v:%v", title, articleId, authorId)
			err := g.Client.SAdd(context.Background(), tagSet, valString).Err()
			if err != nil {
				fmt.Println("标签插入出错...")
				return err
			}
		}
	}
	return nil
}

// SearchArticleByTag 已知tag,获取tagList
func SearchArticleByTag(tag string) []model.TagList {
	tagSet := fmt.Sprintf("%s:%s", "Tag", tag)
	tagResultList, err := g.Client.SMembers(context.Background(), tagSet).Result()
	if err != nil {
		return []model.TagList{}
	}
	var tagList []model.TagList
	for _, recordString := range tagResultList {
		v := utils.ParseTagFromRedisRecord(recordString)
		tagList = append(tagList, *v)
	}
	return tagList
}

// RedisIncreaseArticleAmountById 根据文章Id更新文章的阅读适量
func RedisIncreaseArticleAmountById(info *model.OtherArticleInfo, increment int) error {
	info.ReadAmount = info.ReadAmount + increment
	err := InsertArticleInfo(info)
	return err
}

// RedisIsArticleExist 已知Id值，判断该文章是否存在
func RedisIsArticleExist(articleId int) bool {
	_, err := GetArticleInfo(articleId)
	if err != nil {
		return false
	} else {
		return true
	}
}

// GetACollectionInfo 获取某一篇文章的收藏量
func GetACollectionInfo(articleId int) int {
	key := "ACollectionInfo"
	key = fmt.Sprintf("%s:%v", key, articleId)
	v, err := g.Client.Get(context.Background(), key).Result()
	if err != nil {
		return 0
	}
	var amount int64
	amount, _ = strconv.ParseInt(v, 10, 64) //假设不存在转换错误问题
	res := int(amount)
	return res
}

// UpdateACollectionInfo 更新一篇文章的收藏数量
func UpdateACollectionInfo(articleId int, increment int) error {
	res := GetACollectionInfo(articleId)
	if res+increment < 0 {
		return fmt.Errorf("wrong collectionInfo")
	}

	// 更新文章的收藏数量
	key := "ACollectionInfo"
	key = fmt.Sprintf("%s:%v", key, articleId)
	value := fmt.Sprintf("%v", res+increment)
	err := g.Client.Set(context.Background(), key, value, 0).Err()
	return err
}
