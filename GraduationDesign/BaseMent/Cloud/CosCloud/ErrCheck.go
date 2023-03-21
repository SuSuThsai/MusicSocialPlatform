package CosCloud

import (
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"log"
)

func CheckErr(err error) bool {
	if err == nil {
		return true
	}
	if cos.IsNotFoundError(err) {
		// WARN
		log.Println("WARN: Resource is not existed")
		return false
	} else if e, ok := cos.IsCOSError(err); ok {
		log.Println(fmt.Sprintf("ERROR: Code: %v\n", e.Code))
		log.Println(fmt.Sprintf("ERROR: Message: %v\n", e.Message))
		log.Println(fmt.Sprintf("ERROR: Resource: %v\n", e.Resource))
		log.Println(fmt.Sprintf("ERROR: RequestID: %v\n", e.RequestID))
		return false
		// ERROR
	} else {
		log.Println(fmt.Sprintf("ERROR: %v\n", err))
		return false
		// ERROR
	}
}
