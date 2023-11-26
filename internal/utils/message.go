package utils

import "fmt"

// CombinePrivateLetterField 组合私信查询字段
func CombinePrivateLetterField(sId int, rId int) string {
	patternString := fmt.Sprintf("%s#%v#%v%s", "RSPrivate", sId, rId, "*")
	return patternString
}

// CombineAllPrivateLetterField 组合所有私信字段
func CombineAllPrivateLetterField(rId int) string {
	patternString := fmt.Sprintf("%s#%s#%v%s", "RSPrivate", "*", rId, "*")
	return patternString
}
