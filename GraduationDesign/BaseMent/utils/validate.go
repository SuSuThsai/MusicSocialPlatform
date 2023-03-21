package utils

import (
	"fmt"
	"github.com/go-playground/locales/zh_Hans_CN"
	unTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"mime/multipart"
)

func Validate(data interface{}) (string, int) {
	validate := validator.New()
	uni := unTrans.New(zh_Hans_CN.New())
	trans, _ := uni.GetTranslator("zh_Hans_CN")

	err := zhTrans.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println("err:", err)
	}
	err = validate.Struct(data)
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			return v.Translate(trans), ERROR
		}
	}
	return "", SUCCESS
}

func CheckPicturePFPIsValidate(file *multipart.FileHeader) int {
	if file == nil {
		return TargetNotExit
	}
	fileType := file.Header["Content-Type"][0]
	if fileType != "image/png" && fileType != "image/jpg" && fileType != "image/jpeg" && fileType != "image/gif" && fileType != "image/webp" {
		return ErrorPicturePFPNotFix
	}
	if file.Size/1024 > 4000 {
		return ErrorPicturePFPTooBig
	}
	return SUCCESS
}

func CheckPictureBackgroundIsValidate(file *multipart.FileHeader) int {
	if file == nil {
		return TargetNotExit
	}
	fileType := file.Header["Content-Type"][0]
	if fileType != "image/png" && fileType != "image/jpg" && fileType != "image/jpeg" && fileType != "image/gif" && fileType != "image/webp" {
		return ErrorPictureBackgroundNotFix
	}
	if file.Size/1024 > 8000 {
		return ErrorPictureBackgroundTooBig
	}
	return SUCCESS
}

func CheckMusicLIsValidate(file *multipart.FileHeader) int {
	if file == nil {
		return TargetNotExit
	}
	fileType := file.Header["Content-Type"][0]
	//(aac,m4a) mp3 flac amr ape wma mp4
	if fileType != "audio/mp4" && fileType != "audio/mpeg" && fileType != "audio/x-flac" && fileType != "audio/amr" && fileType != "audio/x-ape" && fileType != "audio/x-ms-wma" && fileType != "video/mp4" {
		return ErrorMusicLNotFix
	}
	if file.Size/1024 > 30000 {
		return ErrorMusicLNotTooBig
	}
	return SUCCESS
}

func CheckMusicWIsValidate(file *multipart.FileHeader) int {
	if file == nil {
		return TargetNotExit
	}
	fileType := file.Header["Content-Type"][0]
	//lrc,txt,doc,docx,xml
	if fileType != "ation/octet-stream" && fileType != "text/html" && fileType != "text/plain" && fileType != "application/msword" && fileType != "application/vnd.openxmlformats-officedocument.wordprocessingml.document" && fileType != "text/xml" {
		return ErrorMusicWNotFix
	}
	if file.Size/1024 > 2000 {
		return ErrorMusicWNotTooBig
	}
	return SUCCESS
}
