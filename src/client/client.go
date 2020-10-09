package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// numéro de port établi au préalable
	PORT := "8080"

	// connexion au serveur
	conn, err := net.Dial("tcp", "localhost:"+PORT)
	defer conn.Close()

	if err != nil {
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
	//time.Sleep(1 * time.Second)

	// attente réception nom image modifiée
	fmt.Println("Attente réception nom de l'image modifiée")
	//filename_modified := receiveString(conn, '\n')
	//fmt.Println(filename_modified)

	// attente réception image modifiée
	fmt.Println("Attente de l'image modifiée")
	//receiveFile(conn, filename_modified)
	receiveFile(conn, image_path[:len(image_path)-4]+"_modified.txt")
}

func sendString(conn net.Conn, chaine string) {
	io.WriteString(conn, fmt.Sprint(chaine))
}

func inputString() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	filtre, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return filtre
}

func saisieFiltre(conn net.Conn) {
	filtre := inputString()
	sendString(conn, filtre)

	validationServer := receiveString(conn, '\n')
	filtre_valide := strings.Compare(validationServer[0:1], "1")

	if filtre_valide != 0 {
		saisieFiltre(conn)
	}
}

func inputImagePath() (string, string) {
	fmt.Print("Saisie du chemin de l'image: ")
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

func receiveFile(conn net.Conn, dstFile string) {
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

func receiveString(conn net.Conn, delimiter byte) string {
	message, err := bufio.NewReader(conn).ReadString(delimiter)
	if err != nil {
		panic(err)
	}
	return message
}
