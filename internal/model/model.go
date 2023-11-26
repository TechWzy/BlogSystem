package model

import "github.com/golang-jwt/jwt/v4"

// LoginInfo 处理登录业务的数据结构(account,password)
type LoginInfo struct {
	Account  string `json:"account" form:"account" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// UploadArticle 绑定上传文章所需要的字段
type UploadArticle struct {
	Title   string `json:"title" form:"title" binding:"required"`
	Content string `json:"content" form:"content" binding:"required"`
	Class   string `json:"class" form:"class" binding:"required"`
	Tag     string `json:"tag" form:"tag" binding:"required"`
}

// MyClaims MyClaim存放jwt加密信息
type MyClaims struct {
	ID      int64  `json:"id" form:"id"`
	Account string `json:"account" form:"account"`
	jwt.RegisteredClaims
}

// ArticleTitle 文章标题结构体
type ArticleTitle struct {
	Title string `json:"title" form:"title" binding:"required"`
	Mode  int    `json:"mode" form:"mode" binding:"required,oneof=0 1 2"` // response的输出模式
}

// Tag 按照标签查找文章
type Tag struct {
	Tag string `json:"tag" form:"tag" binding:"required"`
}

// ArticleId 文章的Id
type ArticleId struct {
	Id int `json:"id" form:"id" binding:"required"`
}

// GetPrivateLetter 获取所有私信
type GetPrivateLetter struct {
	RID  int `json:"rid" form:"rid"`
	SID  int `json:"sid" form:"sid" binding:"required"`
	Exp  int `json:"exp" form:"exp"`
	Mode int `json:"mode" form:"mode" binding:"required,oneof= 1 2"`
}
