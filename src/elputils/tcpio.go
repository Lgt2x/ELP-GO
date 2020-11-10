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

const BUFFERSIZE = 1024

// Send a string using the specified connection object
func SendString(conn net.Conn, chaine string) {
	// send the string chaine
	_, err := io.WriteString(conn, fmt.Sprint(chaine))
	if err != nil {
		panic(err)
	}
}

// Receive a string
func ReceiveString(conn net.Conn, delimiter byte) string {
	message, err := bufio.NewReader(conn).ReadString(delimiter)
	if err != nil {
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
func ReceiveArray(conn net.Conn, delimiter string, delimEnd byte) []string {
	message, err := bufio.NewReader(conn).ReadString(delimEnd)
	if err != nil {
		panic(err)
	}
	return strings.Split(message, delimiter)
}

// Send a file specified a filename
func UploadFile(conn net.Conn, srcFile string) {
	file, err := os.Open(srcFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := FillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := FillString(fileInfo.Name(), 64)
	fmt.Printf("Sending filename %s and filesize!\n", fileName)
	_, _ = conn.Write([]byte(fileSize))
	_, _ = conn.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)

		if err == io.EOF {
			break
		}

		_, _ = conn.Write(sendBuffer)
	}
	fmt.Println("File has been sent")
	return
}

// Receive a file and copy it to specificed location
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
		if (fileSize - receivedBytes) < BUFFERSIZE {
			fmt.Println(fileSize, receivedBytes)
			_, _ = io.CopyN(newFile, conn, fileSize-receivedBytes)
			_, _ = conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		_, _ = io.CopyN(newFile, conn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Printf("Received %d bytes\n", fileSize)

}
