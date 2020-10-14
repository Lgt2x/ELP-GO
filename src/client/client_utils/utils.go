package client_utils

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func SendString(conn net.Conn, chaine string) {
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
		// Try again
		InputFilter(conn)
	}
}

func InputImagePath() (string, string) {
	fmt.Print("Saisie du chemin de l'image: ")
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
	// open file to upload
	fi, err := os.Open(srcFile)
	defer fi.Close()

	if err != nil {
		return
	}

	// upload
	_, err = io.Copy(conn, fi)
	if err != nil {
		return
	}
	fmt.Println("Fin de l'upload")

}

func ReceiveFile(conn net.Conn, dstFile string) {
	// create new file
	fo, err := os.Create(dstFile)
	if err != nil {
		return
	}
	os.Open(dstFile)

	// accept file from client & write to new file
	_, err = io.Copy(fo, conn)
	if err != nil {
		return
	}
}

func ReceiveString(conn net.Conn, delimiter byte) string {
	message, err := bufio.NewReader(conn).ReadString(delimiter)
	if err != nil {
		panic(err)
	}
	return message
}
