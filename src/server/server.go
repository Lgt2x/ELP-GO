package main

import (
	"ELP-GO/src/elputils"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {

	// Get port from argc, or use default value 8080
	var PORT string
	if len(os.Args) == 2 {
		PORT = os.Args[1]
	} else {
		PORT = "8080"
	}

	// Create temp dir for files
	_ = os.Mkdir("tmp", 0755)
	//defer os.RemoveAll("tmp")

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
// Client must send a filter id and a image blob
func handleConnection(connection net.Conn, connId int) {
	fmt.Printf("New connection with a client, id %d\n", connId)

	// Send available filters as contatenated array
	fmt.Println("Sending filter list to the client")
	elputils.SendArray(connection, elputils.FilterList)

	// Receive filter and send back 1 or 0 wether it's valid or not
	filter := elputils.ReceiveFilter(connection, len(elputils.FilterList))

	// Receive image blob and store it in a temp file
	fmt.Println("Receiving image...")
	og_name := "tmp/og_" + strconv.Itoa(connId) + ".jpg"
	modif_name := "tmp/modif_" + strconv.Itoa(connId) + ".jpg"
	elputils.ReceiveFile(connection, og_name)

	// Apply filter
	fmt.Println("Applying filter")
	imageTest := elputils.FileToImage(og_name)
	convert := elputils.Dispatch(imageTest, filter)
	elputils.ImageToFile(convert, modif_name)
	//fileModified := "img_modif.jpg"

	// Send back the file
	fmt.Printf("Sending back the modified image\n")
	elputils.UploadFile(connection, modif_name)

	// Close connection
	fmt.Printf("Closing connection with client %d\n\n", connId)
	connection.Close()

	// Delete temp files
	elputils.DeleteFile(og_name)
	elputils.DeleteFile(modif_name)
}
