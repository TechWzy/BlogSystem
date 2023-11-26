package mysqls

import (
	g "Project/BlogSystem/internal/global"
	"Project/BlogSystem/internal/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// UpdateLikeAndDislike 更新某篇文章的Like与dislike
func UpdateLikeAndDislike(articleId int, delta1 int, delta2 int) error {
	err := g.Db.Model(&model.OtherArticleInfo{}).Where("article_id = ?", articleId).Updates(map[string]interface{}{
		"like":    gorm.Expr("`like`+?", delta1),
		"dislike": gorm.Expr("dislike+?", delta2),
	}).Error
	return err
}

// InsertArticleComment 文章评论插入数据库
func InsertArticleComment(comment *model.ArticleComment) error {
	e := g.Db.Create(&model.ArticleComment{ArticleId: comment.ArticleId, UserId: comment.UserId, Comment: comment.Comment}).Error
	return e
}

// GetArticleCommentList 查询一篇文章的所有评论
func GetArticleCommentList(articleId int) ([]model.ArticleComment, error) {
	var articleCommentList []model.ArticleComment
	e := g.Db.Model(&model.ArticleComment{ArticleId: articleId}).Find(&articleCommentList).Error
	if e != nil {
		return []model.ArticleComment{}, e
	} else {
		return articleCommentList, e
	}
}

// UpdateOrCreateTagReadAmount 插入或者更新interestedTag字段的数据
func UpdateOrCreateTagReadAmount(userId int, tag string) error {
	result := g.Db.Model(&model.InterestedTag{}).Where("user_id = ? AND tag = ?", userId, tag).Updates(map[string]interface{}{
		"read_amount": gorm.Expr("read_amount+?", 1),
	})
	fmt.Println(result.RowsAffected)
	if result.RowsAffected == 0 {
		fmt.Println("OK")
		record := model.InterestedTag{
			UserId:     userId,
			Tag:        tag,
			ReadAmount: 1,
		}
		if e := g.Db.Model(&model.InterestedTag{}).Create(&record).Error; e != nil {
			return e
		}
	}
	return nil
}

// GetArticleTagAmount 获取某人的喜欢标签的频次
func GetArticleTagAmount(userId int) map[string]int {
	var tagList []model.InterestedTag
	g.Db.Model(&model.InterestedTag{}).Where("user_id = ?", userId).Find(&tagList)
	var tagAmountMap map[string]int
	tagAmountMap = make(map[string]int)
	for _, v := range tagList {
		tagAmountMap[v.Tag] = v.ReadAmount
	}
	//fmt.Println("tagAmountMap:", tagAmountMap)
	return tagAmountMap
}

// InsertReadArticleId 插入已被user读的文章ID
func InsertReadArticleId(userId int, articleId int) {
	// 若数据已经存在，那么直接更新
	createdTime := time.Now()
	expirationTime := time.Now().Add(72 * time.Hour)

	result := g.Db.Model(&model.PassReaderRecord{}).Where("user_id = ? and article_id = ?", userId, articleId).Updates(map[string]interface{}{
		"created_time":    createdTime,
		"expiration_time": expirationTime,
	})

	if result.RowsAffected == 0 {
		record := model.PassReaderRecord{
			UserId:         userId,
			ArticleId:      articleId,
			CreatedTime:    createdTime,
			ExpirationTime: expirationTime,
		}
		g.Db.Model(&model.PassReaderRecord{}).Create(&record)
	}

}

// JudgeIsExpiration 判断被被插入文章是否过期,对不存在记录作过期处理
func JudgeIsExpiration(userId int, articleId int) (bool, error) {
	var record model.PassReaderRecord
	result := g.Db.Model(&model.PassReaderRecord{}).Where("user_id = ? and article_id = ?", userId, articleId).First(&record)

	if result.Error != nil {
		return true, nil
	} else {
		if record.ExpirationTime.Before(time.Now()) {
			// 在当前时间之前，还没有过期
			return false, nil
		} else {
			return true, nil
		}
	}
}

// DeleteExpirationData 删除所有过期的数据
func DeleteExpirationData() {
	g.Db.Model(&model.PassReaderRecord{}).Where("expiration_time<?", time.Now()).Delete(&model.PassReaderRecord{})
}

// GetAllReadArticleId 获取所有的已读且未过期的文章
func GetAllReadArticleId(userId int) map[int]bool {
	DeleteExpirationData()
	var recordList []model.PassReaderRecord
	var recordMap map[int]bool
	recordMap = make(map[int]bool)
	g.Db.Model(&model.PassReaderRecord{UserId: userId}).Where("user_id = ?", userId).Find(&recordList)
	for _, v := range recordList {
		recordMap[v.ArticleId] = true
	}
	return recordMap
}

// RandomGetArticleByTag 根据tag值随机获取一篇文章ID值
func RandomGetArticleByTag(tag string) int {
	var articleId int
	// 随机获取一条tag包含tag的数据
	e := g.Db.Model(&model.OtherArticleInfo{}).Where("tag like ?", "%"+tag+"%").Order("RAND()").Select("article_id").First(&articleId).Error
	if e != nil {
		return 0
	} else {
		return articleId
	}
}
