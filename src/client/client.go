package main

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

func main() {
	// numéro de port établi au préalable
	PORT := "8080"

	// connexion au serveur
	conn, err := net.Dial("tcp", "localhost:"+PORT)
	defer conn.Close()

	if (err != nil){
		return
	}

	// attendre réception liste filtres serveurs
	listeFiltres := receiveString(conn, '\t')
	fmt.Println(listeFiltres)

	// saisie du filtre
	saisieFiltre(conn)

	// saisie nom fichier image + validation (exist or not)
	image_path, image_path_abs := inputImagePath()

	// envoi nom image on envoie le chemin relatif car on est en local + soucis de creation
	fmt.Println("Envoi du nom de l'image")
	sendString(conn, image_path+"\n")

	// envoi de l'image
	// time.Sleep(1 * time.Second)
	fmt.Println("Envoi de l'image:", image_path_abs)
	uploadFile(conn, image_path)

	// attente réception image modifiée
	fmt.Println("Attente de l'image modifiée")
	//receiveFile(conn, filename_modified)
	receiveFile(conn)
}

func sendString(conn net.Conn, chaine string) {
	// send the string chaine
	io.WriteString(conn, fmt.Sprint(chaine))
}

func inputString() (string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	filtre, err := reader.ReadString('\n')
	if (err != nil) {
		panic(err)
	}
	return filtre
}

func saisieFiltre(conn net.Conn) {
	filtre := inputString()
	sendString(conn, filtre)

	validationServer := receiveString(conn, '\n')
	filtre_valide := strings.Compare(validationServer[0:1], "1")

	if (filtre_valide != 0) {
		fmt.Println("Saisie invalide")
		saisieFiltre(conn)
	}
}

func inputImagePath() (string, string) {
	fmt.Print("Saisie du chemin de l'image (relatif): ")
	image_path := inputString()
	image_path_abs, _ := filepath.Abs(image_path[:len(image_path)-1])

	if fileExists(image_path_abs) {
		return image_path[:len(image_path)-1], image_path_abs
	}
	return inputImagePath()
}

func fileExists(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	}
	return false
}

func uploadFile(conn net.Conn, srcFile string) {
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
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
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

func fillString(retunString string, toLength int) string {
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

func receiveFile(conn net.Conn) {
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

func receiveString(conn net.Conn, delimiter byte) string {
	message, err := bufio.NewReader(conn).ReadString(delimiter)
	if (err != nil) {
		panic(err)
	}
	return message
}