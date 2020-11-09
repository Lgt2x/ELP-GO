package main

import (
	"ELP-GO/src/elputils"
	"fmt"
	"net"
	"os"
)

func main() {

	// Get port from argc, or use default value 8080
	var PORT string
	if len(os.Args) == 2 {
		PORT = os.Args[1]
	} else {
		PORT = "8080"
	}

	// Listen on TCP port PORT
	ln, err := net.Listen("tcp", ":"+PORT)

	if err != nil {
		fmt.Printf("Couldn't listen on port %d. This port may already be in use\n", PORT)
		return
	}

	fmt.Println("Server TCP created on port " + PORT)

	// Connection number, to identify unique connections
	conn_id := 0

	// Infinite loop for TCP message handling
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Printf("Couldn't accept connection %v\n", conn_id)
		}
		// Goroutine to handle connection
		go handleConnection(conn, conn_id)

		// Increase id for next connection
		conn_id++
	}
}

// Handles a new connection initiated with a client
// Client must send a filter id, a filename and a image blob
func handleConnection(connection net.Conn, conn_id int) {
	fmt.Printf("New connection with a client, id %d\n", conn_id)

	// Send available filters
	liste_filtres := "1. Black & White \n2. Something else"
	elputils.SendString(connection, liste_filtres+"\t")

	// Receive filter and send back "1" when valid, "0" otherwise
	elputils.ValideFiltre(connection)

	// Receive modified image name
	fmt.Println("Nom image")
	fileName := elputils.ReceiveString(connection, '\n')
	fmt.Printf("Target image %s\n", fileName)

	// Receive image blob
	fmt.Println("Receiving image...")
	elputils.ReceiveFile(connection)

	// Apply filter
	fmt.Println("Applying filter")
	fileModified := elputils.NewName(fileName)

	// Rename & send the file
	fmt.Println("Sending back %s", fileName)
	os.Rename(fileName, fileModified)
	elputils.UploadFile(connection, fileModified)

	// Close connection
	fmt.Printf("Closing connection with client %d\n", conn_id)
	connection.Close()

	// Delete temp file
	elputils.DeleteFile(fileModified)
}
