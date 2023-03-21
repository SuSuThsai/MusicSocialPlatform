package CosCloud

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/Model"
	"context"
	"mime/multipart"
	"strings"
)

var c = context.Background()

// Upload 单文件上传
func Upload(fileName string, file multipart.File) (string, bool) {
	_, err := Config.TCCos.Object.Put(c, fileName, file, nil)
	url := Config.TCCos.Object.GetObjectURL(fileName)
	return url.String(), CheckErr(err)
}

func UpLoadMusicL(file *multipart.FileHeader) (string, bool) {
	filename1 := file.Filename
	f, _ := file.Open()
	filename := "YamadaUsers/MusicUpload/MusicL/" + filename1
	url, code := Upload(filename, f)
	return url, code
}

func UpLoadMusicW(file *multipart.FileHeader) (string, bool) {
	filename1 := file.Filename
	f, _ := file.Open()
	filename := "YamadaUsers/MusicUpload/MusicW/" + filename1
	url, code := Upload(filename, f)
	return url, code
}

func UpLoadFace(file *multipart.FileHeader, userId string) (string, bool) {
	f, _ := file.Open()
	xxx := file.Header["Content-Type"][0][6:]
	filename := "YamadaUsers/" + userId + "/" + userId + "_face." + xxx
	url, code := Upload(filename, f)
	return url, code
}

func UpLoadArticlePicture(file *multipart.FileHeader, userId string, articleId string) (string, bool) {
	f, _ := file.Open()
	defer f.Close()
	xxx := file.Header["Content-Type"][0][6:]
	filename := "YamadaUsers/" + userId + "/Articles" + "/" + userId + "_" + articleId + "_articlePicture." + xxx
	url, code := Upload(filename, f)
	return url, code
}

func UpLoadBackGround(file *multipart.FileHeader, userId string) (string, bool) {
	f, _ := file.Open()
	defer f.Close()
	xxx := file.Header["Content-Type"][0][6:]
	filename := "YamadaUsers/" + userId + "/" + userId + "_background." + xxx
	url, code := Upload(filename, f)
	return url, code
}

func Delete(userId string, code int) bool {
	_, data, _, _, _ := Model.GetUser(userId)
	var url []string
	if code == 0 {
		url = strings.Split(data.Pfp, "/")
	} else {
		url = strings.Split(data.Background, "/")
	}
	fileName := "YamadaUsers/" + userId + "/" + url[len(url)-1]
	_, err := Config.TCCos.Object.Delete(c, fileName, nil)
	return CheckErr(err)
}

func DeleteArticlePicture(articleId uint) bool {
	data, _ := Model.CheckAArticle(articleId)
	var err error
	var url []string
	s := strings.Split(data.Img, " ")
	for i := 0; i < len(s); i++ {
		url = strings.Split(s[i], "/")
		fileName := "YamadaUsers/" + data.UserId + "/Articles" + "/" + url[len(url)-1]
		_, err = Config.TCCos.Object.Delete(c, fileName, nil)
	}
	return CheckErr(err)
}
