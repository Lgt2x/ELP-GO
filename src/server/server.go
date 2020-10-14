package main

import (
	utils "ELP-GO/src/server/server_utils"
	"fmt"
	"net"
)

func main() {
	// numéro de port établi au préalable
	// TODO : pouvoir le changer via un paramètre argc
	const PORT string = "8080"

	// le processus actuel sera dédié à l'écoute des échanges TCP sur le port PORT
	ln, err := net.Listen("tcp", ":"+PORT)

	if err != nil {
		fmt.Println("Couldn't listen on " + PORT + ". Is this port already in use ?")
		return
	}
	fmt.Println("Server TCP created on port " + PORT)

	// Connection number, to identify unique connection
	conn_id := 1

	// Infinite loop for TCP message handling
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Printf("Couldn't accept connection %v\n", conn_id)
		}
		// goroutine to handle connection
		go handleConnection(conn, conn_id)

		conn_id += 1
	}
}

func handleConnection(c net.Conn, conn_id int) {
	fmt.Println("Handling connection %v", conn_id)

	/******* INPUT HANDLING *****/
	// envoi des possibilités de filtres
	liste_filtres := "1. Noir et blanc \n2. Autre filtre"
	utils.SendString(c, liste_filtres+"\t")

	// réception du filtre, vérification de la validité et renvoi d'un caractère (1: choix valide, 0: choix invalide)
	// si 1, passage au point suivant
	validationFiltre := false
	for validationFiltre != true {
		choixFiltre := utils.ReceiveString(c)
		fmt.Println(choixFiltre)

		switch choixFiltre {
		// enumeration des choix valides
		case "1":
			// validation et go pour machin truc
			fmt.Println("Choix valide")
			validationFiltre = true
			utils.SendString(c, "1\n")
		default:
			fmt.Println("Choix invalide")
			utils.SendString(c, "0\n")
		}
	}

	// attente de réception du nom de l'image
	fmt.Println("Attente de la réception du nom de l'image")
	filename := utils.ReceiveString(c)
	fmt.Println("Nom image: ", filename)

	// attente de réception de l'image
	fmt.Println("Réception de l'image")
	utils.ReceiveFile(c, filename)

	/**** SERVER RESPONSE ****/
	// appliquer le filtre
	fmt.Println("Application du filtre")

	// envoi du nouveau nom
	//sendString(c, "image_modifiee.txt\n")

	// renvoyer le fichier avec le nom modifié
	//uploadFile(c, "image_modifiee.txt")

	// fermer la connection
	c.Close()

	// supprimer le fichier d'image
	utils.DeleteFile("image_modifiee.txt")
}
