package Config

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"log"
	"net/http"
	"net/url"
)

func InitTencentCos() {
	cosUrl, _ := url.Parse(Conf.TencentCos.Url)
	secretId := Conf.TencentCos.SecretId
	secretKey := Conf.TencentCos.SecretKey
	b := &cos.BaseURL{BucketURL: cosUrl}
	TCCos = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})
	log.Println("TencentCosInfo")
}
