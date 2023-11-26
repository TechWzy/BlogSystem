package api

import (
	"Project/BlogSystem/internal/controller"
	"Project/BlogSystem/internal/logic"
	"github.com/gin-gonic/gin"
	"log"
)

func InitRouter() {
	r := gin.Default()

	//r.GET("/test", func(c *gin.Context) {
	//	mysqls.Test("云原生")
	//})

	user := r.Group("/user")
	{
		// 1.处理用户注册行为
		user.POST("/register", controller.CreateInfo, logic.FillInfo())
		// 2.处理用户登录行为，所生成的token的解析字段应该包含user的大部分信息
		user.POST("/login", controller.UserLogin, logic.GetAuthTokenMiddleware())
	}

	article := r.Group("/art")
	{
		// 1.上传文章
		article.POST("/upload", logic.ParseAuthInfoMiddleware(), controller.UploadArticle)
		// 2.获取阅读量排名前十的文章
		article.GET("/getTopTenArticle", controller.SelectTop10ArticleByReadAmount)
		// 3.根据文章标签，获取排名前五的文章
		article.POST("/getTopFive", controller.GetTop5ArticleByTag)
	}

	UaIt := r.Group("/UaIt")
	{
		// 1.按照文章标题搜索文章列表
		UaIt.POST("/searchByTitle", controller.SearchArticleByTitle, logic.OutputArticleByMode())
		// 2.按照文章标签搜索文章列表
		UaIt.POST("/searchByTag", controller.SearchByTag)
		// 3.根据Id值查询文章
		UaIt.POST("/readArticleById", logic.ParseAuthInfoMiddleware(), controller.ReadArticleById)
		// 4.用户文章交互操作，负责更新(like,dislike,isCollected)
		UaIt.POST("/UaIt1", logic.ParseAuthInfoMiddleware(), controller.UserArticleInteraction1)
		// 5.用户设置评论
		UaIt.POST("/setComment", logic.ParseAuthInfoMiddleware(), controller.SetComment)
		// 6.获取文章的所有评论
		UaIt.POST("/GetArticleComment", controller.GetArticleComment)
		// 7.获取5篇推荐文章
		UaIt.POST("/recommendFiveArticleByTag", logic.ParseAuthInfoMiddleware(), controller.RecommendFiveArticleByTag)
	}

	UUIt := r.Group("/UUIt")
	{
		// 1.设置粉丝状态
		UUIt.POST("/setFollowCond", logic.ParseAuthInfoMiddleware(), controller.SetFollowCond)
	}

	message := r.Group("/message")
	{
		// 1.发送一条私信
		message.POST("/privateLetter", logic.ParseAuthInfoMiddleware(), controller.SendPrivateLetter)
		// 2.获取私信
		message.POST("/getPrivateLetter", logic.ParseAuthInfoMiddleware(), controller.GetPLetter)
	}

	// 启动监听
	if err := r.Run(":8080"); err != nil {
		// 直接退出程序，结束各种连接
		log.Println("Router doesn't initialize successfully!")
		log.Fatalln(err)
		return
	}
	log.Println("Router initialize successfully!")
}
