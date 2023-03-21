package main

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/Grpc"
	"fmt"
)

func main() {
	Config.InitsConfig()
	fmt.Println(Config.Conf.FileGrpc, "11111")
	Config.CreateClientConn()
	fmt.Println(Config.FileRpc, "33333")
	Grpc.TransferFile(nil, "2")
}
