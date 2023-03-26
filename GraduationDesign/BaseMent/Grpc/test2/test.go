package main

import (
	"io"
	"log"
	"net"

	pb "GraduationDesign/BaseMent/Grpc"
	"google.golang.org/grpc"
)

const (
	port = "localhost:5051"
)

type server struct{}

func (s *server) TransferFile(stream pb.FileTransfer_TransferFileServer) error {
	var requestData []byte
	var requestType string
	for {
		req, err := stream.Recv()
		if err == io.EOF { // 当客户端关闭连接时，停止读取并退出循环
			break
		}
		if err != nil {
			log.Printf("failed to receive request: %v", err)
			return err
		}
		requestData = append(requestData, req.Data...)
		requestType = req.Type
		//if err := stream.SendMsg(&pb.FileResponse{Type: "12345", Strings: []string{"received"}}); err != nil {
		//	log.Printf("failed to send response: %v", err)
		//	return err
		//}
		//调用GetResult()函数处理数据并返回结果
	}
	result, err := s.GetResult(requestData, requestType)
	if err != nil {
		log.Printf("failed to get result: %v", err)
		return err
	}
	if err = stream.SendMsg(&pb.FileResponse{Type: "12345", Strings: result}); err != nil {
		log.Printf("failed to send response: %v", err)
		return err
	}
	return nil
}

func (s *server) GetResult(data []byte, dataType string) ([]string, error) {
	//处理数据并返回结果
	return []string{"result1", "result2"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFileTransferServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
