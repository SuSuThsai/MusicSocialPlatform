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
	fmt.Println(client, "44444")
	file, _ := filename.Open()
	defer file.Close()

	// Create a context with a cancel function for the RPC call
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call the TransferFile RPC method and get a stream object
	stream, err := client.TransferFile(ctx)
	if err != nil {
		log.Printf("Failed to call TransferFile: %v", err)
	}

	// Create a buffer of 64KB to read the file data
	buffer := make([]byte, 64*1024)

	for {
		// Read a chunk of data from the file into the buffer
		n, err := file.Read(buffer)

		// Check if we have reached the end of the file or encountered an error
		if err == io.EOF {
			break // Break out of the loop if we have reached the end of the file
		}
		if err != nil {
			log.Fatalf("Failed to read file: %v", err) // Log an error and exit if we encountered an error other than EOF
		}

		// Create a request message with the data chunk and the filetype identifier (only for the first chunk)
		req := &FileRequest{
			Data: buffer[:n], // Slice the buffer to get only the valid bytes read from the file
			Type: "",         // Set an empty string for type by default
		}

		// Set filetype only for first chunk
		if n == len(buffer) {
			req.Type = filetype
		}

		// Send the request message to the stream
		if err = stream.Send(req); err != nil {
			log.Printf("Failed to send request: %v", err)
		}

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
