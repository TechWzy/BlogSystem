package mysqls

import (
	g "Project/BlogSystem/internal/global"
	"Project/BlogSystem/internal/model"
)

// InsertInfo 创建的用户数据插入数据库
func InsertInfo(user *model.User) {
	g.Db.Create(&model.User{Name: user.Name, Account: user.Account, Password: user.Password, Sex: user.Sex, Age: user.Age, Introduction: user.Introduction, FollowMessage: user.FollowMessage})
	//g.Db.Select("Name", "Account", "Password", "Sex", "Age", "Introduction", "FollowMessage").Create(user)
	//fmt.Println("Insert successfully!")
}

// GetAuthInfo 获取用户登录后存储在authInfo上的数据
func GetAuthInfo(account string) model.AuthInfo {
	var authInfo model.AuthInfo
	g.Db.Model(&model.User{}).Where("account = ?", account).Find(&authInfo)
	return authInfo
}

// JudgeUserIsExist 判断是否存在article_Id列是否存在特定值
func JudgeUserIsExist(userId int) error {
	var user model.User
	if err := g.Db.Model(&model.User{}).Where("id = ?", userId).First(&user).Error; err != nil {
		return err
	}
	return nil
}

// GetSenderNameById 获取senderName
func GetSenderNameById(senderId int) (string, error) {
	var name string
	if err := g.Db.Model(&model.User{}).Where("id = ?", senderId).Select("name").First(&name).Error; err != nil {
		return "", err
	} else {
		//fmt.Println("name:", name)
		return name, nil
	}
}

// GetFollowMessage 根据id值获取用户的followMessage
func GetFollowMessage(idolId int) (string, error) {
	var followMessage string
	if err := g.Db.Model(&model.User{}).Where("id = ?", idolId).Select("FollowMessage").First(&followMessage).Error; err != nil {
		return "无法获取followMessage", err
	} else {
		return followMessage, err
	}
}
