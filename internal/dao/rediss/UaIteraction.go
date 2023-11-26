package rediss

import (
	g "Project/BlogSystem/internal/global"
	"Project/BlogSystem/internal/model"
	"Project/BlogSystem/internal/utils"
	"context"
	"fmt"
	"strconv"
)

// Interaction1相关数据库操作

// InsertUserResponse1 插入用户文章交互数据
func InsertUserResponse1(response1 *model.UserResponse1) error {
	// 设置key值，key值的前缀为:ua1:userId:articleId
	key := "ua1"
	key = fmt.Sprintf("%s:%v:%v", key, response1.UserId, response1.ArticleId)
	// 设置value值
	value := fmt.Sprintf("%v:%v", response1.LikeCondition, response1.CCondition)
	// 向redis数据库插入key:value,如果key已经存在，那么新value将会更新旧的value
	if err := g.Client.Set(context.Background(), key, value, 0).Err(); err != nil {
		fmt.Println("Error setting value:", err)
		return err
	}
	return nil
}

// GetUserResponse1 查询用户文章交互数据(点赞，收藏，踩)
func GetUserResponse1(userId int, articleId int) (any, error) {
	//设置key值
	key := "ua1"
	key = fmt.Sprintf("%s:%v:%v", key, userId, articleId)
	// 获取值字符串
	v, err := g.Client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	value := utils.ParseUserArticleResponse1(userId, articleId, v)
	if value == nil {
		return nil, fmt.Errorf("fails to parse")
	}
	return value, nil
}

// InsertUACollectionInfo 插入用户文章交互收藏信息
func InsertUACollectionInfo(userId int, articleId int, cond int) error {
	if cond == 3 {
		return nil
	}
	// value限制在1或2
	key := "UACollected"
	key = fmt.Sprintf("%s:%v:%v", key, userId, articleId)
	value := fmt.Sprintf("%v", cond)
	if err := g.Client.Set(context.Background(), key, value, 0).Err(); err != nil {
		return err
	}
	return nil
}

// GetUACollectionInfo 获取用户文章交互收藏信息
func GetUACollectionInfo(userId int, articleId int) (int, error) {
	// key:value记录该用户对于文章的收藏状态,1表示已经收藏，2表示没有收藏
	key := "UACollected"
	key = fmt.Sprintf("%s:%v:%v", key, userId, articleId)
	v, err := g.Client.Get(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	var v1 int64
	v1, _ = strconv.ParseInt(v, 10, 64)
	cond := int(v1)
	return cond, nil
}
