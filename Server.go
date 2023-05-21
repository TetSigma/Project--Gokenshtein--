package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
)

const (
	uploadDirectory = "./uploads"
	fileBufferSize  = 1
)

func main() {
	os.MkdirAll(uploadDirectory, os.ModePerm)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}

	fmt.Println("Server started. Listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading command:", err)
		return
	}
	command := string(buf[0])

	switch command {
	case "U":
		handleUpload(conn)
	case "D":
		handleDownload(conn)
	case "L":
		fileList, err := getFileList()
		if err != nil {
			fmt.Println("Error getting file list:", err)
			return
		}
		_, err = conn.Write([]byte(fileList))
		if err != nil {
			fmt.Println("Error sending file list:", err)
			return
		}
		fmt.Println("File list sent successfully")
	default:
		fmt.Println("Unknown command:", command)
	}
}

func getFileList() (string, error) {
	files, err := ioutil.ReadDir(uploadDirectory)
	if err != nil {
		fmt.Println("Error while getting file list:", err)
		return "", err
	}

	fileList := ""
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileList += file.Name() + "\n"
	}

	return fileList, nil
}

func handleUpload(conn net.Conn) {
	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading filename length:", err)
		return
	}
	filenameLength := int(buf[0])

	filenameBuf := make([]byte, filenameLength)
	_, err = conn.Read(filenameBuf)
	if err != nil {
		fmt.Println("Error reading filename:", err)
		return
	}
	filename := string(filenameBuf)

	filePath := filepath.Join(uploadDirectory, filename)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, fileBufferSize)
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading file data:", err)
			}
			break
		}

		_, err = file.Write(buffer)
		if err != nil {
			fmt.Println("Error writing file data:", err)
			return
		}
	}

	fmt.Println("File uploaded successfully:", filename)
}

func handleDownload(conn net.Conn) {
	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading filename length:", err)
		return
	}
	filenameLength := int(buf[0])

	filenameBuf := make([]byte, filenameLength)
	_, err = conn.Read(filenameBuf)
	if err != nil {
		fmt.Println("Error reading filename:", err)
		return
	}
	filename := string(filenameBuf)

	filePath := filepath.Join(uploadDirectory, filename)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, fileBufferSize)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading file data:", err)
			}
			break
		}

		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error writing file data:", err)
			return
		}
	}

	fmt.Println("File downloaded successfully:", filename)
}
