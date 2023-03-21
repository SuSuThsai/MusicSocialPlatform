package utils

const (
	SUCCESS = 200
	ERROR   = 500

	//code=1000...用户模块错误
	ErrorUsernameUsed            = 1001
	ErrorPasswordWrong           = 1002
	ErrorUserNotExist            = 1003
	ErrorUserNotRight            = 1008
	ErrorTokenExist              = 1004
	ErrorTokenRuntime            = 1005
	ErrorTokenWrong              = 1006
	ErrorTokenTypeWrong          = 1007
	ErrorUserCreatFail           = 1009
	ErrorPicturePFPNotFix        = 1010
	ErrorPictureBackgroundNotFix = 1011
	ErrorPicturePFPTooBig        = 1012
	ErrorPictureBackgroundTooBig = 1013
	TargetNotExit                = 1014
	UserInfoConstrict            = 1015

	//code=4000 音乐模块
	ErrorMusicLNotFix    = 4001
	ErrorMusicWNotFix    = 4002
	ErrorMusicLNotTooBig = 4003
	ErrorMusicWNotTooBig = 4004
	//code=2000...分类模块错误
	ErrorCategorynameUsed = 2001
	ErrorCategoryNotExist = 2003

	//code=3000...文章模块错误
	ErrorArticleNotExist   = 3003
	ErrorMusicListNotExist = 3004
)

var codeMsg = map[int]string{
	SUCCESS:                      "OK",
	ERROR:                        "FAIL",
	ErrorUsernameUsed:            "用户名已存在",
	ErrorPasswordWrong:           "密码错误",
	ErrorUserNotExist:            "用户不存在",
	ErrorUserNotRight:            "用户权限不足",
	ErrorUserCreatFail:           "创建用户失败",
	ErrorPicturePFPNotFix:        "上传失败，头像非jpg,png,jpeg,gif,WebP,请重新上传！",
	ErrorPictureBackgroundNotFix: "上传失败，背景非jpg,png,jpeg,gif,WebP,请重新上传！",
	ErrorPicturePFPTooBig:        "上传失败，头像不能超过4M！",
	ErrorPictureBackgroundTooBig: "上传失败，背景不能超过8M！",
	TargetNotExit:                "对象不存在",
	UserInfoConstrict:            "用户信息冲突",

	ErrorTokenExist:     "TOKEN不存在",
	ErrorTokenRuntime:   "TOKEN过期",
	ErrorTokenWrong:     "TOKEN不正确",
	ErrorTokenTypeWrong: "TOKEN格式错误",

	ErrorMusicLNotFix:    "上传失败,音乐非mp3,flac,aac,amr,ape,wma,mpeg-4请重新上传！",
	ErrorMusicWNotFix:    "上传失败,歌词非lrc,txt,doc,docx！",
	ErrorMusicLNotTooBig: "上传失败，音乐不能超过30M！",
	ErrorMusicWNotTooBig: "上传失败,音乐不能超过2M！",

	ErrorCategorynameUsed: "分类已存在",
	ErrorCategoryNotExist: "分类不存在",

	ErrorArticleNotExist:   "文章不存在",
	ErrorMusicListNotExist: "歌单不存在",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
