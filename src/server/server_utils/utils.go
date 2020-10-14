package server_utils

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func ReceiveFile(conn net.Conn, dstFile string) {
	// create new file
	fo, err := os.Create(dstFile)
	if err != nil {
		fmt.Println("Erreur création fichier")
		return
	}
	//fmt.Println("Création fichier réussie")

	// accept file from client & write to new file
	_, err = io.Copy(fo, conn)
	if err != nil {
		fmt.Println("Erreur réception fichier")
		return
	}
	fmt.Println("Fin de réception du fichier")
}

func UploadFile(conn net.Conn, srcFile string) {
	// open file to upload
	fi, err := os.Open(srcFile)
	if err != nil {
		return
	}

	// upload
	_, err = io.Copy(conn, fi)
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
