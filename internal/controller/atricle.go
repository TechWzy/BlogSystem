package controller

import (
	"Project/BlogSystem/internal/dao/mysqls"
	"Project/BlogSystem/internal/dao/rediss"
	"Project/BlogSystem/internal/logic"
	"Project/BlogSystem/internal/model"
	"Project/BlogSystem/internal/service"
	"github.com/gin-gonic/gin"
)

// UploadArticle 上传文章
func UploadArticle(c *gin.Context) {
	// 简化原则，只需绑定required字段
	var articleInfo model.UploadArticle
	if err := c.ShouldBind(&articleInfo); err != nil {
		c.JSON(200, gin.H{
			"msg": "Binding operation fails!",
			"err": err,
		})
		c.Abort()
		return
	}

	// authInfo来自tokenString解析而来的结构体
	v, ok := c.Get("authInfo")
	if !ok {
		c.JSON(200, gin.H{
			"msg": "token过期",
		})
	}

	authInfo := v.(model.AuthInfo)
	var article *model.Article

	// 补充信息
	article = logic.FillArticleInfo(&authInfo, &articleInfo)
	// Article插入Mysql数据库
	otherInfo := mysqls.InsertArticleInfo(article)

	// otherInfo插入redis数据库
	if err := rediss.InsertArticleInfo(otherInfo); err != nil {
		c.JSON(200, gin.H{
			"msg": "otherInfo无法转化为redis数据",
		})
		c.Abort()
		return
	}

	// Tag类型数据插入redis数据库
	if err := rediss.RedisInsertTag(otherInfo.Tag, article.Title, article.ArticleId, article.AuthorId); err != nil {
		c.JSON(200, gin.H{
			"msg": "标签插入失败...",
		})
		c.Abort()
		return
	}
	uploadResponse := logic.CreateResponse(article, otherInfo)

	// 向每一位粉丝推荐自己的文章
	service.RecommendArticle(authInfo.ID, &articleInfo)
	c.JSON(200, gin.H{
		"Response": uploadResponse,
	})

}

// SelectTop10ArticleByReadAmount 根据文章阅读量选择Top10的文章
func SelectTop10ArticleByReadAmount(c *gin.Context) {
	var articleList []model.Top10Article
	var err error
	articleList, err = service.GetTopTenAmountArticle()
	if err != nil {
		c.JSON(200, gin.H{
			"msg": "查阅失败...",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"articleList": articleList,
	})
}

func GetTop5ArticleByTag(c *gin.Context) {
	var tag model.Tag
	if err := c.ShouldBind(&tag); err != nil {
		c.JSON(200, gin.H{
			"msg": "Tag绑定失败...",
		})
		c.Abort()
		return
	}
	articleList := service.GetTopFiveArticle(tag.Tag)
	c.JSON(200, gin.H{
		"articleList": articleList,
	})
}
