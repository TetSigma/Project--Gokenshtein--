package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

const serverAddress = "localhost:8080"

func main() {
	fmt.Println("1. Upload file")
	fmt.Println("2. Download file")
	fmt.Println("3. List all files")
	fmt.Print("Choose an option: ")

	var option int
	fmt.Scanln(&option)

	switch option {
	case 1:
		uploadFile()
	case 2:
		downloadFile()
	case 3:
		getList()
	default:
		fmt.Println("Invalid option")
	}
}

func getList() {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("L"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error receiving file list:", err)
		return
	}

	fileList := string(buffer[:n])
	fmt.Println("File list:")
	fmt.Println(fileList)
}

func uploadFile() {
	fmt.Print("Enter the file path: ")
	var filePath string
	fmt.Scanln(&filePath)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("U"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return
	}

	filenameLength := len(filepath.Base(filePath))
	_, err = conn.Write([]byte{byte(filenameLength)})
	if err != nil {
		fmt.Println("Error sending filename length:", err)
		return
	}

	_, err = conn.Write([]byte(filepath.Base(filePath)))
	if err != nil {
		fmt.Println("Error sending filename:", err)
		return
	}

	buffer := make([]byte, 1)
	for {
		_, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading file data:", err)
			}
			break
		}

		_, err = conn.Write(buffer)
		if err != nil {
			fmt.Println("Error sending file data:", err)
			return
		}
	}

	fmt.Println("File uploaded successfully")
}

func downloadFile() {
	fmt.Print("Enter the filename: ")
	var filename string
	fmt.Scanln(&filename)

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("D"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return
	}

	filenameLength := len(filename)
	_, err = conn.Write([]byte{byte(filenameLength)})
	if err != nil {
		fmt.Println("Error sending filename length:", err)
		return
	}

	_, err = conn.Write([]byte(filename))
	if err != nil {
		fmt.Println("Error sending filename:", err)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 1)
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error receiving file data:", err)
			}
			break
		}

		_, err = file.Write(buffer)
		if err != nil {
			fmt.Println("Error writing file data:", err)
			return
		}
	}

	fmt.Println("File downloaded successfully")
}
