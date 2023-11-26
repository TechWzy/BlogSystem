package model

import "time"

// 用户模型

type User struct {
	ID            int    `gorm:"primaryKey" json:"id" form:"id"`
	Name          string `json:"name" form:"name" binding:"required"`
	Account       string `json:"account" form:"account"`
	Password      string `json:"password" form:"password" binding:"required"`
	Sex           string `json:"sex" form:"sex" binding:"required,oneof=man woman secret" `
	Age           int    `json:"age" form:"age" binding:"required,min=1,max=120"`
	Introduction  string `json:"introduction" form:"introduction" gorm:"type:varchar(200);default:'此人很懒，什么也没有填写.....';comment:'自我介绍'"`
	FollowMessage string `json:"followMessage" form:"followMessage" gorm:"type:varchar(200);default:'你好，朋友！很荣幸得到你的认可，希望在往后的日子中与你一起交流进步！';comment:'粉丝关注留言'"`
}

// AuthInfo 与Jwt搭配使用的结构体
type AuthInfo struct {
	ID           int
	Name         string
	Sex          string
	Age          string
	Introduction string
}

// Article 文章模型
type Article struct {
	ArticleId        int              `gorm:"primaryKey" json:"articleId" form:"articleId"`
	AuthorId         int              `json:"authorId" form:"authorId"`
	AuthorName       string           `json:"authorName" form:"authorName"`
	CreatedTime      string           `json:"createdTime" form:"createdTime"`
	Title            string           `json:"title" form:"title"`
	Content          string           `json:"content" form:"content" binding:"required"`
	OtherArticleInfo OtherArticleInfo `gorm:"foreignKey:ArticleId"`
}

type OtherArticleInfo struct {
	ArticleId  int
	ReadAmount int
	Like       int
	Dislike    int
	Class      string `json:"class" form:"class" binding:"required"`
	Tag        string `json:"tag" form:"tag" binding:"required"`
}

// UserResponse1 记录用户对文章的响应(Redis,其中UA1:userId:articleId将作为键)
type UserResponse1 struct {
	UserId    int `json:"userId" form:"userId" `
	ArticleId int `json:"articleId" form:"articleId" binding:"required"`
	// 1表示复位，2表示赞，3表示踩，4表示无任何绑定
	LikeCondition int `json:"likeCondition" form:"likeCondition"`
	// 1表示收藏，2表示无收藏,3表示无绑定
	CCondition int `json:"CCondition" form:"CCondition"`
}

// Top10Article 获取top10文章，该结构体存储需要展示的信息
type Top10Article struct {
	Rank       int
	Title      string
	ArticleId  int
	AuthorId   int
	AuthorName string
	ReadAmount int
}

// Top10AmountRecord 获取阅读量排名前十的记录
type Top10AmountRecord struct {
	ArticleId  int
	ReadAmount int
}

// Top5Article 获取某标签的top5文章
type Top5Article struct {
	Rank       int
	Title      string
	ArticleId  int
	AuthorId   int
	AuthorName string
	ReadAmount int
}

// BaseMessage 与消息相关的结构体
type BaseMessage struct {
	ID          int // 消息ID号，新键更新旧键
	SId         int `json:"sId" form:"sId"`
	RId         int `json:"rId" form:"rId" binding:"required"`
	Expiration  int `json:"exp" form:"exp"`
	CreatedTime string
}

// PrivateLetter 私信结构体
type PrivateLetter struct {
	BaseMessage BaseMessage
	Subject     string `json:"subject" form:"subject" binding:"required"`
	Content     string `json:"content" form:"content" binding:"required"`
}

// 粉丝关注模块

// SetFollowCondition 设置粉丝状态
type SetFollowCondition struct {
	FollowId int
	IdolId   int  `json:"idolId" form:"idolId" binding:"required"`
	Cond     bool `json:"cond" form:"cond"`
}

// ArticleComment 评论结构体
type ArticleComment struct {
	ID        int    `gorm:"primaryKey" json:"id" form:"id"`
	UserId    int    `json:"userId" form:"userId"`
	ArticleId int    `json:"articleId" form:"articleId" binding:"required"`
	Comment   string `json:"comment" form:"comment" binding:"required"`
}

// InterestedTag 记录用户某标签文章的阅读次数
type InterestedTag struct {
	UserId     int
	Tag        string
	ReadAmount int
}

// PassReaderRecord 定义user已读文章的id号
type PassReaderRecord struct {
	UserId         int
	ArticleId      int
	CreatedTime    time.Time
	ExpirationTime time.Time
}
