package server_utils

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

func NewName(filename string) string {
	// find the extension, store it then delete it
	indexExt := strings.Index(filename, ".")
	return filename[:indexExt] + "_modified" + filename[indexExt:]
}

func ValideFiltre(c net.Conn) {
	validationFiltre := false
	for validationFiltre != true {
		choixFiltre := ReceiveString(c)
		fmt.Println(choixFiltre)

		switch choixFiltre {
		// enumeration des choix valides
		case "1":
			// validation et go pour machin truc
			fmt.Println("Choix valide")
			validationFiltre = true
			SendString(c, "1\n")
		default:
			fmt.Println("Choix invalide")
			SendString(c, "0\n")
		}
	}
}

func ReceiveFile(conn net.Conn) {
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	conn.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	conn.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64

	fmt.Println("Start receiving")
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, conn, (fileSize - receivedBytes))
			conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, conn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")
}

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
	fmt.Println("Sending filename and filesize!")
	conn.Write([]byte(fileSize))
	conn.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)

		if err == io.EOF {
			break
		}

		conn.Write(sendBuffer)
	}
	fmt.Println("File has been sent")
	return
}

func DeleteFile(filename string) {
	// delete the file filename
	err := os.Remove(filename)
	if err != nil {
		fmt.Println("Suppression impossible")
	}
}

func ReceiveString(conn net.Conn) string {
	// read the buffer
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		panic(err)
	}
	// return the message minus the last element (like \n)
	return message[:len(message)-1]
}

func SendString(conn net.Conn, chaine string) {
	// send the string chaine on the connection conn
	io.WriteString(conn, fmt.Sprint(chaine))
}

func FillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
