package main

import (
	utils "server/server_utils"
	"fmt"
	"net"
	"os"
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

func handleConnection(c net.Conn, numconn int) {
	fmt.Println("New connection with a client")

	// envoi des possibilités de filtres
	liste_filtres := "1. Noir et blanc \n2. Autre filtre"
	utils.SendString(c, liste_filtres+"\t")

	// réception du filtre, vérification de la validité et renvoi d'un caractère (1: choix valide, 0: choix invalide)
	// si 1, passage au point suivant
	utils.ValideFiltre(c)

	// réception nom image
	fmt.Println("Nom image")
	fileName := utils.ReceiveString(c)
	fmt.Println(fileName)

	// attente de réception de l'image
	fmt.Println("Réception de l'image")
	utils.ReceiveFile(c)

	// appliquer le filtre
	fmt.Println("Application du filtre")
	fileModified := utils.NewName(fileName)

	// rename the file
	fmt.Println("Rename the file")
	os.Rename(fileName, fileModified)

	// renvoyer le fichier avec le nom modifié
	utils.UploadFile(c, fileModified)

	// fermer la connection
	c.Close()
	fmt.Println("Goodbye", numconn)

	// supprimer le fichier d'image
	utils.DeleteFile(fileModified)
}
