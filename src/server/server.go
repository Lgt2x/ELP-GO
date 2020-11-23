package main

import (
	"ELP-GO/src/elputils"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {

	var PORT string
	var ROUTINES int64

	if len(os.Args) == 3 {
		PORT = os.Args[1]
		ROUTINES, _ = strconv.ParseInt(os.Args[2], 10, 32)
	} else {
		log.Println("No parameters provided, using default server config")
		log.Println("Usage : server <port> <routines/img>")
		PORT = "8000"
		ROUTINES = 4
	}

	// Create temp dir for files
	_ = os.Mkdir("/tmp/ELP_GO", 0755)

	// Listen on TCP port
	ln, err := net.Listen("tcp", ":"+PORT)

	if err != nil {
		log.Printf("Couldn't listen on port %s. This port may already be in use\n", PORT)
		return
	}

	log.Println("Server TCP created on port " + PORT)
	log.Printf("Using %d routines per image\n", ROUTINES)

	// Connection number, to identify unique connections
	connId := 0

	// Infinite loop for TCP message handling
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			log.Printf("Couldn't accept connection %v\n", connId)
		}
		// Goroutine to handle connection
		go handleConnection(conn, connId, int(ROUTINES))

		// Increase id for next connection
		connId++
	}
}

// Function to avoid the server to crash if a panic statement is raised due to a client
func defusePanic(connId int) {
	if r := recover(); r != nil {
		log.Println("recovered from ", r)
		elputils.PrintRedLn("Panic due to client " + strconv.Itoa(connId) + ". Stopping connection.")
	}
}

// Handles a new connection initiated with a client
// Client must send a filter id and a image blob
func handleConnection(connection net.Conn, connId int, routines int) {
	// defusePanic will be executed if there is a panic statement raised
	defer defusePanic(connId)

	start := time.Now()

	log.Printf("New connection with a client, id %d\n", connId)

	// Send available filters as contatenated array
	log.Println("Sending filter list to the client")
	elputils.SendArray(connection, elputils.FilterList)
	defer connection.Close()

	// Receive filter and send back 1 or 0 wether it's valid or not
	filter := elputils.ReceiveFilter(connection, len(elputils.FilterList))

	// Receive image blob and store it in a temp file
	log.Println("Receiving image...")
	ogName := "/tmp/ELP_GO/og_" + strconv.Itoa(connId) + ".jpg"
	modifiedName := "/tmp/ELP_GO/modified_" + strconv.Itoa(connId) + ".jpg"
	elputils.ReceiveFile(connection, ogName)

	// Apply filter
	log.Println("Applying filter")
	imageTest := elputils.FileToImage(ogName)
	convert := elputils.ApplyFilterAsync(imageTest, filter, routines)
	elputils.ImageToFile(convert, modifiedName)

	// Send back the file
	log.Printf("Sending back the modified image\n")
	elputils.UploadFile(connection, modifiedName)

	elapsed := time.Since(start)
	log.Printf("Image took %s to process", elapsed)

	// Close connection
	log.Printf("Closing connection with client %d\n\n", connId)
	connection.Close()

	// Delete temp files
	elputils.DeleteFile(ogName)
	elputils.DeleteFile(modifiedName)
}
