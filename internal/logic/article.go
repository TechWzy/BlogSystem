package logic

import (
	"Project/BlogSystem/internal/dao/rediss"
	"Project/BlogSystem/internal/model"
	"time"
)

// FillArticleInfo 补充文章的信息
func FillArticleInfo(authInfo *model.AuthInfo, articleInfo *model.UploadArticle) *model.Article {
	var article model.Article
	article.AuthorId = authInfo.ID
	article.AuthorName = authInfo.Name
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04")
	article.CreatedTime = formattedTime
	article.Title = articleInfo.Title
	article.Content = articleInfo.Content
	article.OtherArticleInfo.Tag = articleInfo.Tag
	article.OtherArticleInfo.Like = 0
	article.OtherArticleInfo.Dislike = 0
	article.OtherArticleInfo.ReadAmount = 0
	article.OtherArticleInfo.Class = articleInfo.Class
	return &article
}

// CreateResponse 返回报文
func CreateResponse(article *model.Article, otherInfo *model.OtherArticleInfo) *model.UploadResponse {
	response := &model.UploadResponse{
		Code:        200,
		ArticleId:   article.ArticleId,
		AuthorId:    article.AuthorId,
		AuthorName:  article.AuthorName,
		CreatedTime: article.CreatedTime,
		Title:       article.Title,
		Content:     article.Content,
		ReadAmount:  otherInfo.ReadAmount,
		Like:        otherInfo.Like,
		Dislike:     otherInfo.Dislike,
		Class:       otherInfo.Class,
		Tag:         otherInfo.Tag,
		Msg:         "Success",
	}
	return response
}

// GetArticleExhibition 获取一个articleExhibition
func GetArticleExhibition(article *model.Article) (*model.ArticleExhibition, error) {
	otherArticleInfo, err := rediss.GetArticleInfo(article.ArticleId)
	if err != nil {
		return &model.ArticleExhibition{}, err
	}
	var amount int
	// 无论真假
	amount = rediss.GetACollectionInfo(otherArticleInfo.ArticleId)
	return &model.ArticleExhibition{
		Title:            article.Title,
		Content:          article.Content,
		AuthorName:       article.AuthorName,
		CreatedTime:      article.CreatedTime,
		ReadAmount:       otherArticleInfo.ReadAmount,
		CollectionAmount: amount,
		Like:             otherArticleInfo.Like,
		Dislike:          otherArticleInfo.Dislike,
		Tag:              otherArticleInfo.Tag,
		ArticleId:        article.ArticleId,
		AuthorId:         article.AuthorId,
	}, nil
}
