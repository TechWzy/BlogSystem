package controller

import (
	"Project/BlogSystem/internal/dao/mysqls"
	"Project/BlogSystem/internal/dao/rediss"
	"Project/BlogSystem/internal/logic"
	"Project/BlogSystem/internal/model"
	"Project/BlogSystem/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

// SearchArticleByTitle 按照文章标题，搜索文章，返回值是一个文章列表
func SearchArticleByTitle(c *gin.Context) {
	var title model.ArticleTitle
	if err := c.ShouldBind(&title); err != nil {
		c.JSON(200, gin.H{
			"msg": "标题绑定失败...",
		})
		c.Abort()
		return
	}
	v, err := mysqls.GetArticleListByTitle(&title)
	if err != nil {
		c.JSON(200, gin.H{
			"msg": "文章列表获取失败...",
		})
		c.Abort()
		return
	}
	c.Set("articleList", v)
	c.Set("mode", title.Mode)
	c.Next()
}

// SearchByTag 按照标签查找文章，文章整合为一个列表
func SearchByTag(c *gin.Context) {
	var tag model.Tag
	if err := c.ShouldBind(&tag); err != nil {
		c.JSON(200, gin.H{
			"msg": "绑定tag失败...",
		})
		c.Abort()
		return
	}

	// 获取ArticleList
	articleList := rediss.SearchArticleByTag(tag.Tag)
	response := &model.TagResponse{
		Tag:         tag.Tag,
		ArticleList: articleList,
	}
	c.JSON(200, gin.H{
		"Response": response,
	})
}

// ReadArticleById 根据文章Id值阅读一篇文章
func ReadArticleById(c *gin.Context) {
	var articleId model.ArticleId
	if err := c.ShouldBind(&articleId); err != nil {
		c.JSON(200, gin.H{
			"msg": "绑定失败...",
		})
		c.Abort()
		return
	}
	// 获取otherInfo,此数据包含阅读量信息
	otherInfo, err := rediss.GetArticleInfo(articleId.Id)
	if err != nil {
		c.JSON(200, gin.H{
			"Id":  articleId.Id,
			"msg": "该文章不存在...",
		})
		c.Abort()
		return
	}
	if err := service.UpdateArticleReadAmount(otherInfo, 1); err != nil {
		c.JSON(200, gin.H{
			"msg": "阅读量更新失败...",
		})
		c.Abort()
		return
	}
	var article *model.Article
	article, err = mysqls.GetArticleById(otherInfo.ArticleId)
	if err != nil {
		c.JSON(200, gin.H{
			"msg": "文章展示失败...",
		})
		c.Abort()
		return
	}
	// 此处的查询本没有必要，主要是为了验证redis是否更新成功
	var articleExhibition *model.ArticleExhibition
	articleExhibition, err = logic.GetArticleExhibition(article)
	if err != nil {
		c.JSON(200, gin.H{
			"msg": "展示失败...",
		})
		c.Abort()
		return
	}
	// 展示文章
	c.JSON(200, gin.H{
		"ArticleExhibition": articleExhibition,
	})

	// 将已阅读文章插入数据库
	v, _ := c.Get("authInfo")
	authInfo := v.(model.AuthInfo)
	service.UpdateReadArticleAndTag(authInfo.ID, articleId.Id, otherInfo.Tag)
}

// UserArticleInteraction1 绑定用户文章操作数据（包括:Like,Dislike,IsCollected）
func UserArticleInteraction1(c *gin.Context) {
	// 解析token字符串
	v, ok := c.Get("authInfo")
	if !ok {
		c.JSON(200, gin.H{
			"msg": "无法解析用户token",
		})
		c.Abort()
		return
	}
	authInfo := v.(model.AuthInfo)

	// 绑定数据
	var userResponse model.UserResponse1
	userResponse.UserId = authInfo.ID
	userResponse.LikeCondition = 4
	userResponse.CCondition = 3
	if err := c.ShouldBind(&userResponse); err != nil {
		c.JSON(200, gin.H{
			"msg": "数据绑定失败...",
		})
		c.Abort()
		return
	}

	if flag := rediss.RedisIsArticleExist(userResponse.ArticleId); !flag {
		c.JSON(200, gin.H{
			"msg": "交互的文章不存在...",
		})
		c.Abort()
		return
	}
	// 更新用户文章交互收藏信息
	if err := service.UpdateUACollectionInfo(&userResponse); err != nil {
		c.JSON(200, gin.H{
			"msg": "文章收藏信息更新失败...",
		})
		c.Abort()
		return
	}
	// 判断交互信息是否已经存在于redis中
	val, err := rediss.GetUserResponse1(userResponse.UserId, userResponse.ArticleId)
	if err != nil {
		if userResponse.LikeCondition == 4 {
			userResponse.LikeCondition = 1
		}
		if err := rediss.InsertUserResponse1(&userResponse); err != nil {
			c.JSON(200, gin.H{
				"msg": "插入数据失败...",
			})
			c.Abort()
			return
		}
		// 更新Mysql中的other_article)infos
		if err := service.MUpdateUA11(userResponse.ArticleId, userResponse.LikeCondition); err != nil {
			c.JSON(200, gin.H{
				"msg": "数据更新失败...",
			})
			c.Abort()
			return
		}
	} else {
		// 更新redis与mysql中的字段信息
		if userResponse.LikeCondition == 4 {
			c.JSON(200, gin.H{
				"msg": "无需更新",
			})
			return
		}
		// 插入redis数据库
		if err := rediss.InsertUserResponse1(&userResponse); err != nil {
			c.JSON(200, gin.H{
				"msg": "插入数据失败...",
			})
			c.Abort()
			return
		}

		// 断言value值
		ua2 := val.(*model.UserResponse1)
		if err := service.MUpdateUA12(userResponse.ArticleId, userResponse.LikeCondition, ua2.LikeCondition); err != nil {
			c.JSON(200, gin.H{
				"msg": "Mysql数据库更新失败...",
			})
			c.Abort()
			return
		}
	}

	// 更新redis
	if err := service.RUpdateUA12(userResponse.ArticleId); err != nil {
		c.JSON(200, gin.H{
			"msg": "Redis数据库更新失败...",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"msg": "更新成功...",
	})
}

// SetComment 某人对某篇文章进行评论
func SetComment(c *gin.Context) {
	v, _ := c.Get("authInfo")
	authInfo := v.(model.AuthInfo)

	var articleComment model.ArticleComment
	if e := c.ShouldBind(&articleComment); e != nil {
		c.JSON(200, gin.H{
			"Msg": "无效绑定...",
		})
		c.Abort()
		return
	}

	// 插入数据库
	articleComment.UserId = authInfo.ID
	if e := mysqls.InsertArticleComment(&articleComment); e != nil {
		c.JSON(200, gin.H{
			"Msg": "数据插入失败...",
		})
		c.Abort()
		return
	}

	// 数据插入成功
	c.JSON(200, gin.H{
		"Msg": "数据插入成功...",
	})
}

// GetArticleComment 获取文章的所有评论
func GetArticleComment(c *gin.Context) {
	var articleId model.ArticleId
	fmt.Println("OK")
	if e := c.ShouldBind(&articleId); e != nil {
		c.JSON(200, gin.H{
			"Msg": "无法绑定文章ID",
		})
		c.Abort()
		return
	}
	articleCommentList, e := service.GetArticleCommentResponse(articleId.Id)
	if e != nil {
		c.JSON(200, gin.H{
			"Msg": "无法获取文章评论响应信息",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"articleComment": articleCommentList,
	})
}

// RecommendFiveArticleByTag 为用户推荐i篇文章
func RecommendFiveArticleByTag(c *gin.Context) {
	v, _ := c.Get("authInfo")
	authInfo := v.(model.AuthInfo)

	// 获取推荐文章
	list, e := service.GetRecommendArticleList(authInfo.ID)
	if e != nil {
		c.JSON(200, gin.H{
			"Msg": "出现了Wrong articleId",
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"response": *list,
	})
}
