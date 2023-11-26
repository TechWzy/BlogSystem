package service

import (
	"Project/BlogSystem/internal/dao/mysqls"
	"Project/BlogSystem/internal/dao/rediss"
	"Project/BlogSystem/internal/model"
	"fmt"
)

// SortArticleByAmount 对文章列表进行排序
func SortArticleByAmount(articleList []model.Top5Article) []model.Top5Article {
	sz := len(articleList)
	for i := 0; i < sz-1; i++ {
		for j := 0; j < sz-i-1; j++ {
			if articleList[j].ReadAmount < articleList[j+1].ReadAmount {
				article := articleList[j]
				articleList[j] = articleList[j+1]
				articleList[j+1] = article
			}
		}
	}
	if sz >= 5 {
		return articleList[0:5]
	} else {
		return articleList
	}

}

// GetTopTenAmountArticle 获取阅读量排名前十的文章
func GetTopTenAmountArticle() ([]model.Top10Article, error) {
	var articleList []model.Top10Article
	articleIdList, err := mysqls.GetTopTenArticleId()
	if err != nil {
		return articleList, err
	}
	for key, val := range articleIdList {
		fmt.Println(key, val)
		v, e := mysqls.GetTop10ArticleById(val.ArticleId)
		if e != nil {
			return []model.Top10Article{}, e
		}
		var article model.Top10Article
		article.Rank = key + 1
		article.Title = v.Title
		article.ArticleId = val.ArticleId
		article.ReadAmount = val.ReadAmount
		article.AuthorId = v.AuthorId
		article.AuthorName = v.AuthorName
		articleList = append(articleList, article)
	}
	return articleList, nil
}

func GetTopFiveArticle(tag string) []model.Top5Article {
	articleIdList := rediss.SearchArticleByTag(tag)
	var articleList []model.Top5Article
	for _, v := range articleIdList {
		var article model.Top5Article
		readAmount, _ := mysqls.GetArticleAmountById(v.ArticleId)
		authorName, _ := mysqls.GetAuthorNameById(v.ArticleId)
		fmt.Println("authorName", authorName)
		article.ArticleId = v.ArticleId
		article.AuthorId = v.AuthorId
		article.Title = v.Title
		article.ReadAmount = readAmount
		article.AuthorName = authorName
		articleList = append(articleList, article)
	}
	articleList = SortArticleByAmount(articleList)
	sz := len(articleList)
	for i := 0; i < sz; i++ {
		articleList[i].Rank = i + 1
	}
	return articleList
}
