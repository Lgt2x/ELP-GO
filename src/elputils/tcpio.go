// Utility function for tcp networking
// used both client-side and server-side
package elputils

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const BufferSize = 1024

// Send a string using the specified connection object
func SendString(conn net.Conn, string string) {
	// send the string string
	_, err := io.WriteString(conn, fmt.Sprint(string))
	if err != nil {
		fmt.Println("Error - SendString")
		panic(err)
	}
}

// Receive a string
func ReceiveString(conn net.Conn, delimiter byte) string {
	message, err := bufio.NewReader(conn).ReadString(delimiter)
	if err != nil {
		fmt.Println("Error - ReceiveString")
		panic(err)
	}
	return message
}

// Send an array of strings using a semi-colon as a separator
func SendArray(conn net.Conn, array []string) {
	_, err := io.WriteString(conn, fmt.Sprint(strings.Join(array, ";")+"\n"))
	if err != nil {
		panic(err)
	}
}

// Receive an array of strings
func ReceiveArray(conn net.Conn, delimiter string, delimitEnd byte) []string {
	message, err := bufio.NewReader(conn).ReadString(delimitEnd)
	if err != nil {
		panic(err)
	}
	return strings.Split(strings.Trim(message, "\n"), delimiter)
}

// Send a file specified a filename
func UploadFile(conn net.Conn, srcFile string) {
	file, err := os.Open(srcFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := FillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := FillString(fileInfo.Name(), 64)

	_, _ = conn.Write([]byte(fileSize))
	_, _ = conn.Write([]byte(fileName))
	sendBuffer := make([]byte, BufferSize)
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}

		_, _ = conn.Write(sendBuffer)
	}
}

// Receive a file and copy it to specified location
func ReceiveFile(conn net.Conn, destination string) {
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	_, _ = conn.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	_, _ = conn.Read(bufferFileName)
	//fileName := strings.Trim(string(bufferFileName), ":")
	newFile, err := os.Create(destination)

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64

	fmt.Println("Start receiving")
	for {
		if (fileSize - receivedBytes) < BufferSize {
			_, _ = io.CopyN(newFile, conn, fileSize-receivedBytes)
			_, _ = conn.Read(make([]byte, (receivedBytes+BufferSize)-fileSize))
			break
		}
		_, _ = io.CopyN(newFile, conn, BufferSize)
		receivedBytes += BufferSize
	}
	fmt.Printf("Received %d bytes\n", fileSize)

}
