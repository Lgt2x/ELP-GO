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
		fmt.Printf("Couldn't listen on port %s. This port may already be in use\n", PORT)
		return
	}

	fmt.Println("Server TCP created on port " + PORT)

	// Connection number, to identify unique connections
	connId := 0

	// Infinite loop for TCP message handling
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Printf("Couldn't accept connection %v\n", connId)
		}
		// Goroutine to handle connection
		go handleConnection(conn, connId)

		// Increase id for next connection
		connId++
	}
}

// Handles a new connection initiated with a client
// Client must send a filter id, a filename and a image blob
func handleConnection(connection net.Conn, connId int) {
	fmt.Printf("New connection with a client, id %d\n", connId)

	// Send available filters as an array
	filterList := []string{"Black & white", "Invert color"}
	fmt.Println("Sending filter list to the client")
	elputils.SendArray(connection, filterList)

	// Receive filter and send back "1" when valid, "0" otherwise
	elputils.ReceiveFilter(connection)

	// Receive modified image name
	fmt.Println("Nom image")
	fileName := elputils.ReceiveString(connection, '\n')
	fmt.Printf("Target image %s\n", fileName)

	// Receive image blob
	fmt.Println("Receiving image...")
	name := elputils.ReceiveFile(connection, "serv_rec.jpg")
	fmt.Println(name)
	// Apply filter
	fmt.Println("Applying filter")

	imageTest := elputils.ImportImage("serv_rec.jpg")
	elputils.WriteToFile(elputils.NegativeRGB(imageTest))
	//fileModified := "img_modif.jpg"

	// Rename & send the file
	fmt.Printf("Sending back %s\n", fileName)

	elputils.UploadFile(connection, "img_modif.jpg")

	// Close connection
	fmt.Printf("Closing connection with client %d\n", connId)
	connection.Close()

	// Delete temp files
	elputils.DeleteFile("img_modif.jpg")
	//elputils.DeleteFile("serv_rec.jpg")
}
