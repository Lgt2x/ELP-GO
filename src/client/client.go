package main

import (
	utils "client/client_utils"
	"fmt"
	"net"
	"os"
)

const BUFFERSIZE = 1024

func main() {
	// numéro de port établi au préalable
	PORT := "8080"

	// connexion au serveur
	conn, err := net.Dial("tcp", "localhost:"+PORT)

	if err != nil {
		fmt.Println("Couldn't connect to server")
		return
	}

	defer conn.Close()

	// attendre réception liste filtres serveurs
	listeFilters := utils.ReceiveString(conn, '\t')
	fmt.Println(listeFilters)

	// saisie du filtre
	utils.InputFilter(conn)

	// affichage chemin
	dir, err := os.Getwd()
	fmt.Println(dir)

	// saisie nom fichier image + validation (exist or not)
	image_path, image_path_abs := utils.InputImagePath()

	// envoi nom image on envoie le chemin relatif car on est en local + soucis de creation
	fmt.Println("Envoi du nom de l'image")
	utils.SendString(conn, image_path+"\n")

	// envoi de l'image
	// time.Sleep(1 * time.Second)
	fmt.Println("Envoi de l'image:", image_path_abs)
	utils.UploadFile(conn, image_path)
	//time.Sleep(1 * time.Second)

	// attente réception nom image modifiée
	fmt.Println("Attente réception nom de l'image modifiée")
	//filename_modified := receiveString(conn, '\n')
	//fmt.Println(filename_modified)

	//utils.ReceiveFile(conn, image_path[:len(image_path)-4]+"_modified.txt")
	utils.ReceiveFile(conn)
}
