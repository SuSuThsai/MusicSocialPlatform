package FtpAndSsh

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/utils"
	"io"
	"log"
	"mime/multipart"
	"path"
)

func UploadFileMusicL(file *multipart.FileHeader) int {
	filename := file.Filename
	//audio/xxx
	//x := file.Header["Content-Type"][0]
	//xxx := ""
	//switch x {
	//case "audio/mp4":
	//	xxx = "m4a" //"aac is the same"
	//case "audio/mpeg":
	//	xxx = "mp3"
	//case "audio/x-flac":
	//	xxx = "flac"
	//case "audio/amr":
	//	xxx = "amr"
	//case "audio/x-ape":
	//	xxx = "ape"
	//case "audio/x-ms-wma":
	//	xxx = "wma"
	//case "video/mp4":
	//	xxx = "mp4"
	//}
	//fmt.Println(xxx)
	//filename = filename + "." + xxx
	f, _ := file.Open()
	defer f.Close()
	dstFile, err := Config.SftpClient.Create(path.Join("/MusicPlatform/Musics/MusicL/", filename))
	//dstFile.Chmod(os.ModePerm)
	if err != nil {
		log.Println("sftpClient.Create error : ", err)
		return utils.ERROR
	}
	defer dstFile.Close()
	ff, err := io.ReadAll(f)
	if err != nil {
		log.Println("ReadAll error : ", err)
		return utils.ERROR
	}
	_, err = dstFile.Write(ff)
	if err != nil {
		log.Println("Write error : ", err)
		return utils.ERROR
	}
	return utils.SUCCESS
}

func UploadFileMusicW(file *multipart.FileHeader) int {
	filename := file.Filename
	//audio/xxx
	//x := file.Header["Content-Type"][0]
	//fmt.Println(x, filename)
	//xxx := ""
	//switch x {
	//case "application/octet-stream":
	//	xxx = "lrc"
	//case "text/html":
	//	xxx = "html"
	//case "text/plain":
	//	xxx = "txt"
	//case "application/msword":
	//	xxx = "doc"
	//case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
	//	xxx = "docx"
	//case "text/xml":
	//	xxx = "xml"
	//}
	//fmt.Println(xxx)
	//filename = filename + "." + xxx
	f, _ := file.Open()
	defer f.Close()
	dstFile, err := Config.SftpClient.Create(path.Join("/MusicPlatform/Musics/MusicW/", filename))
	//dstFile.Chmod(os.ModePerm)
	if err != nil {
		log.Println("sftpClient.Create error : ", err)
		return utils.ERROR
	}
	defer dstFile.Close()
	ff, err := io.ReadAll(f)
	if err != nil {
		log.Println("ReadAll error : ", err)
		return utils.ERROR
	}
	_, err = dstFile.Write(ff)
	if err != nil {
		log.Println("Write error : ", err)
		return utils.ERROR
	}
	return utils.SUCCESS
}
