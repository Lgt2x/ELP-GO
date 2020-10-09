package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func handleConnection(c net.Conn, numconn int){
	fmt.Println("New connection with a client")

	// envoi des possibilités de filtres
	liste_filtres := "1. Noir et blanc \n2. Autre filtre"
	sendString(c, liste_filtres+"\t")

	// réception du filtre, vérification de la validité et renvoi d'un caractère (1: choix valide, 0: choix invalide)
	// si 1, passage au point suivant
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

	// attente de réception du nom de l'image
	fmt.Println("Attente de la réception du nom de l'image")
	filename := receiveString(c)
	fmt.Println("Nom image: ", filename)

	// attente de réception de l'image
	fmt.Println("Réception de l'image")
	receiveFile(c, filename)

	// appliquer le filtre
	fmt.Println("Application du filtre")

	// envoi du nouveau nom
	//sendString(c, "image_modifiee.txt\n")

	// renvoyer le fichier avec le nom modifié
	//uploadFile(c, "image_modifiee.txt")

	// fermer la connection
	c.Close()

	// supprimer le fichier d'image
	deleteFile("image_modifiee.txt")
}

func valideFiltre(c net.Conn) {

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

func receiveFile(conn net.Conn, dstFile string) {
	// create new file
	fo, err := os.Create(dstFile)
	if (err != nil) {
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

func uploadFile(conn net.Conn, srcFile string) {
	// open file to upload
	fi, err := os.Open(srcFile)
	if (err != nil) {
		return
	}

	// upload
	_, err = io.Copy(conn, fi)
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