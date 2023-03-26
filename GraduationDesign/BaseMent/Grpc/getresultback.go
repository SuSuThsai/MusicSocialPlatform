package Grpc

import (
	"GraduationDesign/BaseMent/Config"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
)

// TransferFile A function to send a file data and an identifier to the server and receive a string array and an identifier from the server
func TransferFile(filename *multipart.FileHeader, filetype string) ([]string, string) {
	client := NewFileTransferClient(Config.FileRpc)
	file, _ := filename.Open()
	defer file.Close()

	data, _ := io.ReadAll(file)

	// Call the TransferFile RPC method and get a stream object
	stream, err := client.TransferFile(context.Background())
	if err != nil {
		log.Printf("Failed to call TransferFile: %v", err)
	}
	for len(data) > 0 {
		chunkSize := 64 * 1024
		if len(data) < chunkSize {
			chunkSize = len(data)
		}
		req := &FileRequest{Data: data[:chunkSize], Type: filetype}
		data = data[chunkSize:]
		if err = stream.Send(req); err != nil {
			log.Println("failed to send request:", err)
		}
		//接收服务器的响应，并将结果打印出来
		//var respond interface{}
		//err = stream.RecvMsg(respond)
		//if err != nil {
		//	log.Fatalf("failed to receive response: %v", err)
		//}
		//log.Println("response id = , result = ", respond)
	}

	// Close and receive final response from server
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("Failed to receive response: %v", err)
	}
	// Print out response message fields
	fmt.Println("Received string array:", res.Strings)
	fmt.Println("Received type:", res.Type)
	return res.Strings, res.Type
}
