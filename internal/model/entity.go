package model

import "github.com/golang-jwt/jwt/v4"

// KeyValue 解析KeyValue字符串为结构体
type KeyValue struct {
	Key   string
	Value string
}

type IDAccount struct {
	ID      int64
	Account string
}

// AuthInfoClaims 生成登录成功后的token
type AuthInfoClaims struct {
	AuthInfo AuthInfo
	jwt.RegisteredClaims
}

// UploadResponse UploadArticle后的响应报文
type UploadResponse struct {
	Code        int
	ArticleId   int
	AuthorId    int
	AuthorName  string
	CreatedTime string
	Title       string
	Content     string
	ReadAmount  int
	Like        int
	Dislike     int
	Class       string
	Tag         string
	Msg         string
}

// SimpleArticleResponse 获取简洁版的搜索结果，仅仅包括文章标题，用户ID与文章ID
type SimpleArticleResponse struct {
	Title     string
	ArticleId int
	AuthorId  int
}

// TagList 标签查找的响应报文结构体
type TagList struct {
	Title     string
	ArticleId int
	AuthorId  int
}

// TagResponse Tag获取article的响应报文
type TagResponse struct {
	Tag         string
	ArticleList []TagList
	Msg         string
}

// ArticleExhibition 展示文章的详细信息
type ArticleExhibition struct {
	Title            string
	Content          string
	AuthorName       string
	CreatedTime      string
	ReadAmount       int
	CollectionAmount int
	Like             int
	Dislike          int
	Tag              string
	ArticleId        int
	AuthorId         int
}

// Top10ArticleExhibition 展示阅读量top10的文章
type Top10ArticleExhibition struct {
	ArticleId  int
	AuthorId   int
	AuthName   string
	Title      string
	ReadAmount int
}

// PrivateLetterField 启动getPrivateLetter的响应报文
type PrivateLetterField struct {
	Subject     string
	Content     string
	SenderId    int
	Sender      string
	CreatedTime string
}

// PrivateLetterResponse 获取私信的报文结构体
type PrivateLetterResponse struct {
	Type              string
	PrivateLetterList []PrivateLetterField
}

// AComment 文章评论结构体
type AComment struct {
	CommentatorName string
	CommentatorId   int
	Comment         string
}

// ArticleCommentList 文章评论展示列表
type ArticleCommentList struct {
	Title       string
	CommentList []AComment
}

// RecommendArticle 推荐文章展示
type RecommendArticle struct {
	Title     string
	ArticleId int
	AuthId    int
	Tag       string
}

// RecommendArticleList 推荐文章展示列表
type RecommendArticleList struct {
	Msg  string
	List []RecommendArticle
}
