package mysqls

import (
	g "Project/BlogSystem/internal/global"
	"Project/BlogSystem/internal/model"
	"fmt"
)

// InsertArticleInfo 当成功上传一篇文章后，将数据存储到数据库
func InsertArticleInfo(article *model.Article) *model.OtherArticleInfo {
	// 插入数据
	result := g.Db.Model(model.Article{}).Select("author_id", "author_name", "created_time", "title",
		"content").Create(article)

	if result.Error != nil {
		// 出现错误，插入失败
		fmt.Printf("插入文章失败: %s\n", result.Error)
	} else {
		// 插入成功，result.RowsAffected 表示影响的行数
		fmt.Printf("成功插入文章，影响 %d 行\n", result.RowsAffected)
	}

	// 获取新记录的Id值
	var id int
	g.Db.Model(&model.Article{}).Select("article_id").Where("author_id = ? and created_time = ? and title = ?",
		article.AuthorId, article.CreatedTime, article.Title).First(&id)

	otherInfo := &model.OtherArticleInfo{
		ArticleId:  id,
		ReadAmount: 0,
		Like:       0,
		Dislike:    0,
		Class:      article.OtherArticleInfo.Class,
		Tag:        article.OtherArticleInfo.Tag,
	}
	// 数据插入other_article_infos
	g.Db.Model(&model.OtherArticleInfo{}).Select("article_id", "read_amount", "like",
		"dislike", "class", "tag").Create(otherInfo)
	return otherInfo
}

// GetArticleById 根据文章的ID值获取一篇文章
func GetArticleById(articleId int) (*model.Article, error) {
	var article model.Article
	if err := g.Db.Model(&model.Article{}).Where("article_id = ?", articleId).Find(&article).Error; err != nil {
		return &model.Article{}, err
	}
	return &article, nil
}

// GetArticleListByTitle 根据文章标题搜索相关文章，返回结构体
func GetArticleListByTitle(title *model.ArticleTitle) (any, error) {
	var articleList []model.Article
	keyTitle := "%" + title.Title + "%" // 模糊搜索
	if err := g.Db.Model(&model.Article{}).Where("title LIKE ?", keyTitle).Find(&articleList).Error; err != nil {
		fmt.Println("Wrong")
		return nil, err
	}
	return articleList, nil
}

// IncreaseArticleAmountById 通过文章ID号码更新文章阅读量，阅读量任意
func IncreaseArticleAmountById(articleId int, currentAmount int, increment int) error {
	err := g.Db.Model(&model.OtherArticleInfo{}).Where("article_id = ?", articleId).Update("read_amount", currentAmount+increment).Error
	return err
}

// GetOtherArticleInfos 获取文章的OtherArticleInfos
func GetOtherArticleInfos(articleId int) (*model.OtherArticleInfo, error) {
	var infos model.OtherArticleInfo
	if err := g.Db.Model(&model.OtherArticleInfo{}).Where("article_Id = ?", articleId).First(&infos).Error; err != nil {
		return &model.OtherArticleInfo{}, err
	}
	return &infos, nil
}

// GetTopTenArticleId 获取阅读量排名前十的文章id
func GetTopTenArticleId() ([]model.Top10AmountRecord, error) {
	var articleIdList []model.Top10AmountRecord
	result := g.Db.Model(&model.OtherArticleInfo{}).Order("read_amount desc").Limit(10).Find(&articleIdList)
	if result.Error != nil {
		//fmt.Println("GetTopTenArticleId,wrong")
		return articleIdList, result.Error
	}
	return articleIdList, nil
}

// GetTop10ArticleById 根据文章Id返回文章
func GetTop10ArticleById(articleId int) (*model.Top10Article, error) {
	var article model.Top10Article
	if err := g.Db.Model(&model.Article{}).Where("article_id = ?", articleId).Select("author_id", "author_name", "title").First(&article).Error; err != nil {
		return &model.Top10Article{}, err
	}
	return &article, nil
}

// GetArticleAmountById 根据id至获取文章的阅读量
func GetArticleAmountById(articleId int) (int, error) {
	var readAmount int
	if err := g.Db.Model(&model.OtherArticleInfo{}).Where("article_id = ?", articleId).Select("read_amount").First(&readAmount).Error; err != nil {
		return 0, err
	}
	return readAmount, nil
}

// GetAuthorNameById 根据id值获取作者名
func GetAuthorNameById(articleId int) (string, error) {
	var authorName string
	if err := g.Db.Model(&model.Article{}).Where("article_id = ?", articleId).Select("author_name").First(&authorName).Error; err != nil {
		return "", err
	}
	return authorName, nil
}

func GetArticleTitle(articleId int) string {
	var title string
	// 选择单个字段时一定要用到select
	g.Db.Model(&model.Article{}).Where("article_id = ?", articleId).Select("title").First(&title)
	return title
}

// GetArticleTag 获取文章的标签
func GetArticleTag(articleId int) string {
	var tag string
	g.Db.Model(&model.OtherArticleInfo{}).Where("article_id = ?", articleId).Select("tag").First(&tag)
	return tag
}

// GetArticleAuthId 获取AuthId
func GetArticleAuthId(articleId int) int {
	//fmt.Println("articleId in GetAuthId:", articleId)
	var authId int
	g.Db.Model(&model.Article{}).Where("article_id = ?", articleId).Select("author_id").First(&authId)
	return authId
}
