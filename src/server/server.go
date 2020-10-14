package main

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

func handleConnection(c net.Conn, numconn int){
	fmt.Println("New connection with a client")

	// envoi des possibilités de filtres
	liste_filtres := "1. Noir et blanc \n2. Autre filtre"
	sendString(c, liste_filtres+"\t")

	// réception du filtre, vérification de la validité et renvoi d'un caractère (1: choix valide, 0: choix invalide)
	// si 1, passage au point suivant
	valideFiltre(c)

	// réception nom image
	fmt.Println("Nom image")
	fileName := receiveString(c)
	fmt.Println(fileName)

	// attente de réception de l'image
	fmt.Println("Réception de l'image")
	receiveFile(c)

	// appliquer le filtre
	fmt.Println("Application du filtre")
	fileModified := newName(fileName)

	// rename the file
	fmt.Println("Rename the file")
	os.Rename(fileName, fileModified)

	// renvoyer le fichier avec le nom modifié
	uploadFile(c, fileModified)


	// fermer la connection
	c.Close()
	fmt.Println("Goodbye", numconn)

	// supprimer le fichier d'image
	deleteFile(fileModified)
}

func newName(filename string) string {
	// find the extension, store it then delete it
	indexExt := strings.Index(filename, ".")
	return filename[:indexExt] + "_modified"+filename[indexExt:]
}

func valideFiltre(c net.Conn) {
	validationFiltre := false
	for (validationFiltre != true){
		choixFiltre := receiveString(c)
		fmt.Println(choixFiltre)

		switch choixFiltre {
		// enumeration des choix valides
		case "1":
			// validation et go pour machin truc
			fmt.Println("Choix valide")
			validationFiltre = true
			sendString(c, "1\n")
		default:
			fmt.Println("Choix invalide")
			sendString(c, "0\n")
		}
	}
}

func main() {
	// numéro de port établi au préalable
	PORT := "8080"

	// le processus actuel sera dédié à l'écoute des échanges TCP sur le port PORT
	ln, err := net.Listen("tcp", ":"+PORT)

	if err != nil {
		// handle error
		return
	}

	numconn := 1

	fmt.Println("Server TCP created on port "+PORT)

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println("Connexion impossible")
		}
		go handleConnection(conn, numconn)

		numconn += 1
	}
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
	fmt.Println("File has been sent")
	return
}

func deleteFile(filename string) {
	// delete the file filename
	err := os.Remove(filename)
	if err != nil {
		fmt.Println("Suppression impossible")
	}
}

func receiveString(conn net.Conn) string {
	// read the buffer
	message, err := bufio.NewReader(conn).ReadString('\n')
	if (err != nil) {
		panic(err)
	}
	// return the message minus the last element (like \n)
	return message[:len(message)-1]
}

func sendString(conn net.Conn, chaine string) {
	// send the string chaine on the connection conn
	io.WriteString(conn, fmt.Sprint(chaine))
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