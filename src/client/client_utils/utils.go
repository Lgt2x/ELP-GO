package client_utils

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const BUFFERSIZE = 1024

func SendString(conn net.Conn, chaine string) {
	// send the string chaine
	io.WriteString(conn, fmt.Sprint(chaine))
}

func InputString() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	filtre, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return filtre
}

func InputFilter(conn net.Conn) {
	filtre := InputString()
	SendString(conn, filtre)

	validationServer := ReceiveString(conn, '\n')
	filtre_valide := strings.Compare(validationServer[0:1], "1")

	if filtre_valide != 0 {
		fmt.Println("Saisie invalide")
		InputFilter(conn)
	}
}

func InputImagePath() (string, string) {
	fmt.Print("Saisie du chemin de l'image (relatif): ")
	image_path := InputString()
	image_path_abs, _ := filepath.Abs(image_path[:len(image_path)-1])

	if FileExists(image_path_abs) {
		return image_path[:len(image_path)-1], image_path_abs
	}
	return InputImagePath()
}

func FileExists(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	}
	return false
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
	fmt.Println("File has been sent, closing connection!")
	return
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

func ReceiveString(conn net.Conn, delimiter byte) string {
	message, err := bufio.NewReader(conn).ReadString(delimiter)
	if err != nil {
		panic(err)
	}
	return message
}
