package service

import (
	"Project/BlogSystem/internal/dao/mysqls"
	"Project/BlogSystem/internal/dao/rediss"
	"Project/BlogSystem/internal/model"
	"Project/BlogSystem/internal/utils"
	"fmt"
	"strings"
)

// service层使用dao层或logic层提供的功能实现服务，注意dao层与logic层不应该使用service层的服务

// UpdateArticleReadAmount 根据ArticleId更新文章的阅读量
func UpdateArticleReadAmount(info *model.OtherArticleInfo, increment int) error {
	if err := mysqls.IncreaseArticleAmountById(info.ArticleId, info.ReadAmount, increment); err != nil {
		return err
	}
	if err := rediss.RedisIncreaseArticleAmountById(info, increment); err != nil {
		return err
	}
	return nil
}

// MUpdateUA11 当用户第一次与文章交互时，更新Mysql中的otherInfo数据
func MUpdateUA11(articleId int, cond int) error {
	var err error
	err = nil
	if cond == 2 {
		err = mysqls.UpdateLikeAndDislike(articleId, 1, 0)
	} else if cond == 3 {
		err = mysqls.UpdateLikeAndDislike(articleId, 0, 1)
	}
	return err
}

// MUpdateUA12 当用户不止第一次与该文章交互
func MUpdateUA12(articleId int, cond1 int, cond2 int) error {
	var err error
	err = nil
	// 当cond1==1表示取消点赞或踩
	if cond1 == 1 {
		if cond2 == 2 {
			err = mysqls.UpdateLikeAndDislike(articleId, -1, 0)
		} else if cond2 == 3 {
			err = mysqls.UpdateLikeAndDislike(articleId, 0, -1)
		}
		return err
	} else if cond1 == 2 {
		if cond2 == 1 {
			err = mysqls.UpdateLikeAndDislike(articleId, 1, 0)
		} else if cond2 == 3 {
			err = mysqls.UpdateLikeAndDislike(articleId, 1, -1)
		}
		return err
	} else if cond1 == 3 {
		if cond2 == 1 {
			err = mysqls.UpdateLikeAndDislike(articleId, 0, 1)
		} else if cond2 == 2 {
			err = mysqls.UpdateLikeAndDislike(articleId, -1, 1)
		}
		return err
	} else {
		return err
	}
}

// RUpdateUA12 更新redis中的otherInfo字段中的like与dislike
func RUpdateUA12(articleId int) error {
	var err error
	err = nil
	// 先获取mysql中的otherInfo
	var info *model.OtherArticleInfo
	info, err = mysqls.GetOtherArticleInfos(articleId)
	if err != nil {
		return err
	}
	err = rediss.InsertArticleInfo(info)
	if err != nil {
		return err
	}
	return err
}

// UpdateUACollectionInfo 更新用户文章的收藏信息
func UpdateUACollectionInfo(response1 *model.UserResponse1) error {
	// 先更新用户文章交互字段，再更新收藏量
	cFlag, err := rediss.GetUACollectionInfo(response1.UserId, response1.ArticleId)
	// 用户文章字段不存在
	//fmt.Println("Err:", err, "CFlag:", cFlag)
	if err != nil {
		if response1.CCondition == 3 {
			response1.CCondition = 2
		}
		if e := rediss.InsertUACollectionInfo(response1.UserId, response1.ArticleId, response1.CCondition); e != nil {
			return e
		}
		if response1.CCondition == 1 {
			if e := rediss.UpdateACollectionInfo(response1.ArticleId, 1); e != nil {
				return e
			}
		}
	} else {
		if response1.CCondition == 3 {
			return nil
		} else if response1.CCondition == 1 && cFlag == 2 {
			if e := rediss.UpdateACollectionInfo(response1.ArticleId, 1); e != nil {
				return e
			}
			if e := rediss.InsertUACollectionInfo(response1.UserId, response1.ArticleId, response1.CCondition); e != nil {
				return e
			}
		} else if response1.CCondition == 2 && cFlag == 1 {
			if e := rediss.UpdateACollectionInfo(response1.ArticleId, -1); e != nil {
				return e
			}
			if e := rediss.InsertUACollectionInfo(response1.UserId, response1.ArticleId, response1.CCondition); e != nil {
				return e
			}
		}
	}
	return nil
}

// GetArticleCommentResponse 获取文章评论列表的响应报文
func GetArticleCommentResponse(articleId int) (model.ArticleCommentList, error) {
	articleCommentList, e := mysqls.GetArticleCommentList(articleId)
	if e != nil || len(articleCommentList) == 0 {
		return model.ArticleCommentList{}, e
	}
	var commentList []model.AComment
	for _, v := range articleCommentList {
		var comment model.AComment
		userName, _ := mysqls.GetSenderNameById(v.UserId)
		comment.CommentatorId = v.UserId
		comment.CommentatorName = userName
		comment.Comment = v.Comment
		commentList = append(commentList, comment)
	}
	var articleComment model.ArticleCommentList
	articleComment.Title = mysqls.GetArticleTitle(articleCommentList[0].ArticleId)
	articleComment.CommentList = commentList
	return articleComment, nil
}

// GetRecommendArticle 获取推荐文章
func GetRecommendArticle(articleId int) *model.RecommendArticle {
	//fmt.Println("articleId:", articleId)
	title := mysqls.GetArticleTitle(articleId)
	tag := mysqls.GetArticleTag(articleId)
	authId := mysqls.GetArticleAuthId(articleId)
	recommendArticle := model.RecommendArticle{
		Title:     title,
		ArticleId: articleId,
		AuthId:    authId,
		Tag:       tag,
	}
	return &recommendArticle
}

// GetRecommendArticleList 获取推荐文章列表
func GetRecommendArticleList(userId int) (*model.RecommendArticleList, error) {
	// 用户文章tag频率
	tagAmountMap := mysqls.GetArticleTagAmount(userId)
	//fmt.Println("tagAmountMap:", tagAmountMap)
	if len(tagAmountMap) == 0 {
		return &model.RecommendArticleList{}, fmt.Errorf("no tag recommend")
	}

	// 获取前缀切片与索引映射
	proList, tagMap := utils.GetArticleTagSlice(tagAmountMap)
	//fmt.Println("proList:", proList, "tagMap:", tagMap)
	// 获取一个随机值,并将该随机值映射为标签
	var tagList []string
	num := 5
	for i := 1; i <= num; i++ {
		v := utils.GetRandomVal(proList)
		tagList = append(tagList, tagMap[v])
	}
	//fmt.Println("tagList:", tagList)
	// 获取所有已读且未过期的文章
	articleMap := mysqls.GetAllReadArticleId(userId)
	//fmt.Println("articleMap:", articleMap)
	// 寻找与tagList相对应的文章
	var articleIdList []int
	count := 50

	selectedArticleId := make(map[int]bool)
	for i := 0; i < num; i++ {
		//fmt.Println(tagList[i])
		id := mysqls.RandomGetArticleByTag(tagList[i]) //根据tag随机获取一个文章Id
		_, exist := articleMap[id]

		if id == 0 {
			return &model.RecommendArticleList{}, fmt.Errorf("wrong articleId")
		}

		// 如果文章已读，那么重新获取id值
		if exist {
			i--
		} else {
			_, exist := selectedArticleId[id]
			if exist {
				i--
			} else {
				selectedArticleId[id] = true
				articleIdList = append(articleIdList, id)
			}
		}
		count--
		if count == 0 {
			break
		}
	}

	//fmt.Println("articleIdList:", articleIdList)

	// 寻找articleIdList里面的文章
	var recommendArticleList model.RecommendArticleList
	var articleList []model.RecommendArticle
	userName, _ := mysqls.GetSenderNameById(userId)
	recommendArticleList.Msg = userName + "文章推荐列表"
	for v := range articleIdList {
		recommendArticle := GetRecommendArticle(articleIdList[v])
		articleList = append(articleList, *recommendArticle)
	}
	recommendArticleList.List = articleList
	return &recommendArticleList, nil
}

// UpdateReadArticleAndTag 更新已读文章与兴趣标签
func UpdateReadArticleAndTag(userId int, articleId int, tag string) {
	mysqls.InsertReadArticleId(userId, articleId)
	parts := strings.Split(tag, ".")
	for _, v := range parts {
		if e := mysqls.UpdateOrCreateTagReadAmount(userId, v); e != nil {
			fmt.Println(e)
		}
	}
}
