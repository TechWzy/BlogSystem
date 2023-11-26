package logic

import (
	"Project/BlogSystem/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

// FillInfo 填充未绑定信息
func FillInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		rand.Seed(time.Now().UnixNano()) // 设置随机种子

		account := getAccount("AccountPassword")
		c.Set("account", account)
		c.Next()
	}
}

// OutputArticleByMode 根据不同输出模式，输出不同格式的报文
func OutputArticleByMode() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("articleList")
		if !ok {
			c.JSON(200, gin.H{
				"msg": "无法获取articleList",
			})
			c.Abort()
			return
		}
		articleList := v.([]model.Article)
		v, ok = c.Get("mode")
		if !ok {
			c.JSON(200, gin.H{
				"msg": "无法获取mode",
			})
		}
		mode := v.(int)
		// 全套输出
		if mode == 2 {
			articleExhibitionList := make([]model.ArticleExhibition, len(articleList))
			// 循环遍历补充article的其余信息
			for num, article := range articleList {
				articleExhibition, err := GetArticleExhibition(&article)
				if err != nil {
					fmt.Printf("ArticleId:%v,获取失败...", article.ArticleId)
				} else {
					articleExhibitionList[num] = *articleExhibition
				}
			}

			c.JSON(200, gin.H{
				"msg":                   "文章列表获取成功...",
				"articleExhibitionList": articleExhibitionList,
			})
		} else if mode == 1 {
			var simpleArticleList []model.SimpleArticleResponse
			for _, article := range articleList {
				simpleArticleRecord := model.SimpleArticleResponse{
					Title:     article.Title,
					ArticleId: article.ArticleId,
					AuthorId:  article.AuthorId,
				}
				simpleArticleList = append(simpleArticleList, simpleArticleRecord)
			}
			c.JSON(200, gin.H{
				"msg":               "获取成功...",
				"simpleArticleList": simpleArticleList,
			})
		}
	}
}
